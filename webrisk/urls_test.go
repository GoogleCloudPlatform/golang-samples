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
package webrisk

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestGeneratePatterns(t *testing.T) {
	vectors := []struct {
		url    string
		output []string
		fail   bool
	}{{
		url:    "http://a.b.c/1/2.html?param=1/2",
		output: []string{"a.b.c/1/2.html?param=1/2", "a.b.c/1/2.html", "a.b.c/1/", "a.b.c/", "b.c/1/2.html?param=1/2", "b.c/1/2.html", "b.c/1/", "b.c/"},
	}, {
		url:    "http://a.b.c/1/2/3/4/5/",
		output: []string{"a.b.c/1/2/3/4/5/", "a.b.c/1/2/3/", "a.b.c/1/2/", "a.b.c/1/", "a.b.c/", "b.c/1/2/3/4/5/", "b.c/1/2/3/", "b.c/1/2/", "b.c/1/", "b.c/"},
	}, {
		url:    "http://a.b.c.d.e.f.g.h.i/",
		output: []string{"a.b.c.d.e.f.g.h.i/", "c.d.e.f.g.h.i/", "d.e.f.g.h.i/", "e.f.g.h.i/", "f.g.h.i/", "g.h.i/", "h.i/"},
	}, {
		url:    "http://a.b.c.d.e/1.html",
		output: []string{"a.b.c.d.e/1.html", "a.b.c.d.e/", "b.c.d.e/1.html", "b.c.d.e/", "c.d.e/1.html", "c.d.e/", "d.e/1.html", "d.e/"},
	}, {
		url:    "http://b.c/1/2/3.html?param=1/2",
		output: []string{"b.c/1/2/3.html?param=1/2", "b.c/1/2/3.html", "b.c/1/2/", "b.c/1/", "b.c/"},
	}, {
		url:    "http://a.b/?",
		output: []string{"a.b/"},
	}, {
		url:    "http://[2001:470:1:18::114]/a/b",
		output: []string{"[2001:470:1:18::114]/a/b", "[2001:470:1:18::114]/a/", "[2001:470:1:18::114]/"},
	}, {
		url:    "http://1.2.3.4/a/b",
		output: []string{"1.2.3.4/a/b", "1.2.3.4/a/", "1.2.3.4/"},
	}, {
		url:    "http://a.b/",
		output: []string{"a.b/"},
	}, {
		url:    "http://b/",
		output: []string{"b/"},
	}, {
		url:    "https://a.b.c.d.e.f.g.h.i/",
		output: []string{"a.b.c.d.e.f.g.h.i/", "c.d.e.f.g.h.i/", "d.e.f.g.h.i/", "e.f.g.h.i/", "f.g.h.i/", "g.h.i/", "h.i/"},
	}, {
		url:    "a.b.c.d.e.f.g.h.i/",
		output: []string{"a.b.c.d.e.f.g.h.i/", "c.d.e.f.g.h.i/", "d.e.f.g.h.i/", "e.f.g.h.i/", "f.g.h.i/", "g.h.i/", "h.i/"},
	}, {
		url:    "[2001:470:1:18::114]/a/b",
		output: []string{"[2001:470:1:18::114]/a/b", "[2001:470:1:18::114]/a/", "[2001:470:1:18::114]/"},
	}, {
		url:  "/asdf",
		fail: true,
	}}

	for i, v := range vectors {
		patterns, err := generatePatterns(v.url)
		if err != nil != v.fail {
			if err != nil {
				t.Errorf("test %d, unexpected error: %v", i, err)
			} else {
				t.Errorf("test %d, unexpected success", i)
			}
			continue
		}
		sort.Strings(patterns)
		sort.Strings(v.output)
		if !reflect.DeepEqual(patterns, v.output) {
			t.Errorf("test %d, generatePatterns(%q):\ngot  %q\nwant %q", i, v.url, patterns, v.output)
		}
	}
}

