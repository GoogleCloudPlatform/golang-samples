// Copyright 2026 Google LLC
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

package developerknowledge

// [START developerknowledge_answer_query]
import (
	"context"
	"fmt"
	"io"

	developerknowledge "google.golang.org/api/developerknowledge/v1"
)

// answerQuery answers a developer question grounded in Google developer documentation.
func answerQuery(w io.Writer, query string) (*developerknowledge.AnswerQueryResponse, error) {
	ctx := context.Background()

	svc, err := developerknowledge.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("developerknowledge.NewService: %w", err)
	}

	req := &developerknowledge.AnswerQueryRequest{
		Query: query,
	}

	resp, err := svc.V1.AnswerQuery(req).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("AnswerQuery: %w", err)
	}

	if resp.Answer != nil {
		fmt.Fprintf(w, "Answer:\n%s\n\n", resp.Answer.AnswerText)
		fmt.Fprintf(w, "Citations count: %d\n", len(resp.Answer.Citations))
		fmt.Fprintf(w, "References count: %d\n", len(resp.Answer.References))
	}

	return resp, nil
}

// [END developerknowledge_answer_query]
