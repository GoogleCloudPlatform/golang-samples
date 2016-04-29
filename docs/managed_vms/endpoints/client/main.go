// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command client performs authenticated requests against an Endpoints API server.
package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jws"
)

var (
	host   = flag.String("host", "", "The API host. Required.")
	apiKey = flag.String("api-key", "", "Your API key. Required.")

	echo           = flag.String("echo", "", "Message to echo. Cannot be used with -service-account")
	serviceAccount = flag.String("service-account", "", "Path to service account JSON file. Cannot be used with -echo.")
)

func main() {
	flag.Parse()

	if *apiKey == "" || *host == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *serviceAccount == "" && *echo == "" {
		fmt.Fprint(os.Stderr, "Provide one of -echo or -service-account.")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *serviceAccount != "" && *echo != "" {
		fmt.Fprint(os.Stderr, "Provide only one of -echo or -service-account.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var resp *http.Response
	var err error
	if *echo != "" {
		resp, err = doEcho()
	} else if *serviceAccount != "" {
		resp, err = doJWT()
	}
	if err != nil {
		log.Fatal(err)
	}
	b, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(b)
}

// doEcho performs an authenticated echo request using an API key.
func doEcho() (*http.Response, error) {
	msg := map[string]string{
		"message": *echo,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(msg); err != nil {
		return nil, err
	}
	return http.Post(*host+"/echo?key="+*apiKey, "application/json", &buf)
}

// doJWT performs an authenticated request using the credentials in the service account file.
func doJWT() (*http.Response, error) {
	sa, err := ioutil.ReadFile(*serviceAccount)
	if err != nil {
		return nil, fmt.Errorf("Could not read service account file: %v", err)
	}
	conf, err := google.JWTConfigFromJSON(sa)
	if err != nil {
		return nil, fmt.Errorf("Could not parse service account JSON: %v", err)
	}
	rsaKey, err := parseKey(conf.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Could not get RSA key: %v", err)
	}

	iat := time.Now()
	exp := iat.Add(time.Hour)

	jwt := &jws.ClaimSet{
		Iss:   "jwt-client.endpoints.sample.google.com",
		Sub:   "foo!",
		Aud:   "echo.endpoints.sample.google.com",
		Scope: "email",
		Iat:   iat.Unix(),
		Exp:   exp.Unix(),
	}
	jwsHeader := &jws.Header{
		Algorithm: "RS256",
		Typ:       "JWT",
	}

	msg, err := jws.Encode(jwsHeader, jwt, rsaKey)
	if err != nil {
		return nil, fmt.Errorf("Could not encode JWT: %v", err)
	}

	req, _ := http.NewRequest("GET", *host+"/auth/info/googlejwt?key="+*apiKey, nil)
	req.Header.Add("Authorization", "Bearer "+msg)
	return http.DefaultClient.Do(req)
}

// The following code is copied from golang.org/x/oauth2/internal
// Copyright (c) 2009 The oauth2 Authors. All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:

//   * Redistributions of source code must retain the above copyright
//notice, this list of conditions and the following disclaimer.
//   * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// parseKey converts the binary contents of a private key file
// to an *rsa.PrivateKey. It detects whether the private key is in a
// PEM container or not. If so, it extracts the the private key
// from PEM container before conversion. It only supports PEM
// containers with no passphrase.
func parseKey(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block != nil {
		key = block.Bytes
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(key)
	if err != nil {
		parsedKey, err = x509.ParsePKCS1PrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("private key should be a PEM or plain PKSC1 or PKCS8; parse error: %v", err)
		}
	}
	parsed, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is invalid")
	}
	return parsed, nil
}
