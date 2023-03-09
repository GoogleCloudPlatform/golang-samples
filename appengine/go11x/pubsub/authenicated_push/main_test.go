// Copyright 2020 Google LLC
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
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/pubsub/v1"
)

func TestReceiveMessagesHandler(t *testing.T) {
	testToken := "test-verification-token"

	tests := []struct {
		name    string
		email   string
		aud     string
		token   string
		wantErr bool
	}{
		{
			name:    "works",
			email:   "test-service-account-email@example.com",
			aud:     "http://example.com",
			token:   testToken,
			wantErr: false,
		},
		{
			name:    "bad email",
			email:   "bad-email@example.com",
			aud:     "http://example.com",
			token:   testToken,
			wantErr: true,
		},
		{
			name:    "bad token sent",
			email:   "test-service-account-email@example.com",
			aud:     "http://example.com",
			token:   "bad token",
			wantErr: true,
		},
		{
			name:    "mismatched aud claim in auth token",
			email:   "test-service-account-email@example.com",
			aud:     "http://mismatched.com",
			token:   testToken,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authToken, pk := createRS256JWT(t, tt.email, tt.aud)
			app := &app{pubsubVerificationToken: testToken}
			app.defaultHTTPClient = createClient(t, pk)
			pr := &pushRequest{
				Message: pubsub.PubsubMessage{
					Attributes: map[string]string{"key": "value"},
					Data:       "Hello Cloud Pub/Sub! Here is my message!",
					MessageId:  "136969346945",
				},
				Subscription: "test-sub",
			}
			body, err := json.Marshal(pr)
			if err != nil {
				t.Fatalf("json.Marshal(%v): got %v, want nil", pr, err)
			}

			v := url.Values{}
			v.Set("token", tt.token)
			req := httptest.NewRequest("POST", "http://example.com?"+v.Encode(), bytes.NewReader(body))
			req.URL.RawQuery = req.URL.Query().Encode()
			req.Header.Set("Authorization", "Bearer "+authToken)
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			app.receiveMessagesHandler(rr, req)

			if !tt.wantErr && rr.Code != 200 {
				t.Fatalf("code: got %v, want 200", rr.Code)
			}
			if tt.wantErr && rr.Body.String() == "OK" {
				t.Fatalf("code: got 200, want %v", rr.Code)
			}
			// only assert further for the happy path
			if tt.wantErr {
				return
			}
			if len(app.messages) != 1 {
				t.Fatalf("len(messages): got %v, want 1", len(app.messages))
			}
			if !cmp.Equal(app.messages[0], pr.Message.Data) {
				t.Fatalf("got %+v, want %+v", app.messages[0], pr)
			}
		})
	}
}

type certResponse struct {
	Keys []jwk `json:"keys"`
}

type jwt struct {
	header    string
	payload   string
	signature string
}

func (j *jwt) String() string {
	return fmt.Sprintf("%s.%s.%s", j.header, j.payload, j.signature)
}

// hashedContent gets the SHA256 checksum for verification of the JWT.
func (j *jwt) hashedContent() []byte {
	signedContent := j.header + "." + j.payload
	hashed := sha256.Sum256([]byte(signedContent))
	return hashed[:]
}

type jwtHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
	KeyID     string `json:"kid"`
}
type jwk struct {
	Alg string `json:"alg"`
	Crv string `json:"crv"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Use string `json:"use"`
	E   string `json:"e"`
	N   string `json:"n"`
}

func createRS256JWT(t *testing.T, email string, aud string) (string, rsa.PublicKey) {
	t.Helper()
	token := createAuthToken(t, email, aud)
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("unable to generate key: %v", err)
	}
	sig, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, token.hashedContent())
	if err != nil {
		t.Fatalf("unable to sign content: %v", err)
	}
	token.signature = base64.RawURLEncoding.EncodeToString(sig)
	return token.String(), privateKey.PublicKey
}

// Same as `idtoken.Payload` with the addition of `email` and `email_verified` claims
// present in Cloud Pub/Sub JWT tokens.
type ExtendedPayload struct {
	idtoken.Payload
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

func createAuthToken(t *testing.T, email string, aud string) *jwt {
	t.Helper()
	header := jwtHeader{
		KeyID:     "123",
		Algorithm: "RS256",
		Type:      "JWT",
	}
	payload := ExtendedPayload{
		Payload: idtoken.Payload{
			Issuer:   "https://accounts.google.com",
			Audience: aud,
			Expires:  time.Now().Add(1 * time.Minute).Unix(),
		},
		Email:         email,
		EmailVerified: true,
	}

	hb, err := json.Marshal(&header)
	if err != nil {
		t.Fatalf("unable to marshall header: %v", err)
	}
	pb, err := json.Marshal(&payload)
	if err != nil {
		t.Fatalf("unable to marshall payload: %v", err)
	}
	eb := base64.RawURLEncoding.EncodeToString(hb)
	ep := base64.RawURLEncoding.EncodeToString(pb)
	return &jwt{
		header:  eb,
		payload: ep,
	}
}

func createClient(t *testing.T, pk rsa.PublicKey) *http.Client {
	return &http.Client{
		Transport: RoundTripFn(func(req *http.Request) *http.Response {
			cr := certResponse{
				Keys: []jwk{
					{
						Kid: "123",
						N:   base64.RawURLEncoding.EncodeToString(pk.N.Bytes()),
						E:   base64.RawURLEncoding.EncodeToString(new(big.Int).SetInt64(int64(pk.E)).Bytes()),
					},
				},
			}
			b, err := json.Marshal(&cr)
			if err != nil {
				t.Fatalf("unable to marshal response: %v", err)
			}
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(b)),
				Header:     make(http.Header),
			}
		}),
	}
}

type RoundTripFn func(req *http.Request) *http.Response

func (f RoundTripFn) RoundTrip(req *http.Request) (*http.Response, error) { return f(req), nil }
