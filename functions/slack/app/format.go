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

// [START functions_slack_format]

package slack

import (
	"fmt"

	"google.golang.org/api/kgsearch/v1"
)

func formatSlackMessage(query string, response *kgsearch.SearchResponse) (*Message, error) {
	var entity interface{}
	if len(response.ItemListElement) > 0 {
		entity = response.ItemListElement[0]
	}
	resp, ok := entity.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("error formatting response entity")
	}
	result, ok := resp["result"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("error formatting response result")
	}
	message := &Message{
		ResponseType: "in_channel",
		Text:         fmt.Sprintf("Query: %s", query),
		Attachments:  []attachment{},
	}

	attachment := attachment{}
	if entity != nil {
		if name, ok := result["name"].(string); ok {
			if description, ok := result["description"].(string); ok {
				attachment.Title = fmt.Sprintf("%s: %s", name, description)
			} else {
				attachment.Title = name
			}
		}
		if detailedDesc, ok := result["detailedDescription"].(map[string]interface{}); ok {
			if url, ok := detailedDesc["url"].(string); ok {
				attachment.TitleLink = url
			}
			if article, ok := detailedDesc["articleBody"].(string); ok {
				attachment.Text = article
			}
		}
		if image, ok := result["image"].(map[string]interface{}); ok {
			if imageURL, ok := image["contentUrl"].(string); ok {
				attachment.ImageURL = imageURL
			}
		}
		attachment.Color = "//3367d6"
	} else {
		attachment.Text = "No results match your query."
	}
	message.Attachments = append(message.Attachments, attachment)

	return message, nil
}

// [END functions_slack_format]
