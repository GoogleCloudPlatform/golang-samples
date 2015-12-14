package sample

// [START OAuth_and_App_Engine]
import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

func welcomeOAuth(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	u, err := user.CurrentOAuth(ctx, "")
	if err != nil {
		http.Error(w, "OAuth Authorization header required", http.StatusUnauthorized)
		return
	}
	if !u.Admin {
		http.Error(w, "Admin login only", http.StatusUnauthorized)
		return
	}
	fmt.Fprintf(w, `Welcome, admin user %s!`, u)
}

// [END OAuth_and_App_Engine]