func TestParseIPAddress(t *testing.T) {
	vectors := []struct {
		url    string
		output string
	}{
		{"123.123.0.0.1", ""},
		{"255.0.0.1", "255.0.0.1"},
		{"12.0x12.01234", "12.18.2.156"},
		{"276.2.3", "20.2.0.3"},
		{"012.034.01.055", "10.28.1.45"},
		{"0x12.0x43.0x44.0x01", "18.67.68.1"},
		{"167838211", "10.1.2.3"},
		{"3279880203", "195.127.0.11"},
		{"4294967295", "255.255.255.255"},
		{"10.192.95.89 xy", "10.192.95.89"},
		{"1.2.3.00x0", ""},
	}
	for i, v := range vectors {
		iphost := parseIPAddress(v.url)
		if iphost != v.output {
			t.Errorf("test %d, parseIPAddress(%q) = %q, want %q", i, v.url, iphost, v.output)
		}
	}
}

func TestCanonicalHost(t *testing.T) {
	vectors := []struct {
		url    string
		output string
		fail   bool
	}{
		{"http://www.google.com/foo.html", "www.google.com", false},
		{"http://google.com./foo.html", "google.com", false},
		{"http://google.com.:8080/foo.html", "google.com", false},
		{"http://google...com/foo.html", "google.com", false},
		{"http://..google.com/foo.html", "google.com", false},
		{"http://[FEDC:BA98:7654:3210:FEDC:BA98:7654:3210]:80/index.html",
			strings.ToLower("[FEDC:BA98:7654:3210:FEDC:BA98:7654:3210]"), false},
		{"http://[::192.9.5.5]/ipng", "[::192.9.5.5]", false},
		{"http://0x12.0x43.0x44.0x01", "18.67.68.1", false},
		{"http://192.168.0.1:80/index.html", "192.168.0.1", false},
		{"/asdf", "", true},
	}

	for i, v := range vectors {
		host, err := canonicalHost(v.url)
		if err != nil != v.fail {
			if err != nil {
				t.Errorf("test %d url %v, unexpected error: %v", i, v.url, err)
			} else {
				t.Errorf("test %d, unexpected success", i)
			}
			continue
		}
		if host != v.output {
			t.Errorf("test %d, canonicalHost(%q) = %q, want %q", i, v.url, host, v.output)
		}
	}
}

func TestGenerateLookupHosts(t *testing.T) {
	vectors := []struct {
		url    string
		output []string
		fail   bool
	}{{
		url:    "http://www.google.com/foo.html",
		output: []string{"www.google.com", "google.com"},
	}, {
		url:    "http://a.b.c.com/foo.html",
		output: []string{"a.b.c.com", "b.c.com", "c.com"},
	}, {
		url:    "http://a.b.c.d.e.f.kita.tokyo.jp",
		output: []string{"a.b.c.d.e.f.kita.tokyo.jp", "c.d.e.f.kita.tokyo.jp", "d.e.f.kita.tokyo.jp", "e.f.kita.tokyo.jp", "f.kita.tokyo.jp", "kita.tokyo.jp", "tokyo.jp"},
	}, {
		url:    "http://[::192.9.5.5]/ipng",
		output: []string{"[::192.9.5.5]"},
	}, {
		url:  "/asdf",
		fail: true,
	}}

	for i, v := range vectors {
		hosts, err := generateLookupHosts(v.url)
		if err != nil != v.fail {
			if err != nil {
				t.Errorf("test %d, unexpected error: %v", i, err)
			} else {
				t.Errorf("test %d, unexpected success", i)
			}
			continue
		}
		if !reflect.DeepEqual(hosts, v.output) {
			t.Errorf("test %d, generateLookupHosts(%q):\ngot  %q\nwant %q", i, v.url, hosts, v.output)
		}
	}
}

