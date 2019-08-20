// Copyright 2018, Google, LLC.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// [START functions_verify_webhook]

package slack

import (
	"fmt"
	"net/url"
)

func verifyWebhook(form url.Values) error {
	t, ok := form["token"]
	if !ok || len(t) == 0 {
		return fmt.Errorf("empty form token")
	}
	if t[0] != config.Token {
		return fmt.Errorf("invalid request/credentials")
	}
	return nil
}

// [END functions_verify_webhook]
