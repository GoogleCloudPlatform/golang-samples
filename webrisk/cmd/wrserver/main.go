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

// Command wrserver is an application for serving URL lookups via a simple API.
//
// In order to abstract away the complexities of the Web Risk API v4, the
// wrserver application can be used to serve a subset of API v4 over HTTP.
// This subset is intentionally small so that it would be easy to implement by
// a client. It is intended for wrserver to either be running locally on a
// client's machine or within the same local network. That way it can handle
// most local API calls before resorting to making an API call to the actual
// Web Risk API over the internet.
//
// Usage of wrserver looks something like this:
//	             _________________
//	            |                 |
//	            |  Web Risk  |
//	            |  API v4 servers |
//	            |_________________|
//	                     |
//	            ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
//	               The Internet
//	            ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
//	                     |
//	              _______V_______
//	             |               |
//	             |   SBServer    |
//	      +------|  Application  |------+
//	      |      |_______________|      |
//	      |              |              |
//	 _____V_____    _____V_____    _____V_____
//	|           |  |           |  |           |
//	|  Client1  |  |  Client2  |  |  Client3  |
//	|___________|  |___________|  |___________|
//
// In theory, each client could directly use the Go WebriskClient implementation,
// but there are situations where that is not desirable. For example, the client
// may not be using the language Go, or there may be multiple clients in the
// same machine or local network that would like to share a local database
// and cache. The wrserver was designed to address these issues.
//
// The wrserver application is technically a proxy since it is itself actually
// an API v4 client. It connects the Web Risk API servers using an API key
// and maintains a local database and cache. However, it is also a server since
// it re-serves a subset of the API v4 endpoints. These endpoints are minimal
// in that they do not require each client to maintain state between calls.
//
// The assumption is that communication between SBServer and Client1, Client2,
// and Client3 is inexpensive, since they are within the same machine or same
// local network. Thus, the wrserver can efficiently satisfy some requests
// without talking to the global Web Risk servers since it has a
// potentially larger cache. Furthermore, it can multiplex multiple requests to
// the Web Risk servers on fewer TCP connections, reducing the cost for
// comparatively more expensive internet transfers.
//
// By default, the wrserver listens on localhost:8080 and serves the following
// API endpoints:
//	/v4/threatMatches:find
//	/v4/threatLists
//	/status
//	/r
//
//
// Endpoint: /v4/threatMatches:find
//
// This is a lightweight implementation of the API v4 threatMatches endpoint.
// Essentially, it takes in a list of URLs, and returns a list of threat matches
// for those URLs. Unlike the Web Risk API, it does not require an API key.
//
// Example usage:
//	# Send request to server:
//	$ curl \
//	  -H "Content-Type: application/json" \
//	  -X POST -d '{
//	      "threatInfo": {
//	          "threatTypes":      ["UNWANTED_SOFTWARE", "MALWARE"],
//	          "platformTypes":    ["ANY_PLATFORM"],
//	          "threatEntryTypes": ["URL"],
//	          "threatEntries": [
//	              {"url": "google.com"},
//	              {"url": "bad1url.org"},
//	              {"url": "bad2url.org"}
//	          ]
//	      }
//	  }' \
//	  localhost:8080/v4/threatMatches:find
//
//	# Receive response from server:
//	{
//	    "matches": [{
//	        "threat":          {"url": "bad1url.org"},
//	        "platformType":    "ANY_PLATFORM",
//	        "threatType":      "UNWANTED_SOFTWARE",
//	        "threatEntryType": "URL"
//	    }, {
//	        "threat":          {"url": "bad2url.org"},
//	        "platformType":    "ANY_PLATFORM",
//	        "threatType":      "UNWANTED_SOFTWARE",
//	        "threatEntryType": "URL"
//	    }, {
//	        "threat":          {"url": "bad2url.org"},
//	        "platformType":    "ANY_PLATFORM",
//	        "threatType":      "MALWARE",
//	        "threatEntryType": "URL"
//	    }]
//	}
//
//
// Endpoint: /v4/threatLists
//
// The endpoint returns a list of the threat lists that the wrserver is
// currently subscribed to. The threats returned by the earlier threatMatches
// API call may only be one of these types.
//
// Example usage:
//	# Send request to server:
//	$ curl -X GET localhost:8080/v4/threatLists
//
//	# Receive response from server:
//	{
//	    "threatLists": [{
//	        "threatType":      "MALWARE"
//	        "platformType":    "ANY_PLATFORM",
//	        "threatEntryType": "URL",
//	    }, {
//	        "threatType":      "SOCIAL_ENGINEERING",
//	        "platformType":    "ANY_PLATFORM"
//	        "threatEntryType": "URL",
//	    }, {
//	        "threatType":      "UNWANTED_SOFTWARE"
//	        "platformType":    "ANY_PLATFORM",
//	        "threatEntryType": "URL",
//	    }]
//	}
//
//
// Endpoint: /status
//
// The status endpoint allows a client to obtain some statistical information
// regarding the health of wrserver. It can be used to determine how many
// requests were satisfied locally by wrserver alone and how many requests
// were forwarded to the Web Risk API servers.
//
// Example usage:
//	$ curl localhost:8080/status
//	{
//	    "Stats" : {
//	        "QueriesByDatabase" : 132,
//	        "QueriesByCache" : 31,
//	        "QueriesByAPI" : 6,
//	        "QueriesFail" : 0,
//	    },
//	    "Error" : ""
//	}
//
//
// Endpoint: /r
//
// The redirector endpoint allows a client to pass in a query URL.
// If the URL is safe, the client is automatically redirected to the target.
// If the URL is unsafe, then an interstitial warning page is shown instead.
//
// Example usage:
//	$ curl -i localhost:8080/r?url=http://google.com
//	HTTP/1.1 302 Found
//	Location: http://google.com
//
//	$ curl -i localhost:8080/r?url=http://bad1url.org
//	HTTP/1.1 200 OK
//	Date: Wed, 13 Apr 2016 21:29:33 GMT
//	Content-Length: 1783
//	Content-Type: text/html; charset=utf-8
//
//	<!-- Warning interstitial page shown -->
//	...
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/GoogleCloudPlatform/golang-samples/webrisk"
	pb "github.com/GoogleCloudPlatform/golang-samples/webrisk/internal/webrisk_proto"

	_ "github.com/GoogleCloudPlatform/golang-samples/webrisk/cmd/wrserver/statik"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/rakyll/statik/fs"
)

