// Copyright 2019 Google LLC
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

package tictactoe

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/zabawaba99/firego"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine"
)

func firebase(ctx context.Context) (*firego.Firebase, error) {
	hc, err := google.DefaultClient(ctx,
		"https://www.googleapis.com/auth/firebase.database",
		"https://www.googleapis.com/auth/userinfo.email",
	)
	if err != nil {
		return nil, err
	}
	base := os.Getenv("FIREBASE_BASE")
	if base == "" || strings.Contains(base, "YOUR-PROJECT-ID") {
		// Check the environment variable for the base firebase URL.
		//
		// The config should look like:
		//
		// env_variables:
		//    FIREBASE_BASE: https://app-id.firebaseio.com
		//
		return nil, errors.New("Unset FIREBASE_BASE environment variable.")
	}
	return firego.New(base, hc), nil
}

func createToken(ctx context.Context, channelID string) (string, error) {
	iss, err := appengine.ServiceAccount(ctx)
	if err != nil {
		return "", err
	}
	iat := time.Now().Unix()
	jwt := map[string]interface{}{
		"iss": iss,
		"sub": iss,
		"aud": "https://identitytoolkit.googleapis.com/google.identity.identitytoolkit.v1.IdentityToolkit",
		"iat": iat,
		"exp": iat + 3600, // 1 hour
		"uid": channelID,
	}
	body, err := json.Marshal(jwt)
	if err != nil {
		return "", err
	}
	header := base64.StdEncoding.EncodeToString([]byte(`{"typ":"JWT","alg":"RS256"}`))
	payload := append([]byte(header), byte('.'))
	payload = append(payload, []byte(base64.StdEncoding.EncodeToString(body))...)
	_, sig, err := appengine.SignBytes(ctx, payload)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s", payload, base64.StdEncoding.EncodeToString(sig)), nil
}
