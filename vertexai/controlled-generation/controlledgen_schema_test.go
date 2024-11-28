// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controlledgeneration

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func Test_controlledGenerationResponseSchema(t *testing.T) {
	tc := testutil.SystemTest(t)
	w := new(bytes.Buffer)

	location := "us-central1"
	modelName := "gemini-1.5-pro-001"

	err := controlledGenerationResponseSchema(w, tc.ProjectID, location, modelName)
	if err != nil {
		t.Fatalf("controlledGenerationResponseSchema: %v", err.Error())
	}

	// We explicitly requested a response in JSON with a specific schema, so we're
	// expecting to properly decode the output as a slice of structured values.
	type Recipe struct {
		Name string `json:"recipe_name"`
	}
	var recipes []Recipe
	err = json.Unmarshal(w.Bytes(), &recipes)
	if err != nil {
		t.Errorf(`could not unmarshal response:
%s

into slice of Recipe because: %v`, w.Bytes(), err)
	}
	if len(recipes) == 0 {
		t.Errorf("zero recipes generated")
	}
	for i, recipe := range recipes {
		if recipe.Name == "" {
			t.Errorf("recipes[%d] doesn't have a name", i)
		}
	}
}

func Test_controlledGenerationResponseSchema2(t *testing.T) {
	tc := testutil.SystemTest(t)
	w := new(bytes.Buffer)

	location := "us-central1"
	modelName := "gemini-1.5-pro-001"

	err := controlledGenerationResponseSchema2(w, tc.ProjectID, location, modelName)
	if err != nil {
		t.Fatalf("controlledGenerationResponseSchema2: %v", err.Error())
	}

	// We explicitly requested a response in JSON with a specific schema, so we're
	// expecting to properly decode the output as a slice of structured values.
	type Review []struct {
		Rating *int   `json:"rating"`
		Flavor string `json:"flavor"`
	}
	hasEmpty := func(r Review) bool {
		for _, part := range r {
			if part.Rating == nil && part.Flavor == "" {
				return true
			}
		}
		return false
	}

	var reviews []Review
	err = json.Unmarshal(w.Bytes(), &reviews)
	if err != nil {
		t.Errorf(`could not unmarshal response:
%s

into slice of Recipe because: %v`, w.Bytes(), err)
	}
	if len(reviews) == 0 {
		t.Errorf("zero reviews generated")
	}
	for i, review := range reviews {
		if len(review) == 0 {
			t.Errorf("reviews[%d] contains no data", i)
		}
		if hasEmpty(review) {
			t.Errorf("reviews[%d] has an element that contains no data", i)
		}
	}
}

func Test_controlledGenerationResponseSchema3(t *testing.T) {
	tc := testutil.SystemTest(t)
	w := new(bytes.Buffer)

	location := "us-central1"
	modelName := "gemini-1.5-pro-001"

	err := controlledGenerationResponseSchema3(w, tc.ProjectID, location, modelName)
	if err != nil {
		t.Fatalf("controlledGenerationResponseSchema3: %v", err.Error())
	}

	// We explicitly requested a response in JSON with a specific schema, so we're
	// expecting to properly decode the output as a structured object.
	type DayForecast struct {
		Day         string
		Forecast    string
		Humidity    string
		Temperature *int
		WindSpeed   *int `json:"Wind Speed"`
	}
	type Forecast struct {
		Days []DayForecast `json:"forecast"`
	}
	valid := func(f Forecast) bool {
		if len(f.Days) == 0 {
			return false
		}
		for _, d := range f.Days {
			// Check required fields
			if d.Day == "" || d.Forecast == "" || d.Temperature == nil {
				return false
			}
		}
		return true
	}

	var f Forecast
	err = json.Unmarshal(w.Bytes(), &f)
	if err != nil {
		t.Errorf(`could not unmarshal response:
%s

into Forecast because: %v`, w.Bytes(), err)
	}
	if !valid(f) {
		t.Errorf("invalid forecast: %v", f)
	}
}

func Test_controlledGenerationResponseSchema4(t *testing.T) {
	tc := testutil.SystemTest(t)
	w := new(bytes.Buffer)

	location := "us-central1"
	modelName := "gemini-1.5-pro-001"

	err := controlledGenerationResponseSchema4(w, tc.ProjectID, location, modelName)
	if err != nil {
		t.Fatalf("controlledGenerationResponseSchema4: %v", err.Error())
	}

	// We explicitly requested a response in JSON with a specific schema, so we're
	// expecting to properly decode the output as a slice of structured values.
	type Item struct {
		ToDiscard    *int   `json:"to_discard"`
		Subcategory  string `json:"subcategory"`
		SafeHandling *int   `json:"safe_handling"`
		ItemCategory string `json:"item_category"`
		ForResale    *int   `json:"for_resale"`
		Condition    string `json:"condition"`
	}
	var items []Item
	err = json.Unmarshal(w.Bytes(), &items)
	if err != nil {
		t.Errorf(`could not unmarshal response:
%s

into []Item because: %v`, w.Bytes(), err)
	}
	if len(items) == 0 {
		t.Errorf("not items returned")
	}
}

func Test_controlledGenerationResponseSchema6(t *testing.T) {
	tc := testutil.SystemTest(t)
	w := new(bytes.Buffer)

	location := "us-central1"
	modelName := "gemini-1.5-pro-001"

	err := controlledGenerationResponseSchema6(w, tc.ProjectID, location, modelName)
	if err != nil {
		t.Fatalf("controlledGenerationResponseSchema6: %v", err.Error())
	}

	// We explicitly requested a response in JSON with a specific schema, so we're
	// expecting to properly decode the output as a slice.
	type Item struct {
		Object string `json:"object"`
	}
	var items [][]Item
	err = json.Unmarshal(w.Bytes(), &items)
	if err != nil {
		t.Errorf(`could not unmarshal response:
%s

into [][]Item because: %v`, w.Bytes(), err)
	}
	if len(items) == 0 {
		t.Errorf("no items returned")
	}
}

func Test_controlledGenerationResponseSchemaEnum(t *testing.T) {
	tc := testutil.SystemTest(t)
	w := new(bytes.Buffer)

	location := "us-central1"
	modelName := "gemini-1.5-pro-001"

	err := controlledGenerationResponseSchemaEnum(w, tc.ProjectID, location, modelName)
	if err != nil {
		t.Fatalf("controlledGenerationResponseSchemaEnum: %v", err.Error())
	}

	exp := `Candidate label: "documentary"`
	act := w.String()
	if !strings.Contains(w.String(), exp) {
		t.Errorf("expected output to contain text %q, got: %q", exp, act)
	}
}
