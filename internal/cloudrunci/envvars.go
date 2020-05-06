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
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// EnvVars is a collection of environment variables.
type EnvVars map[string]string

func (e EnvVars) String() string {
	s := make([]string, len(e))
	i := 0
	for k := range e {
		s[i] = e.Variable(k)
		i++
	}
	sort.Strings(s)
	return strings.Join(s, ",")
}

// Variable retrieves an environment variable assignment as a string.
func (e EnvVars) Variable(k string) string {
	return fmt.Sprintf("%s=%s", strings.TrimSpace(k), strings.TrimSpace(e[k]))
}

// KeyString converts the environment variables names to a comma-delimited list.
func (e EnvVars) KeyString() string {
	s := make([]string, len(e))
	i := 0
	for k := range e {
		s[i] = k
		i++
	}
	sort.Strings(s)
	return strings.Join(s, ",")
}

// keyRegex defines valid environment variable name.
// Regex arranged so the first submatch will have the name without whitespace.
// https://stackoverflow.com/questions/2821043/allowed-characters-in-linux-environment-variable-names
var keyRegex = regexp.MustCompile(`^\s*([a-zA-Z_]+\w*)\s*$`)

// Validate confirms all environment variables are valid.
func (e EnvVars) Validate() error {
	broken := EnvVars{}
	for k := range e {
		match := keyRegex.FindStringSubmatch(k)
		if len(match) == 0 {
			broken[k] = ""
		}
	}
	if len(broken) > 0 {
		return fmt.Errorf("invalid environment variable names: %s", broken.KeyString())
	}

	return nil
}
