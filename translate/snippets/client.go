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

// Package trsnip contains snippet code for the Translate API.
// The code is not runnable.
package trsnip

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

func createClientWithKey() {
	ctx := context.Background()

	const apiKey = "YOUR_TRANSLATE_API_KEY"
	client, err := translate.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Translate(ctx, []string{"Hello, world!"}, language.Russian, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v", resp)
}