const (
	statusPath     = "/status"
	findThreatPath = "/v1beta1/uris:search"
	redirectPath   = "/r"
)

const (
	mimeJSON  = "application/json"
	mimeProto = "application/x-protobuf"
)

var (
	apiKeyFlag   = flag.String("apikey", "", "specify your Web Risk API key")
	srvAddrFlag  = flag.String("srvaddr", "localhost:8080", "TCP network address the HTTP server should use")
	proxyFlag    = flag.String("proxy", "", "proxy to use to connect to the HTTP server")
	databaseFlag = flag.String("db", "", "path to the Web Risk database.")
)

var threatTemplate = map[webrisk.ThreatType]string{
	webrisk.ThreatTypeMalware:           "/malware.tmpl",
	webrisk.ThreatTypeUnwantedSoftware:  "/unwanted.tmpl",
	webrisk.ThreatTypeSocialEngineering: "/social_engineering.tmpl",
}

const usage = `wrserver: starts a Web Risk API proxy server.

In order to abstract away the complexities of the Web Risk API v4, the
wrserver application can be used to serve a subset of the v4 API.
This subset is intentionally small so that it would be easy to implement by
a client. It is intended for wrserver to either be running locally on a
client's machine or within the same local network so that it can handle most
local API calls before resorting to making an API call to the actual
Web Risk API over the internet.

Usage: %s -apikey=$APIKEY

`

// unmarshal reads pbResp from req. The mime will either be JSON or ProtoBuf.
func unmarshal(req *http.Request, pbReq proto.Message) (string, error) {
	var mime string
	alt := req.URL.Query().Get("alt")
	if alt == "" {
		alt = req.Header.Get("Content-Type")
	}
	switch alt {
	case "json", mimeJSON:
		mime = mimeJSON
	case "proto", mimeProto:
		mime = mimeProto
	default:
		return mime, errors.New("invalid interchange format")
	}

	switch req.Header.Get("Content-Type") {
	case mimeJSON:
		if err := jsonpb.Unmarshal(req.Body, pbReq); err != nil {
			return mime, err
		}
	case mimeProto:
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return mime, err
		}
		if err := proto.Unmarshal(body, pbReq); err != nil {
			return mime, err
		}
	}
	return mime, nil
}

// marshal writes pbResp into resp. The mime can either be JSON or ProtoBuf.
func marshal(resp http.ResponseWriter, pbResp proto.Message, mime string) error {
	resp.Header().Set("Content-Type", mime)
	switch mime {
	case mimeProto:
		body, err := proto.Marshal(pbResp)
		if err != nil {
			return err
		}
		if _, err := resp.Write(body); err != nil {
			return err
		}
	case mimeJSON:
		var m jsonpb.Marshaler
		var b bytes.Buffer
		if err := m.Marshal(&b, pbResp); err != nil {
			return err
		}
		if _, err := resp.Write(b.Bytes()); err != nil {
			return err
		}
	default:
		return errors.New("invalid interchange format")
	}
	return nil
}

