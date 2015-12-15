package sample

// [START creating_namespaces_on_a_per_user_basis]
import (
	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

func namespace(ctx context.Context) context.Context {
	// assumes the user is logged in.
	ctx, err := appengine.Namespace(ctx, user.Current(ctx).ID)
	if err != nil {
		// ...
	}
	return ctx
}

// [END creating_namespaces_on_a_per_user_basis]
