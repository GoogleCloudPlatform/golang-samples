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

package cloudrunci

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// gcloudBin is the path to the gcloud executable.
var gcloudBin string

func init() {
	gcloudBin = os.Getenv("GCLOUD_BIN")
	if gcloudBin == "" {
		gcloudBin = "gcloud"
	}
}

// gcloud provides a common mechanism for executing gcloud commands to handle output and errors.
func gcloud(label string, cmd *exec.Cmd) ([]byte, error) {
	log.Printf("Running %s...", label)
	log.Println("Executing:", cmd.Path, strings.Join(cmd.Args[1:], " "))
	// TODO: add a flag for verbose output (e.g. when running with binary created with `go test -c`)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error Output for %s:", label)
		os.Stderr.Write(out)
		return []byte{}, fmt.Errorf("%s: %q", label, err)
	}

	return bytes.TrimSpace(out), err
}

// CreateIDToken generates an ID token for requests to the fully managed platform.
// In the future the URL of the targeted service will be used to scope the audience.
func CreateIDToken(_ string) (string, error) {
	args := []string{
		"--quiet",
		"auth",
		"print-identity-token",
	}

	cmd := exec.Command(gcloudBin, args...)

	out, err := gcloud("operation [id-token]", cmd)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
