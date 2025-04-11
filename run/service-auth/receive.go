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
	"context"
	"net/http"
	"strings"

	"google.golang.org/api/idtoken"
)

func receiveAuthorizedRequest(w http.ResponseWriter, r *http.Request){
	ctx := context.Background()
	
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || len(strings.Split(authHeader, " ")) != 2 {
		http.Error(w, "Missing Authorization header", http.StatusBadRequest)
		return
	}

	token := strings.Split(authHeader, " ")[1]

	payload, err := idtoken.Validate(context.Background(), splittedAuth[1], audience)
	

}