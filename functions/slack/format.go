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

// [START functions_slack_format]

package slack

import (
	"fmt"

	"google.golang.org/api/kgsearch/v1"
)

func formatSlackMessage(query string, response *kgsearch.SearchResponse) (*Message, error) {
	if response == nil {
		return nil, fmt.Errorf("empty response")
	}

	if response.ItemListElement == nil || len(response.ItemListElement) == 0 {
		message := &Message{
			ResponseType: "in_channel",
			Text:         fmt.Sprintf("Query: %s", query),
			Attachments: []attachment{
				{
					Color: "#d6334b",
					Text:  "No results match your query.",
				},
			},
		}
		return message, nil
	}

	entity, ok := response.ItemListElement[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("could not parse response entity")
	}
	result, ok := entity["result"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("error formatting response result")
	}

	attach := attachment{Color: "#3367d6"}
	if name, ok := result["name"].(string); ok {
		if description, ok := result["description"].(string); ok {
			attach.Title = fmt.Sprintf("%s: %s", name, description)
		} else {
			attach.Title = name
		}
	}
	if detailedDesc, ok := result["detailedDescription"].(map[string]interface{}); ok {
		if url, ok := detailedDesc["url"].(string); ok {
			attach.TitleLink = url
		}
		if article, ok := detailedDesc["articleBody"].(string); ok {
			attach.Text = article
		}
	}
	if image, ok := result["image"].(map[string]interface{}); ok {
		if imageURL, ok := image["contentUrl"].(string); ok {
			attach.ImageURL = imageURL
		}
	}

	message := &Message{
		ResponseType: "in_channel",
		Text:         fmt.Sprintf("Query: %s", query),
		Attachments:  []attachment{attach},
	}
	return message, nil
}

// [END functions_slack_format]
