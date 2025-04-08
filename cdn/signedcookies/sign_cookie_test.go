// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package signedcookie

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestReadKeyFile(t *testing.T) {
	f, err := ioutil.TempFile("", "cdnkey")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.Close()

	key := `nZtRohdNF9m3cKM24IcK4w==`
	expected := []byte{0x9d, 0x9b, 0x51, 0xa2, 0x17, 0x4d, 0x17, 0xd9,
		0xb7, 0x70, 0xa3, 0x36, 0xe0, 0x87, 0x0a, 0xe3}
	if err := ioutil.WriteFile(f.Name(), []byte(key), 0600); err != nil {
		t.Fatal(err)
	}
	b, err := readKeyFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, expected) {
		t.Fatalf("the signed cookie value did not match: got=%v; expected=%v", b, expected)
	}
}

func TestSignCookie(t *testing.T) {
	testKey := []byte{0x9d, 0x9b, 0x51, 0xa2, 0x17, 0x4d, 0x17, 0xd9,
		0xb7, 0x70, 0xa3, 0x36, 0xe0, 0x87, 0x0a, 0xe3} // base64url: nZtRohdNF9m3cKM24IcK4w==

	cases := []struct {
		testName    string
		urlPrefix   string
		keyName     string
		key         []byte
		expiration  time.Time
		out         string
		shouldMatch bool
	}{
		{
			testName:    "Valid Domain and Sig",
			urlPrefix:   "https://media.example.com/segments/",
			keyName:     "my-key",
			key:         testKey,
			expiration:  time.Unix(1558131350, 0),
			out:         "URLPrefix=aHR0cHM6Ly9tZWRpYS5leGFtcGxlLmNvbS9zZWdtZW50cy8=:Expires=1558131350:KeyName=my-key:Signature=_qwhz38bxCKdiDqENLIx4ujrw-U=",
			shouldMatch: true,
		},
		{
			testName:    "Valid Domain #2 & Sig",
			urlPrefix:   "https://video.example.com/manifests/123/",
			keyName:     "my-key",
			key:         testKey,
			expiration:  time.Unix(1558131350, 0),
			out:         "URLPrefix=aHR0cHM6Ly92aWRlby5leGFtcGxlLmNvbS9tYW5pZmVzdHMvMTIzLw==:Expires=1558131350:KeyName=my-key:Signature=0ZSSGg0VdRH29siXG8wiPTV5LOE=",
			shouldMatch: true,
		},
		{
			testName:    "Mismatched Timestamps",
			urlPrefix:   "https://media.example.com/segments/",
			keyName:     "my-key",
			key:         testKey,
			expiration:  time.Unix(1558131350, 0),
			out:         "URLPrefix=aHR0cHM6Ly9tZWRpYS5leGFtcGxlLmNvbS9zZWdtZW50cy8=:Expires=1858131350:KeyName=my-key:Signature=QMZgLb8pS9MkhTxcPOQTM5nzJXc=",
			shouldMatch: false,
		},
		{
			testName:    "Mismatched Key Names",
			urlPrefix:   "https://media.example.com/segments/",
			keyName:     "bad-key",
			key:         testKey,
			expiration:  time.Unix(1558131350, 0),
			out:         "URLPrefix=aHR0cHM6Ly9tZWRpYS5leGFtcGxlLmNvbS9zZWdtZW50cy8=:Expires=1858131350:KeyName=my-key:Signature=QMZgLb8pS9MkhTxcPOQTM5nzJXc=",
			shouldMatch: false,
		},
		{
			testName:  "Invalid Key Material",
			urlPrefix: "https://media.example.com/segments/",
			keyName:   "my-key",
			key: []byte{0x9d, 0x9b, 0x51, 0xa2, 0x17, 0x4d, 0x17, 0xd9,
				0xb7, 0x70, 0xa3, 0x36, 0xe0, 0x87, 0x0a},
			expiration:  time.Unix(1558131350, 0),
			out:         "URLPrefix=aHR0cHM6Ly9tZWRpYS5leGFtcGxlLmNvbS9zZWdtZW50cy8=:Expires=1858131350:KeyName=my-key:Signature=QMZgLb8pS9MkhTxcPOQTM5nzJXc=",
			shouldMatch: false,
		},
	}

	for _, c := range cases {
		t.Run(
			fmt.Sprintf("%s (shouldMatch: %t)", c.testName, c.shouldMatch), func(t *testing.T) {
				signedValue, err := signCookie(c.urlPrefix, c.keyName,
					c.key, c.expiration)
				if err != nil {
					t.Errorf("SignCookie returned an error: %s", err)
				}

				if signedValue != c.out && c.shouldMatch {
					t.Errorf("signed value did not match: got: %s, want: %s", signedValue, c.out)
				}

				// Test for invalid matches - e.g. where the strings are empty,
				// or where the same signature is being calculated across test
				// cases.
				if signedValue == c.out && !c.shouldMatch {
					t.Errorf("signed value incorrectly matched: got %s, want %s", signedValue, c.out)
				}

			})
	}
}