func TestCanonicalPath(t *testing.T) {
	vectors := []struct {
		url    string
		output string
		fail   bool
	}{
		{"http://a.com", "/", false},
		{"http://a.com/foo.html", "/foo.html", false},
		{"http://a.com/foo/.././bar/./../foo.html", "/foo.html", false},
		{"http://a.com/a/b/", "/a/b/", false},
		{"http://a.com/a/b/c", "/a/b/c", false},
		{"http://a.com//a//b///c////", "/a/b/c/", false},
		{"http://%31%36%38%2e%31%38%38%2e%39%39%2e%32%36/%2E%73%65%63%75%72%65/%77%77%77%2E%65%62%61%79%2E%63%6F%6D/?query#fragment",
			"/.secure/www.ebay.com/", false}, // "http://168.188.99.26/.secure/www.ebay.com/"
		{"http://195.127.0.11/uploads/%20%20%20%20/.verify/.eBaysecure=updateuserdataxplimnbqmn-xplmvalidateinfoswqpcmlx=hgplmcx/",
			"/uploads/%20%20%20%20/.verify/.eBaysecure=updateuserdataxplimnbqmn-xplmvalidateinfoswqpcmlx=hgplmcx/", false},
		{"http://host%23.com/%257Ea%2521b%2540c%2523d%2524e%25f%255E00%252611%252A22%252833%252944_55%252B",
			"/~a!b@c%23d$e%25f^00&11*22(33)44_55+", false},
		{"/asdf", "", true},
	}

	for i, v := range vectors {
		path, err := canonicalPath(v.url)
		if err != nil != v.fail {
			if err != nil {
				t.Errorf("test %d, unexpected error: %v", i, err)
			} else {
				t.Errorf("test %d, unexpected success", i)
			}
			continue
		}
		if path != v.output {
			t.Errorf("test %d, canonicalPath(%q) = %q, want %q", i, v.url, path, v.output)
		}
	}
}

func TestGenerateLookupPaths(t *testing.T) {
	vectors := []struct {
		url    string
		output []string
		fail   bool
	}{
		{"http://a.com/a/b/c.html", []string{"/", "/a/", "/a/b/", "/a/b/c.html"}, false},
		{"http://a.b/", []string{"/"}, false},
		{"http://a.com/a/b/c/d/e.html?123", []string{"/", "/a/", "/a/b/", "/a/b/c/", "/a/b/c/d/e.html", "/a/b/c/d/e.html?123"}, false},
		{"/asdf", nil, true},
	}

	for i, v := range vectors {
		paths, err := generateLookupPaths(v.url)
		if err != nil != v.fail {
			if err != nil {
				t.Errorf("test %d, unexpected error: %v", i, err)
			} else {
				t.Errorf("test %d, unexpected success", i)
			}
			continue
		}
		if !reflect.DeepEqual(paths, v.output) {
			t.Errorf("test %d, generateLookupPaths(%q) = %q, want %q", i, v.url, paths, v.output)
		}
	}
}

