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
