package app

// [START sample]
import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func logHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	post := &Post{Body: "sample post"}
	key := datastore.NewIncompleteKey(ctx, "Posts", nil)
	if _, err := datastore.Put(ctx, key, post); err != nil {
		log.Errorf(ctx, "could not put into datastore: %v", err)
		http.Error(w, "An error occurred. Try again.", http.StatusInternalServerError)
		return
	}
	log.Debugf(ctx, "Datastore put successful")

	w.Write([]byte("ok!"))
}

// [END sample]

type Post struct {
	Body string
}

func init() {
	http.HandleFunc("/log", logHandler)
}
