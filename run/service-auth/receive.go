// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

// Parse the authorization header and decode the information beign
// sent by the Bearer Token
func receiveAuthorizedRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Attempt to retrieve and validate the Authorization header.
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.Write([]byte("Hello, anonymus user\n"))
		return
	}

	if len(strings.Split(authHeader, " ")) != 2 {
		http.Error(w, "Malformed Authorization header", http.StatusBadRequest)
		return
	}

	token := strings.Split(authHeader, " ")[1]

	// Verify and decode the JWT
	v, err := idtoken.NewValidator(r.Context(), option.WithHTTPClient(http.DefaultClient))
	if err != nil {
		http.Error(w, "Unable to create Validator", http.StatusBadRequest)
		return
	}

	payload, err := v.Validate(r.Context(), token, "")
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid Token: %v", err), http.StatusBadRequest)
		return
	}

	w.Write([]byte(fmt.Sprintf("Hello, %s!\n", payload.Claims["email"])))
}
