// Copyright 2023 Google LLC
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
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

func TestReceiveAuthorizedGetRequest(t *testing.T) {

	tests := []struct {
		name           string
		header         string
		wantStatusCode int
		wantResponse   string
	}{
		{
			name:           "Empty Authorization header",
			header:         "",
			wantStatusCode: http.StatusOK,
			wantResponse:   "Hello, anonymous user.\n",
		},
		{
			name:           "Invalid Authorization header format",
			header:         "InvalidHeaderFormat",
			wantStatusCode: http.StatusOK,
			wantResponse:   "Unhandled header format (InvalidHeaderFormat).\n",
		},
		{
			name:           "Valid Bearer token",
			header:         "Bearer " + createToken("1234567890@test.com"),
			wantStatusCode: http.StatusOK,
			wantResponse:   "Hello, 1234567890@test.com!\n",
		},
		{
			name:           "Invalid Bearer token",
			header:         "Bearer " + "invalid-token",
			wantStatusCode: http.StatusOK,
			wantResponse:   "Unable to parse token: token contains an invalid number of segments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Authorization", tt.header)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(receiveAuthorizedGetRequest)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatusCode)
			}

			if rr.Body.String() != tt.wantResponse {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.wantResponse)
			}

		})
	}
}

func createToken(email string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
	})
	fmt.Printf("token: %+v\n", token)
	// ReWrite <my-secret-key> for your secret key in the same way to run/service-auth/receive.go
	tokenString, _ := token.SignedString([]byte("my-secret-key"))
	return tokenString
}