// serveStatus writes a simple JSON with server status information to resp.
func serveStatus(resp http.ResponseWriter, req *http.Request, sb *webrisk.WebriskClient) {
	stats, sbErr := sb.Status()
	errStr := ""
	if sbErr != nil {
		errStr = sbErr.Error()
	}
	buf, err := json.Marshal(struct {
		Stats webrisk.Stats
		Error string
	}{stats, errStr})
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	resp.Header().Set("Content-Type", mimeJSON)
	resp.Write(buf)
}

// serveLookups is a light-weight implementation of the "/v4/threatMatches:find"
// API endpoint. This allows clients to look up whether a given URL is safe.
// Unlike the official API, it does not require an API key.
// It supports both JSON and ProtoBuf.
func serveLookups(resp http.ResponseWriter, req *http.Request, sb *webrisk.WebriskClient) {
	if req.Method != "POST" {
		http.Error(resp, "invalid method", http.StatusBadRequest)
		return
	}

	// Decode the request message.
	pbReq := new(pb.SearchUrisRequest)
	mime, err := unmarshal(req, pbReq)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Should this handler use the information in threatTypes,
	// platformTypes, and threatEntryTypes?

	// Parse the request message.
	urls := []string{pbReq.Uri}

	// Lookup the URL.
	utss, err := sb.LookupURLsContext(req.Context(), urls)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	// Compose the response message.
	pbResp := &pb.SearchUrisResponse{
		Threat: &pb.SearchUrisResponse_ThreatUri{},
	}
	for _, uts := range utss {
		// Use map to condense duplicate ThreatDescriptor entries.
		tdm := make(map[webrisk.ThreatType]bool)
		for _, ut := range uts {
			tdm[ut.ThreatType] = true
		}

		for td := range tdm {
			pbResp.Threat.ThreatTypes = append(pbResp.Threat.ThreatTypes, pb.ThreatType(td))
		}
	}

	// Encode the response message.
	if err := marshal(resp, pbResp, mime); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
}

func parseTemplates(fs http.FileSystem, t *template.Template, paths ...string) (*template.Template, error) {
	for _, path := range paths {
		file, err := fs.Open(path)
		if err != nil {
			return nil, err
		}
		tmpl, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		t, err = t.Parse(string(tmpl))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

// serveRedirector implements a basic HTTP redirector that will filter out
// redirect URLs that are unsafe according to the Web Risk API.
func serveRedirector(resp http.ResponseWriter, req *http.Request, sb *webrisk.WebriskClient, fs http.FileSystem) {
	rawURL := req.URL.Query().Get("url")
	if rawURL == "" || req.URL.Path != "/r" {
		http.NotFound(resp, req)
		return
	}
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	threats, err := sb.LookupURLsContext(req.Context(), []string{rawURL})
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(threats[0]) == 0 {
		http.Redirect(resp, req, rawURL, http.StatusFound)
		return
	}

	t := template.New("Web Risk Interstitial")
	for _, threat := range threats[0] {
		if tmpl, ok := threatTemplate[threat.ThreatType]; ok {
			t, err = parseTemplates(fs, t, tmpl, "/interstitial.html")
			if err != nil {
				http.Error(resp, err.Error(), http.StatusInternalServerError)
				return
			}
			err = t.Execute(resp, map[string]interface{}{
				"Threat": threat,
				"Url":    parsedURL})
			if err != nil {
				http.Error(resp, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}
	http.Error(resp, err.Error(), http.StatusInternalServerError)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if *apiKeyFlag == "" {
		fmt.Fprintln(os.Stderr, "No -apikey specified")
		os.Exit(1)
	}
	conf := webrisk.Config{
		APIKey:   *apiKeyFlag,
		ProxyURL: *proxyFlag,
		DBPath:   *databaseFlag,
		Logger:   os.Stderr,
	}
	sb, err := webrisk.NewWebriskClient(conf)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to initialize Web Risk client: ", err)
		os.Exit(1)
	}
	statikFS, err := fs.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to initialize static files: ", err)
		os.Exit(1)
	}

	http.HandleFunc(statusPath, func(w http.ResponseWriter, r *http.Request) {
		serveStatus(w, r, sb)
	})
	http.HandleFunc(findThreatPath, func(w http.ResponseWriter, r *http.Request) {
		serveLookups(w, r, sb)
	})
	http.HandleFunc(redirectPath, func(w http.ResponseWriter, r *http.Request) {
		serveRedirector(w, r, sb, statikFS)
	})
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(statikFS)))

	fmt.Fprintln(os.Stdout, "Starting server at", *srvAddrFlag)
	if err := http.ListenAndServe(*srvAddrFlag, nil); err != nil {
		fmt.Fprintln(os.Stderr, "Server error:", err)
		return
	}
	fmt.Fprintln(os.Stdout, "Stopping server")
}
