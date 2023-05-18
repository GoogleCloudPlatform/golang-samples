// Copyright 2022 Google LLC
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
	"encoding/json"
	"fmt"

	"google.golang.org/api/kgsearch/v1"
)

type ItemList struct {
	Items []ItemListElement `json:"itemListElement"`
}
type ItemListElement struct {
	Result EntitySearchResult `json:"result"`
}
type EntitySearchResult struct {
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	DetailedDesc DetailedDesc `json:"detailedDescription"`
	URL          string       `json:"url"`
	Image        Image
}
type DetailedDesc struct {
	ArticleBody string
	URL         string
}
type Image struct {
	ContentURL string
}

func formatSlackMessage(query string, response *kgsearch.SearchResponse) (*Message, error) {
	if response == nil {
		return nil, fmt.Errorf("empty response")
	}

	if response.ItemListElement == nil || len(response.ItemListElement) == 0 {
		message := &Message{
			ResponseType: "in_channel",
			Text:         fmt.Sprintf("Query: %s", query),
			Attachments: []Attachment{
				{
					Color: "#d6334b",
					Text:  "No results match your query.",
				},
			},
		}
		return message, nil
	}

	// The KnowledgeGraph API returns an empty interface. To make this more
	// useful, we convert it back to json, and unmarshal into specific types.
	jsonstring, err := response.MarshalJSON()
	if err != nil {
		return nil, err
	}
	r := &ItemList{}
	if err := json.Unmarshal(jsonstring, r); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}
	result := r.Items[0].Result

	attach := Attachment{Color: "#3367d6"}
	attach.Title = result.Name
	attach.TitleLink = result.DetailedDesc.URL
	attach.Text = result.DetailedDesc.ArticleBody
	attach.ImageURL = result.Image.ContentURL

	message := &Message{
		ResponseType: "in_channel",
		Text:         fmt.Sprintf("Query: %s", query),
		Attachments:  []Attachment{attach},
	}
	return message, nil
}

// [END functions_slack_format]
