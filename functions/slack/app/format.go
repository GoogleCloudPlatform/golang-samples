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

package slack

// [START functions_slack_format]
import (
	"fmt"

	"google.golang.org/api/kgsearch/v1"
)

func formatSlackMessage(query string, response kgsearch.SearchResponse) (*SlackMessage, error) {
	entity := response.ItemListElement[0].result
	message := &SlackMessage{
		responseType: "in_channel",
		text:         fmt.Sprintf("Query: %s", query),
		attachments:  []attachment{},
	}

	attachment := attachment{}
	if entity != "" {
		name := entity.Name
		description := entity.Description
		detailedDesc := entity.DetailedDescription
		url := detailedDesc.URL
		article := detailedDesc.articleBody
		imageURL := entity.Image.ContentUrl

		attachment.color = "//3367d6"
		if name && description {
			attachment.title = fmt.Sprintf("%s: %s", entity.Name, entity.Description)
		} else if name {
			attachment.title = name
		}
		if url {
			attachment.titleLink = url
		}
		if article {
			attachment.text = article
		}
		if imageURL {
			attachment.imageURL = imageURL
		}
	} else {
		attachment.text = "No results match your query."
	}
	message.attachments = append(message.attachments, attachment)

	return message, nil
}

// [END functions_slack_format]
