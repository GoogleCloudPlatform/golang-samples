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

// Package leaderboard starts a Gopher Run leaderboard server.
package leaderboard

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// ScoreData is a player's score.
type ScoreData struct {
	Name     string  `json:"name"`
	Team     string  `json:"team"`
	Coins    int     `json:"coins"`
	Distance float32 `json:"distance"`
	Combo    float32 `json:"combo"`
}

// TopScores returns the top 10 scores in the leaderboard.
func TopScores(ctx context.Context, client *firestore.Client) ([]ScoreData, error) {
	iter := client.Collection("leaderboard").Query.OrderBy("coins", firestore.Desc).Limit(10).Documents(ctx)
	var top []ScoreData
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("iter.Next: %v", err)
		}
		var d ScoreData
		if err = doc.DataTo(&d); err != nil {
			return nil, fmt.Errorf("doc.DataTo: %v", err)
		}
		top = append(top, d)
	}
	return top, nil
}

// AddScore adds a score to the leaderboard database and returns information about whether it updated an existing score.
func AddScore(ctx context.Context, client *firestore.Client, d ScoreData) (string, error) {
	var oldD ScoreData
	iter := client.Collection("leaderboard").Query.Limit(1).Where("name", "==", d.Name).Documents(ctx)
	doc, err := iter.Next()
	if err != iterator.Done && err != nil {
		return "", fmt.Errorf("iter.Next: %v", err)
	}
	if err != iterator.Done {
		if err = doc.DataTo(&oldD); err != nil {
			return "", fmt.Errorf("doc.DataTo: %v", err)
		}
	}
	s := ""
	if oldD.Coins < d.Coins {
		s = "pb"
		_, err := client.Collection("leaderboard").Doc(d.Name).Set(ctx, map[string]interface{}{
			"name":     d.Name,
			"team":     d.Team,
			"coins":    d.Coins,
			"distance": d.Distance,
			"combo":    d.Combo,
		})
		if err != nil {
			return "", fmt.Errorf("Doc(%v).Set: %v", d.Name, err)
		}
	}
	return s, nil
}
