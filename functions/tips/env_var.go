// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START functions_env_vars]

// Package tips contains tips for writing Cloud Functions in Go.
package tips

import (
	"fmt"
	"net/http"
	"os"
)

// EnvVar is an example of getting an environment variable in a Go function.
func EnvVar(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "FOO: %q", os.Getenv("FOO"))
}

// [END functions_env_vars]
