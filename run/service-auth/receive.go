package main

// [START auth_validate_and_decode_bearer_token_on_go]
import (
	"fmt"
	"net/http"
	"strings"
)

// Parse the authorization header and decode the information beign
// sent by the Bearer Token
func (a *app) receiveAuthorizedRequest(w http.ResponseWriter, r *http.Request) {
	// Allows requests only for the root path ("/") to prevent duplicate calls.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Request method should be GET.
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Attempt to retrieve and validate the Authorization header.
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.Write([]byte("Hello, anonymous user\n"))
		return
	}

	if len(strings.Split(authHeader, " ")) != 2 {
		http.Error(w, "Malformed Authorization header", http.StatusBadRequest)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := strings.Split(authHeader, " ")[1]

	payload, status, err := a.validateToken(token)
	if err != nil {
		http.Error(w, err.Error(), status)
	}

	w.Write(fmt.Appendf(nil, "Hello, %s!\n", payload.Claims["email"]))
}

// [END auth_validate_and_decode_bearer_token_on_go]
