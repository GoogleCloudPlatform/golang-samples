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

package converter

import (
	"testing"

	storagepb "google.golang.org/genproto/googleapis/cloud/bigquery/storage/v1"
)

func TestValidate(t *testing.T) {
	testcases := []struct {
		description      string
		schema           *storagepb.TableSchema
		cols             []string
		wantConverterErr bool
		wantErr          bool
	}{
		{
			description:      "nil schema",
			wantConverterErr: true,
		},
		{
			description: "required field",
			schema: &storagepb.TableSchema{
				Fields: []*storagepb.TableFieldSchema{
					{
						Name: "foo",
						Type: storagepb.TableFieldSchema_STRING,
						Mode: storagepb.TableFieldSchema_NULLABLE,
					},
					{
						Name: "bar",
						Type: storagepb.TableFieldSchema_STRING,
						Mode: storagepb.TableFieldSchema_REQUIRED,
					},
				},
			},
			cols:    []string{"foo"},
			wantErr: true,
		},
		{
			description: "mismatched type",
			schema: &storagepb.TableSchema{
				Fields: []*storagepb.TableFieldSchema{
					{
						Name: "foo",
						Type: storagepb.TableFieldSchema_STRING,
						Mode: storagepb.TableFieldSchema_NULLABLE,
					},
					{
						Name: "bar",
						Type: storagepb.TableFieldSchema_GEOGRAPHY,
						Mode: storagepb.TableFieldSchema_NULLABLE,
					},
				},
			},
			cols:    []string{"foo", "bar"},
			wantErr: true,
		},
		{
			description: "compatible",
			schema: &storagepb.TableSchema{
				Fields: []*storagepb.TableFieldSchema{
					{
						Name: "foo",
						Type: storagepb.TableFieldSchema_STRING,
						Mode: storagepb.TableFieldSchema_NULLABLE,
					},
					{
						Name: "bar",
						Type: storagepb.TableFieldSchema_STRING,
						Mode: storagepb.TableFieldSchema_NULLABLE,
					},
				},
			},
			cols: []string{"foo", "bar"},
		},
	}
	for _, tc := range testcases {
		conv, err := NewCSVConverter(tc.schema)
		if err != nil {

			if tc.wantConverterErr {
				continue
			}
			t.Errorf("%s NewCSVConverter: %v", tc.description, err)
			continue
		}

		if err == nil && tc.wantConverterErr {
			t.Errorf("%s: expected converter err, got success", tc.description)
			continue
		}

		err = conv.Validate(tc.cols)

		if err != nil && !tc.wantErr {
			t.Errorf("%s: wanted success, got err: %v", tc.description, err)
		}
		if err == nil && tc.wantErr {
			t.Errorf("%s: wanted err, got success", tc.description)
		}
	}
}

func TestConvert(t *testing.T) {
	conv, err := NewCSVConverter(&storagepb.TableSchema{
		Fields: []*storagepb.TableFieldSchema{
			{
				Name: "foo",
				Type: storagepb.TableFieldSchema_STRING,
				Mode: storagepb.TableFieldSchema_NULLABLE,
			},
			{
				Name: "bar",
				Type: storagepb.TableFieldSchema_STRING,
				Mode: storagepb.TableFieldSchema_REQUIRED,
			},
		}})
	if err != nil {
		t.Fatalf("failed to instantiate converter")
	}

	if err := conv.Validate([]string{"foo", "bar"}); err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	testMap := map[string]string{
		"foo": "val1",
		"bar": "val2",
	}
	b, err := conv.Convert(testMap)
	if err != nil {
		t.Errorf("conversion failed: %v", err)
	}
	if len(b) == 0 {
		t.Errorf("bytes are empty")
	}

}
