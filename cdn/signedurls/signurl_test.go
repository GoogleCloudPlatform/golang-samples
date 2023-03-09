// Copyright 2019 Google LLC
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
package signedurl

import (
	"bytes"
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

func TestSignURL(t *testing.T) {
	testKey := []byte{0x9d, 0x9b, 0x51, 0xa2, 0x17, 0x4d, 0x17, 0xd9,
		0xb7, 0x70, 0xa3, 0x36, 0xe0, 0x87, 0x0a, 0xe3} // base64url: nZtRohdNF9m3cKM24IcK4w==

	cases := []struct {
		testName   string
		url        string
		keyName    string
		expiration time.Time
		out        string
	}{
		{
			testName:   "Domain and Exact Path",
			url:        "http://35.186.234.33/index.html",
			keyName:    "my-key",
			expiration: time.Unix(1558131350, 0),
			out:        "http://35.186.234.33/index.html?Expires=1558131350&KeyName=my-key&Signature=fm6JZSmKNsB5sys8VGr-JE4LiiE=",
		},
		{
			testName:   "Domain Only",
			url:        "https://www.google.com/",
			keyName:    "my-key",
			expiration: time.Unix(1549751401, 0),
			out:        "https://www.google.com/?Expires=1549751401&KeyName=my-key&Signature=M_QO7BGHi2sGqrJO-MDr0uhDFuc=",
		},
		{
			testName:   "With Query Params",
			url:        "https://www.example.com/some/path?some=query&another=param",
			keyName:    "my-key",
			expiration: time.Unix(1549751461, 0),
			out:        "https://www.example.com/some/path?some=query&another=param&Expires=1549751461&KeyName=my-key&Signature=sTqqGX5hUJmlRJ84koAIhWW_c3M=",
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			signedValue := signURL(
				c.url, c.keyName, testKey, c.expiration,
			)
			if signedValue != c.out {
				t.Errorf("signed value incorrectly matched: got %s, want %s", signedValue, c.out)
			}
		})
	}
}

func TestSignURLWithPrefix(t *testing.T) {
	testKey := []byte{0x9d, 0x9b, 0x51, 0xa2, 0x17, 0x4d, 0x17, 0xd9,
		0xb7, 0x70, 0xa3, 0x36, 0xe0, 0x87, 0x0a, 0xe3} // base64url: nZtRohdNF9m3cKM24IcK4w==

	cases := []struct {
		testName   string
		urlPrefix  string
		keyName    string
		expiration time.Time
		out        string
	}{
		{
			testName:   "Domain and Simple Prefix",
			urlPrefix:  "https://media.example.com/segments/",
			keyName:    "my-key",
			expiration: time.Unix(1558131350, 0),
			out:        "URLPrefix=aHR0cHM6Ly9tZWRpYS5leGFtcGxlLmNvbS9zZWdtZW50cy8=&Expires=1558131350&KeyName=my-key&Signature=HWE5tBTZgnYVoZzVLG7BtRnOsgk=",
		},
		{
			testName:   "Domain Only",
			urlPrefix:  "https://www.google.com/",
			keyName:    "my-key",
			expiration: time.Unix(1549751401, 0),
			out:        "URLPrefix=aHR0cHM6Ly93d3cuZ29vZ2xlLmNvbS8=&Expires=1549751401&KeyName=my-key&Signature=o0zZ77jb7BgtGRPQEaXmX3cCLh8=",
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			signedValue, err := signURLWithPrefix(
				c.urlPrefix, c.keyName, testKey, c.expiration,
			)
			if err != nil {
				t.Error(err)
			}

			if signedValue != c.out {
				t.Errorf("signed value incorrectly matched: got %s, want %s", signedValue, c.out)
			}
		})
	}
}
