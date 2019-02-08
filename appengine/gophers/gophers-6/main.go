// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	firebase "firebase.google.com/go"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	// [START new_imports]
	"context"
	"io"
	"path"
	"strings"

	"cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/apiv1"
	uuid "github.com/gofrs/uuid"
	"google.golang.org/appengine/delay"
	// [END new_imports]
)

var (
	firebaseConfig = &firebase.Config{
		DatabaseURL:   "https://console.firebase.google.com > Overview > Add Firebase to your web app",
		ProjectID:     "https://console.firebase.google.com > Overview > Add Firebase to your web app",
		StorageBucket: "https://console.firebase.google.com > Overview > Add Firebase to your web app",
	}
	indexTemplate = template.Must(template.ParseFiles("index.html"))
)

// [START label_struct]

// A Label is a description for a post's image.
type Label struct {
	Description string
	Score       float32
}

// [END label_struct]

// [START new_post_fields]

type Post struct {
	Author   string
	UserID   string
	Message  string
	Posted   time.Time
	ImageURL string
	Labels   []Label
}

// [END new_post_fields]

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

// [START var_label_func]

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

// [END var_label_func]

// [START upload_image]

// uploadFileFromForm uploads a file if it's present in the "image" form field.
func uploadFileFromForm(ctx context.Context, r *http.Request) (url string, err error) {
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
	storageBucket := client.Bucket(firebaseConfig.StorageBucket)

	// Random filename, retaining existing extension.
	u, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("generating UUID: %v", err)
	}
	name := u.String() + path.Ext(fh.Filename)

	w := storageBucket.Object(name).NewWriter(ctx)

	// Warning: storage.AllUsers gives public read access to anyone.
	w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	w.ContentType = fh.Header.Get("Content-Type")

	// Entries are immutable, be aggressive about caching (1 day).
	w.CacheControl = "public, max-age=86400"

	if _, err := io.Copy(w, f); err != nil {
		w.CloseWithError(err)
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}

	const publicURL = "https://storage.googleapis.com/%s/%s"
	return fmt.Sprintf(publicURL, firebaseConfig.StorageBucket, name), nil
}

// [END upload_image]

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
	app, err := firebase.NewApp(ctx, firebaseConfig)
	if err != nil {
		params.Notice = "Couldn't authenticate. Try logging in again?"
		log.Errorf(ctx, "firebase.NewApp: %v", err)
		params.Message = message // Preserve their message so they can try again.
		indexTemplate.Execute(w, params)
		return
	}
	// Create a new authenticator for the app.
	auth, err := app.Auth(ctx)
	if err != nil {
		params.Notice = "Couldn't authenticate. Try logging in again?"
		log.Errorf(ctx, "app.Auth: %v", err)
		params.Message = message // Preserve their message so they can try again.
		indexTemplate.Execute(w, params)
		return
	}
	// Verify the token passed in by the user is valid.
	tok, err := auth.VerifyIDTokenAndCheckRevoked(ctx, r.FormValue("token"))
	if err != nil {
		params.Notice = "Couldn't authenticate. Try logging in again?"
		log.Errorf(ctx, "auth.VerifyIDAndCheckRevoked: %v", err)
		params.Message = message // Preserve their message so they can try again.
		indexTemplate.Execute(w, params)
		return
	}
	// Use the validated token to get the user's information.
	user, err := auth.GetUser(ctx, tok.UID)
	if err != nil {
		params.Notice = "Couldn't authenticate. Try logging in again?"
		log.Errorf(ctx, "auth.GetUser: %v", err)
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

	// [START image_URL]
	// Get the image if there is one.
	imageURL, err := uploadFileFromForm(ctx, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		params.Notice = "Error saving image: " + err.Error()
		params.Message = post.Message // Preserve their message so they can try again.
		indexTemplate.Execute(w, params)
		return
	}
	// [END image_URL]

	// [START add_image_URL]
	post.ImageURL = imageURL
	// [END add_image_URL]

	key := datastore.NewIncompleteKey(ctx, "Post", nil)
	if key, err = datastore.Put(ctx, key, &post); err != nil {
		log.Errorf(ctx, "datastore.Put: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		params.Notice = "Couldn't add new post. Try again?"
		params.Message = post.Message // Preserve their message so they can try again.
		indexTemplate.Execute(w, params)
		return
	}

	// [START empty_image]
	// Only look for labels if the post has an image.
	if imageURL != "" {
		// Run labelFunc. This will start a new Task in the background.
		if err := labelFunc.Call(ctx, key.IntID()); err != nil {
			log.Errorf(ctx, "delay Call %v", err)
		}
	}
	// [END empty_image]

	// Prepend the post that was just added.
	params.Posts = append([]Post{post}, params.Posts...)
	params.Notice = fmt.Sprintf("Thank you for your submission, %s!", post.Author)
	indexTemplate.Execute(w, params)
}
