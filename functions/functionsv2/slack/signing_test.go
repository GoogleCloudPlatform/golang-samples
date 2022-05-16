// Copyright 2022 Google LLC
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

package slack

import (
	"encoding/hex"
	"fmt"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestSigningWithSecret(t *testing.T) {
	type testCase struct {
		name       string
		signature  string
		timeStamp  string
		wantResult bool
	}

	secret := "talesfromthecrypt"
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	body := "somebody"
	base := fmt.Sprintf("%s:%s:%s", version, ts, body)
	correctSHA2Signature := fmt.Sprintf("%s=%s", version, hex.EncodeToString(getSignature([]byte(base), []byte(secret))))

	tests := []testCase{
		{name: "Good request", signature: correctSHA2Signature, timeStamp: ts, wantResult: true},
		{name: "Bad signature", signature: "v0=146abde6763faeba19adc4d9fe4961668f4be11f7405a1c05b636f29312eac2e", timeStamp: ts, wantResult: false},
		{name: "Old timestamp", signature: correctSHA2Signature, timeStamp: "12345", wantResult: false},
	}

	for _, tc := range tests {
		rq := httptest.NewRequest("POST", "https://someurl.com", strings.NewReader("somebody"))
		rq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rq.Header.Add("X-Slack-Request-Timestamp", tc.timeStamp)
		rq.Header.Add("X-Slack-Signature", tc.signature)

		got, err := verifyWebHook(rq, []byte(""), secret)
		if err != nil {
			// Any error other then the expected one is a failed test.
			if !strings.Contains(err.Error(), "checkTimestamp") {
				t.Errorf("verifyWebHook: %v", err)
			}
		}
		if tc.wantResult != got {
			t.Errorf("Test: %v - Wanted: %v but got: %v", tc.name, tc.wantResult, got)
		}
	}
}
