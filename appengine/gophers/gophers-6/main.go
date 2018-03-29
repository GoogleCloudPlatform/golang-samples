// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/log"

	"cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/apiv1"
	firebase "firebase.google.com/go"
	uuid "github.com/satori/go.uuid"
)

var (
	indexTemplate = template.Must(template.ParseFiles("index.html"))
)

type Label struct {
	Description string
	Score       float32
}

type Post struct {
	Author   string
	UserID   string
	Message  string
	Posted   time.Time
	ImageURL string
	Labels   []Label
}

type templateParams struct {
	Notice  string
	Name    string
	Message string
	Posts   []Post
}

func main() {
	http.HandleFunc("/", indexHandler)
	appengine.Main()
}

// labelFunc will be called asynchronously as a Cloud Task. labelFunc can
// be executed by calling labelFunc.Call(ctx, postID). If an error is returned
// the function will be retried.
var labelFunc = delay.Func("label-image", func(ctx context.Context, id int64) error {
	// Get the post to label.
	k := datastore.NewKey(ctx, "Post", "", id, nil)
	post := Post{}
	if err := datastore.Get(ctx, k, &post); err != nil {
		log.Errorf(ctx, "getting Post to label: %v", err)
		return err
	}
	if post.ImageURL == "" {
		// Nothing to label.
		return nil
	}

	// Create a new vision client.
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Errorf(ctx, "NewImageAnnotatorClient: %v", err)
		return err
	}
	defer client.Close()

	// Get the image and label it.
	image := vision.NewImageFromURI(post.ImageURL)
	labels, err := client.DetectLabels(ctx, image, nil, 5)
	if err != nil {
		log.Errorf(ctx, "Failed to detect labels: %v", err)
		return err
	}

	for _, l := range labels {
		post.Labels = append(post.Labels, Label{
			Description: l.GetDescription(),
			Score:       l.GetScore(),
		})
	}

	// Update the database with the new labels.
	if _, err := datastore.Put(ctx, k, &post); err != nil {
		log.Errorf(ctx, "Failed to update image: %v", err)
		return err
	}
	return nil
})

// storageBucketName is the Google Cloud Storage bucket to store uploaded images under.
var storageBucketName = ""

// uploadFileFromForm uploads a file if it's present in the "image" form field.
func uploadFileFromForm(ctx context.Context, r *http.Request) (url string, err error) {
	if storageBucketName == "" {
		return "", errors.New("storage bucket is missing")
	}
	// Read the file from the form.
	f, fh, err := r.FormFile("image")
	if err == http.ErrMissingFile {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	// Ensure the file is an image. http.DetectContentType only uses 512 bytes.
	buf := make([]byte, 512)
	if _, err := f.Read(buf); err != nil {
		return "", err
	}
	if contentType := http.DetectContentType(buf); !strings.HasPrefix(contentType, "image") {
		return "", fmt.Errorf("not an image: %s", contentType)
	}
	// Reset f so subsequent calls to Read start from the beginning of the file.
	f.Seek(0, 0)

	// Create a storage client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}
	storageBucket := client.Bucket(storageBucketName)

	// Random filename, retaining existing extension.
	u, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("generating UUID: %v", err)
	}
	name := u.String() + path.Ext(fh.Filename)

	w := storageBucket.Object(name).NewWriter(ctx)
	w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	w.ContentType = fh.Header.Get("Content-Type")

	// Entries are immutable, be aggressive about caching (1 day).
	w.CacheControl = "public, max-age=86400"

	if _, err := io.Copy(w, f); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}

	const publicURL = "https://storage.googleapis.com/%s/%s"
	return fmt.Sprintf(publicURL, storageBucketName, name), nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	ctx := appengine.NewContext(r)
	params := templateParams{}

	q := datastore.NewQuery("Post").Order("-Posted").Limit(20)
	if _, err := q.GetAll(ctx, &params.Posts); err != nil {
		log.Errorf(ctx, "Getting posts: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		params.Notice = "Couldn't get latest posts. Refresh?"
		indexTemplate.Execute(w, params)
		return
	}

	if r.Method == "GET" {
		indexTemplate.Execute(w, params)
		return
	}
	// It's a POST request, so handle the form submission.

	message := r.FormValue("message")
	if message == "" {
		w.WriteHeader(http.StatusBadRequest)
		params.Notice = "No message provided"
		indexTemplate.Execute(w, params)
		return
	}

	// Create a new Firebase App.
	app, err := firebase.NewApp(ctx, &firebase.Config{
		DatabaseURL:   "copy from Firebase Console > Overview > Add Firebase to your web app",
		ProjectID:     "copy from Firebase Console > Overview > Add Firebase to your web app",
		StorageBucket: "copy from Firebase Console > Overview > Add Firebase to your web app",
	})
	if err != nil {
		params.Notice = "Couldn't authenticate. Try logging in again?"
		params.Message = message // Preserve their message so they can try again.
		indexTemplate.Execute(w, params)
		return
	}
	// Create a new authenticator for the app.
	auth, err := app.Auth(ctx)
	if err != nil {
		params.Notice = "Couldn't authenticate. Try logging in again?"
		params.Message = message // Preserve their message so they can try again.
		indexTemplate.Execute(w, params)
		return
	}
	// Verify the token passed in by the user is valid.
	tok, err := auth.VerifyIDTokenAndCheckRevoked(ctx, r.FormValue("token"))
	if err != nil {
		params.Notice = "Couldn't authenticate. Try logging in again?"
		params.Message = message // Preserve their message so they can try again.
		indexTemplate.Execute(w, params)
		return
	}
	// Use the validated token to get the user's information.
	user, err := auth.GetUser(ctx, tok.UID)
	if err != nil {
		params.Notice = "Couldn't authenticate. Try logging in again?"
		params.Message = message // Preserve their message so they can try again.
		indexTemplate.Execute(w, params)
		return
	}

	post := Post{
		UserID:  user.UID, // Include UserID in case Author isn't unique.
		Author:  user.DisplayName,
		Message: message,
		Posted:  time.Now(),
	}
	params.Name = post.Author

	// Get the image if there is one.
	imageURL, err := uploadFileFromForm(ctx, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		params.Notice = "Error saving image: " + err.Error()
		params.Message = post.Message // Preserve their message so they can try again.
		indexTemplate.Execute(w, params)
		return
	}
	post.ImageURL = imageURL

	key := datastore.NewIncompleteKey(ctx, "Post", nil)
	if key, err = datastore.Put(ctx, key, &post); err != nil {
		log.Errorf(ctx, "datastore.Put: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		params.Notice = "Couldn't add new post. Try again?"
		params.Message = post.Message // Preserve their message so they can try again.
		indexTemplate.Execute(w, params)
		return
	}

	// Only look for labels if the post has an image.
	if imageURL != "" {
		// Run labelFunc. This will start a new Task in the background.
		if err := labelFunc.Call(ctx, key.IntID()); err != nil {
			log.Errorf(ctx, "delay Call %v", err)
		}
	}

	// Prepend the post that was just added.
	params.Posts = append([]Post{post}, params.Posts...)
	params.Notice = fmt.Sprintf("Thank you for your submission, %s!", post.Author)
	indexTemplate.Execute(w, params)
}
