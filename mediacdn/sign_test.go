// Copyright 2022 Google LLC
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
package snippets

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

var privateTestKey = []byte{34, 31, 185, 24, 168, 225, 242, 115, 112, 155, 38,
	157, 183, 65, 104, 243, 85, 182, 188, 26, 176, 101, 247, 177,
	243, 93, 114, 156, 94, 191, 219, 75, 183, 211, 110, 78, 223,
	133, 62, 172, 159, 217, 158, 126, 34, 6, 254, 108, 57, 194,
	141, 93, 219, 91, 8, 162, 88, 62, 52, 75, 42, 103, 202, 238,
}

var buf = &bytes.Buffer{}

func TestSignURL(t *testing.T) {
	cases := []struct {
		testName   string
		url        string
		keyName    string
		expiration time.Time
		out        string
	}{
		{
			testName:   "Domain and exact path",
			url:        "http://35.186.234.33/index.html",
			keyName:    "my-key",
			expiration: time.Unix(1558131350, 0),
			out:        "http://35.186.234.33/index.html?Expires=1558131350&KeyName=my-key&Signature=bwCkNAIuVneG0cRPwwPDk1vGmMfqR_TbFfLguwdsfF8Pdlk8INOKICYVOTHY5jHlGgwSF2jkRkm8bWZGwu-SAw",
		},
		{
			testName:   "Domain only",
			url:        "https://www.example.com/",
			keyName:    "my-key",
			expiration: time.Unix(1549751401, 0),
			out:        "https://www.example.com/?Expires=1549751401&KeyName=my-key&Signature=xCkdQFmY6zAsWjdEjCpbaZEnlrK0KF3it5nVqcjDo5gvC3LvERf6wn0DdrNt7-ZSAuGzXDF51pWc1Ye0AqNeBw",
		},
		{
			testName:   "With query params",
			url:        "https://www.example.com/some/path?some=query&another=param",
			keyName:    "my-key",
			expiration: time.Unix(1549751461, 0),
			out:        "https://www.example.com/some/path?some=query&another=param&Expires=1549751461&KeyName=my-key&Signature=kM8uoFD9tfNKqOe1ulQpWUutBL4oQERxcR6sCg-brtPOSGJXqvuUOyEP1EsGzVCesI6epkY4AxYC9yCAuY1GDQ",
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			if err := signURL(buf, c.url, c.keyName, privateTestKey, c.expiration); err != nil {
				t.Errorf("signURL got err: %v", err)
			}
			if got := buf.String(); !strings.Contains(got, c.out) {
				t.Errorf("signed value incorrectly matched: got %q, want %q", got, c.out)
			}
		})
	}
}

func TestSignURLPrefix(t *testing.T) {
	cases := []struct {
		testName   string
		url        string
		keyName    string
		expiration time.Time
		out        string
	}{
		{
			testName:   "Domain and exact path",
			url:        "http://35.186.234.33/index.html",
			keyName:    "my-key",
			expiration: time.Unix(1558131350, 0),
			out:        "http://35.186.234.33/index.html?URLPrefix=aHR0cDovLzM1LjE4Ni4yMzQuMzMvaW5kZXguaHRtbA&Expires=1558131350&KeyName=my-key&Signature=EabJ40MSKAB1hb7IkiNVowQI3N0cLapjDeGQ2oJUx_1n5Blypl31V7auj-KaiGMVMpYCSGxY48G6G-x7G6xNBg",
		},
		{
			testName:   "Domain only",
			url:        "https://www.google.com/",
			keyName:    "my-key",
			expiration: time.Unix(1549751401, 0),
			out:        "https://www.google.com/?URLPrefix=aHR0cHM6Ly93d3cuZ29vZ2xlLmNvbS8&Expires=1549751401&KeyName=my-key&Signature=f82Yhq9HrFXuAKNKlKpt7qk3e1BKo2OCtIy6JF0HA2j_l1IUF69ZFBXposUSky_fgvVvTpxi9IOJCONTKiMNDw",
		},
		{
			testName:   "With query params",
			url:        "https://www.example.com/some/path?some=query&another=param",
			keyName:    "my-key",
			expiration: time.Unix(1549751461, 0),
			out:        "https://www.example.com/some/path?some=query&another=param&URLPrefix=aHR0cHM6Ly93d3cuZXhhbXBsZS5jb20vc29tZS9wYXRoP3NvbWU9cXVlcnkmYW5vdGhlcj1wYXJhbQ&Expires=1549751461&KeyName=my-key&Signature=zx7rPX3Zol2y3Uu9HuL_IQFe0LiOp556Z40rlAJjxBwiQ6vwbA2BvqvkrSNb3VWOWbxpnI4ssHxtMn7sQHh9CQ",
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			if err := signURLPrefix(buf, c.url, c.keyName, privateTestKey, c.expiration); err != nil {
				t.Errorf("signURLPrefix got err: %v", err)
			}
			if got := buf.String(); !strings.Contains(got, c.out) {
				t.Errorf("signed value incorrectly matched: got %q, want %q", got, c.out)
			}
		})
	}
}

func TestSignCookie(t *testing.T) {
	cases := []struct {
		testName   string
		url        string
		keyName    string
		expiration time.Time
		out        string
	}{
		{
			testName:   "Domain only",
			url:        "https://www.google.com/",
			keyName:    "my-key",
			expiration: time.Unix(1549751401, 0),
			out:        "Edge-Cache-Cookie=URLPrefix=aHR0cHM6Ly93d3cuZ29vZ2xlLmNvbS8:Expires=1549751401:KeyName=my-key:Signature=O67Laog-pcQ2_RNOuVrgGiN5NS-16I0SOItQRnW0yDkbawgVgX9KfFCgdoqXpY0P3f8ZdMEM2tEVsU6-Saq9BA",
		},
		{
			testName:   "Domain and exact path",
			url:        "https://www.example.com/some",
			keyName:    "my-key",
			expiration: time.Unix(1549751461, 0),
			out:        "Edge-Cache-Cookie=URLPrefix=aHR0cHM6Ly93d3cuZXhhbXBsZS5jb20vc29tZQ:Expires=1549751461:KeyName=my-key:Signature=MjRwgGa4vJJ5lkVt1xJSoi-LyMk5x-bf1AmUBr-2XiB6zP4LSqHsmQZoeZA4fVw6C7HCcNqQT1UzGPgGe7bpAQ",
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			if err := signCookie(buf, c.url, c.keyName, privateTestKey, c.expiration); err != nil {
				t.Errorf("signCookie got err: %v", err)
			}
			if got := buf.String(); !strings.Contains(got, c.out) {
				t.Errorf("signed value incorrectly matched: got %q, want %q", got, c.out)
			}
		})
	}
}
