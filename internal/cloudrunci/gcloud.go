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
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

// gcloudBin is the path to the gcloud executable.
var gcloudBin string

func init() {
	gcloudBin = os.Getenv("GCLOUD_BIN")
	if gcloudBin == "" {
		gcloudBin = "gcloud"
	}
}

// gcloud provides a common mechanism for executing gcloud commands.
// It will attempt to retry failed commands. Use gcloudWithoutRetry() for no retry.
func gcloud(label string, cmd *exec.Cmd) ([]byte, error) {
	var out []byte
	var err error

	delaySeconds := 2 * time.Second
	if strings.Contains(label, labelOperationBuild) {
		delaySeconds = 60 * time.Second
	}

	maxAttempts := 5
	success := testutil.RetryWithoutTest(maxAttempts, delaySeconds, func(r *testutil.R) {
		// exec.Cmd objects cannot be reused once started, so first make a copy.
		cmdCopy := &exec.Cmd{
			Path: cmd.Path,
			Args: cmd.Args,
			Env:  cmd.Env,
			Dir:  cmd.Dir,
		}
		out, err = gcloudExec(fmt.Sprintf("Attempt #%d: ", r.Attempt), label, cmdCopy)
		if err != nil {
			log.Printf("gcloudExec: %v", err)
			r.Fail()
		}
		// Reset stdout for retry.
		cmd.Stdout = nil
		cmd.Stderr = nil
	})

	if success {
		return out, nil
	}

	return out, fmt.Errorf("gcloudExec: %s: gave up after %d failed attempts", label, maxAttempts)
}

// gcloudWithoutRetry provides a common mechanism for executing gcloud commands.
func gcloudWithoutRetry(label string, cmd *exec.Cmd) ([]byte, error) {
	return gcloudExec("", label, cmd)
}

// gcloudExec adds output prefixing to the execution of the provided command.
func gcloudExec(prefix string, label string, cmd *exec.Cmd) ([]byte, error) {
	log.Printf("%sRunning: %s...", prefix, label)
	log.Printf("%sExecuting: %s: %s: %s", prefix, label, cmd.Path, strings.Join(cmd.Args[1:], " "))
	// TODO: add a flag for verbose output (e.g. when running with binary created with `go test -c`)
	out, err := cmd.CombinedOutput()
	if err != nil {
		os.Stderr.Write([]byte(fmt.Sprintf("%s%s: Error Output\n###\n", prefix, label)))
		if len(out) > 0 {
			os.Stderr.Write(out)
		} else {
			os.Stderr.Write([]byte("no output produced"))
		}
		os.Stderr.Write([]byte("\n###\n"))
		return out, fmt.Errorf("%s%s: %q", prefix, label, err)
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
