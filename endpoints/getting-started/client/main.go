// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command client performs authenticated requests against an Endpoints API server.
package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jws"
)

var (
	host   = flag.String("host", "", "The API host. Required.")
	audience = flag.String("audience", "", "The audience for the JWT. equired")
	serviceAccountFile = flag.String("service-account-file", "", "Path to service account JSON file. Required.")
	serviceAccountEmail = flag.String("service-account-email", "", "Path email associated with the service account. Required.")
)

func main() {
	flag.Parse()

	if *audience == "" || *host == "" || *serviceAccountFile == "" || *serviceAccountEmail == "" {
		fmt.Println("requires: --host, --audience, --service-account-file, --service-account-email")
		os.Exit(1)
	}
	jwt, err := generateJWT(*serviceAccountFile, *serviceAccountEmail, *audience, 3600)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(jwt)
	resp, err := makeJWTRequest(jwt, *host)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}

// [START endpoints_generate_jwt_sa]
func generateJWT(saKeyfile, saEmail, audience string, expiryLength int64) (string, error) {
	now := time.Now().Unix()
	
	// build payload
	jwt := &jws.ClaimSet{
		Iat:           now,
		// expires after 'expiraryLength' seconds.
		Exp:           now + expiryLength,
		// Iss must match 'issuer' in the security configuration in your
		// swagger spec (e.g. service account email). It can be any string.
		Iss:           saEmail,
		// Aud must be either your Endpoints service name, or match the value
		// specified as the 'x-google-audience' in the OpenAPI document.
		Aud:           audience,
		// Sub and Email should match the service account's email address
		Sub:           saEmail,
		PrivateClaims: map[string]interface{}{"email": saEmail},
	}
	jwsHeader := &jws.Header{
		Algorithm: "RS256",
		Typ:       "JWT",
	}

	// sign with keyfile
	sa, err := ioutil.ReadFile(saKeyfile)
	if err != nil {
		return "", fmt.Errorf("Could not read service account file: %v", err)
	}
	conf, err := google.JWTConfigFromJSON(sa)
	if err != nil {
		return "", fmt.Errorf("Could not parse service account JSON: %v", err)
	}
	block, _ := pem.Decode(conf.PrivateKey)
	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("private key parse error: %v", err)
	}
	rsaKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("private key failed rsa.PrivateKey type assertion")
	}
	return jws.Encode(jwsHeader, jwt, rsaKey)
}
// [END endpoints_generate_jwt_sa]

// [START endpoints_jwt_request]
func makeJWTRequest(signedJWT, url string) (string, error){
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error building HTTP request: %v", err)
	}
	req.Header.Add("Authorization", "Bearer " + signedJWT)
	req.Header.Add("content-type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		//fix these error messages
		return "", fmt.Errorf("error making HTTP request", err)
	}
	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error parsing HTTP response", err)
	}
	return string(responseData), nil
}
// [END endpoints_jwt_request]