func TestCanonicalURL(t *testing.T) {
	vectors := []struct {
		url    string
		output string
		fail   bool
	}{
		{
			url:    "http://%31%36%38%2e%31%38%38%2e%39%39%2e%32%36/%2E%73%65%63%75%72%65/%77%77%77%2E%65%62%61%79%2E%63%6F%6D/",
			output: "http://168.188.99.26/.secure/www.ebay.com/",
		},
		{
			url:    "http://195.127.0.11/uploads/%20%20%20%20/.verify/.eBaysecure=updateuserdataxplimnbqmn-xplmvalidateinfoswqpcmlx=hgplmcx/",
			output: "http://195.127.0.11/uploads/%20%20%20%20/.verify/.eBaysecure=updateuserdataxplimnbqmn-xplmvalidateinfoswqpcmlx=hgplmcx/",
		},
		{
			url:    "http://host%23.com/%257Ea%2521b%2540c%2523d%2524e%25f%255E00%252611%252A22%252833%252944_55%252B",
			output: "http://host%23.com/~a!b@c%23d$e%25f^00&11*22(33)44_55+",
		},

		{"http://host/%25%32%35", "http://host/%25", false},
		{"http://host/%25%32%35%25%32%35", "http://host/%25%25", false},
		{"http://host/%2525252525252525", "http://host/%25", false},
		{"http://host/asdf%25%32%35asd", "http://host/asdf%25asd", false},
		{"http://host/%%%25%32%35asd%%", "http://host/%25%25%25asd%25%25", false},
		{"http://www.google.com/", "http://www.google.com/", false},
		{"http://3279880203/blah", "http://195.127.0.11/blah", false},
		{"http://www.evil.com/blah#frag", "http://www.evil.com/blah", false},
		{"http://www.GOOgle.com/", "http://www.google.com/", false},
		{"http://www.google.com.../", "http://www.google.com/", false},
		{"http://www.google.com/foo\tbar\rbaz\n2", "http://www.google.com/foobarbaz2", false},
		{"http://www.google.com/q?", "http://www.google.com/q", false},
		{"http://www.google.com/q?r?", "http://www.google.com/q", false},
		{"http://www.google.com/q?r?s", "http://www.google.com/q", false},
		{"http://evil.com/foo#bar#baz", "http://evil.com/foo", false},
		{"http://evil.com/foo;", "http://evil.com/foo;", false},
		{"http://evil.com/foo?bar;", "http://evil.com/foo", false},
		{"http://\x01\x80.com/", "http://%01%80.com/", false},
		{"http://notrailingslash.com", "http://notrailingslash.com/", false},
		{"http://www.gotaport.com:1234/", "http://www.gotaport.com/", false},
		{"  http://www.google.com/  ", "http://www.google.com/", false},
		{"http:// leadingspace.com/", "http://%20leadingspace.com/", false},
		{"http://%20leadingspace.com/", "http://%20leadingspace.com/", false},
		{"%20leadingspace.com/", "http://%20leadingspace.com/", false},
		{"https://www.securesite.com/", "https://www.securesite.com/", false},
		{"ftp://ftp.myfiles.com/", "ftp://ftp.myfiles.com/", false},
		{"http://some%1bhost.com/%1b", "http://some%1bhost.com/%1b", false},
		{"  http://www.google.com/  ", "http://www.google.com/", false},
		{"http://www.google.com/q?r?s%3F", "http://www.google.com/q", false},
		{"http://www.\xC3\xBcmlat.com/", "http://www.xn--mlat-zra.com/", false},
		{"http://[2001:470:1:18::114]/", "http://[2001:470:1:18::114]/", false}, // IPv6 literal.
		{"http%3A%2F%2Fwackyurl.com:80/", "http://wackyurl.com/", false},
		{"http://W!eird<>Ho$^.com/", "http://w!eird<>ho$^.com/", false},
		{"http://i.have.way.too.many.dots.com/", "http://i.have.way.too.many.dots.com/", false},

		// All of these cases are missing a valid hostname and should return empty
		{"", "", true},
		{":", "", true},
		{"/blah", "", true},
		{"#ref", "", true},
		{"/blah#ref", "", true},
		{"?query#ref", "", true},
		{"/blah?query#ref", "", true},
		{"/blah;param", "", true},
		{"http://#ref", "", true},
		{"http:///blah#ref", "", true},
		{"http://?query#ref", "", true},
		{"http:///blah?query#ref", "", true},
		{"http:///blah;param", "", true},
		{"http:///blah;param?query#ref", "", true},
		{"mailto:bryner@google.com", "", true},
	}
	for i, v := range vectors {
		path, err := canonicalURL(v.url)
		if err != nil != v.fail {
			if err != nil {
				t.Errorf("test %d, unexpected error: %v", i, err)
			} else {
				t.Errorf("test %d, unexpected success. URL: %v, got: %v, want: %v", i, v.url, path, v.output)
			}
			continue
		}
		if path != v.output {
			t.Errorf("test %d, canonicalURL(%q) = %q, want %q", i, v.url, path, v.output)
		}
	}
}
