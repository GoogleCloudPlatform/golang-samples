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

package main

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "os"
    "encoding/json"
    "github.com/dgrijalva/jwt-go"
)

var AUDIENCE string
var CERTIFICATES = make(map[string]string)


func main() {
    http.HandleFunc("/", indexHandler)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
        log.Printf("Defaulting to port %s", port)
    }

    log.Printf("Listening on port %s", port)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func certs() map[string]string {
    const url = "https://www.gstatic.com/iap/verify/public_key"

    if len(CERTIFICATES) == 0 {
        resp, err := http.Get(url)
        if err != nil {
            log.Printf("Failed to fetch certificates: %s", err)
            return CERTIFICATES
        }

        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Printf("Error reading certs: %s", err)
            return CERTIFICATES
        }

        err = json.Unmarshal(body, &CERTIFICATES)
        if err != nil {
            log.Printf("Error converting from JSON: %s", err)
            return CERTIFICATES
        }
    }

    return CERTIFICATES
}

func getMetadata(itemName string) string {
    const url = "http://metadata.google.internal/computeMetadata/v1/project/"

    client := &http.Client{}
    req, _ := http.NewRequest("GET", url + itemName, nil)
    req.Header.Add("Metadata-Flavor", "Google")
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error making metadata request: %s", err)
        return "None"
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    if err != nil {
        log.Printf("Error reading metadata: %s", err)
        return "None"
    } else {
        return(string(body))
    }
}

func audience() string {
    if AUDIENCE == "" {
        project_number := getMetadata("numeric-project-id")
        project_id := getMetadata("project-id")
        AUDIENCE = "/projects/" + project_number + "/apps/" + project_id
    }

    return AUDIENCE
}

func validateAssertion(assertion string) (string, string) {
    certificates := certs()

    token, err := jwt.Parse(assertion, func(token *jwt.Token) (interface{}, error) {
        keyId := token.Header["kid"].(string)

        _, ok := token.Method.(*jwt.SigningMethodECDSA)
        if !ok {
            log.Printf("Wrong signing method: %v", token.Header["alg"])
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }

        cert := certificates[keyId]
        return jwt.ParseECPublicKeyFromPEM([]byte(cert))
    })

    if err != nil {
        log.Printf("Failed to validate assertion: %s", assertion)
        return "None", "None"
    }

    claims, _ := token.Claims.(jwt.MapClaims)
    return claims["email"].(string), claims["sub"].(string)
}

// indexHandler responds to requests with our greeting.
func indexHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    assertion := r.Header.Get("X-Goog-IAP-JWT-Assertion")
    email, _ := validateAssertion(assertion)
    fmt.Fprint(w, "Hello " + email)
}
