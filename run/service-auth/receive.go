package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func main() {
	http.HandleFunc("/", receiveAuthorizedGetRequest)
	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	// Start HTTP server.
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func receiveAuthorizedGetRequest(w http.ResponseWriter, r *http.Request) {
	/*
	   receiveAuthorizedGetRequest takes the "Authorization" header from a
	   request, decodes it using the jwt-go library, and returns back the email
	   from the header to the caller.
	*/
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		fmt.Fprintf(w, "Hello, anonymous user.\n")
		return
	}

	// split the auth type and value from the header.
	authValues := strings.SplitN(authHeader, " ", 2)
	if len(authValues) != 2 {
		fmt.Fprintf(w, "Unhandled header format (%v).\n", authHeader)
		return
	}
	authType, creds := authValues[0], authValues[1]
	if authType == "Bearer" {
		token, err := jwt.Parse(creds, func(token *jwt.Token) (interface{}, error) {
			// Check if signing algorithm is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// TODO: Set secret key
			return []byte("my-secret-key"), nil
		})
		
		if err != nil {
			fmt.Fprintf(w, "Unable to parse token: %v", err)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Fprintf(w, "Unable to extract claims from token")
			return
		}
		email, ok := claims["email"].(string)
		if !ok {
			fmt.Fprintf(w, "Unable to extract email from token")
			return
		}
		fmt.Fprintf(w, "Hello, %v!\n", email)
	} else {
		fmt.Fprintf(w, "Unhandled header format (%v).\n", authType)
	}
}
