package main

import (
	"crypto/rsa"
	"crypto/x509"
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

	serviceAccount = flag.String("service-account", "", "Path to service account JSON file. Cannot be used with -oauth-config.")
	oauthConfig    = flag.String("oauth-config", "", "Path to oauth client secrets JSON file. Cannot be used with -service-account.")
)

func main() {
	flag.Parse()

	if *apiKey == "" || *host == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *serviceAccount == "" && *oauthConfig == "" {
		fmt.Fprintln(os.Stderr, "Provide one of -oauth-config and -service-account")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *serviceAccount != "" && *oauthConfig != "" {
		fmt.Fprintln(os.Stderr, "Can't use -oauth-config and -service-account")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var resp *http.Response
	var err error
	if *serviceAccount != "" {
		resp, err = doJWT()
	} else {
		log.Fatal("Not implemented.")
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

func doJWT() (*http.Response, error) {
	sa, err := ioutil.ReadFile(*serviceAccount)
	if err != nil {
		log.Fatalf("Could not read service account file: %v", err)
	}
	conf, err := google.JWTConfigFromJSON(sa)
	if err != nil {
		log.Fatalf("Could not parse service account JSON: %v", err)
	}
	rsaKey, err := parseKey(conf.PrivateKey)
	if err != nil {
		log.Fatalf("Could not get RSA key: %v", err)
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
		log.Fatalf("Could not encode JWT: %v", err)
	}

	req, _ := http.NewRequest("GET", *host+"/auth/info/googlejwt?key="+*apiKey, nil)
	req.Header.Add("Authorization", "Bearer "+msg)
	return http.DefaultClient.Do(req)
}

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
