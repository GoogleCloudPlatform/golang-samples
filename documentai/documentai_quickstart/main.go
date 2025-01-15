// Copyright 2023 Google LLC
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

// quickstart sends a file at a given filePath for online processing.
package documentai_snippets

// [START documentai_quickstart]
import (
	"context"
	"flag"
	"fmt"
	"os"

	documentai "cloud.google.com/go/documentai/apiv1"
	"cloud.google.com/go/documentai/apiv1/documentaipb"
	"google.golang.org/api/option"
)

func main() {
	projectID := flag.String("project_id", "PROJECT_ID", "Cloud Project ID")
	location := flag.String("location", "us", "The Processor location")
	// Create a Processor before running sample
	processorID := flag.String("processor_id", "aaaaaaaa", "The Processor ID")
	filePath := flag.String("file_path", "invoice.pdf", "The path to the file to parse")
	mimeType := flag.String("mime_type", "application/pdf", "The mimeType of the file")
	flag.Parse()

	ctx := context.Background()

	endpoint := fmt.Sprintf("%s-documentai.googleapis.com:443", *location)
	client, err := documentai.NewDocumentProcessorClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		fmt.Println(fmt.Errorf("error creating Document AI client: %w", err))
	}
	defer client.Close()

	// Open local file.
	data, err := os.ReadFile(*filePath)
	if err != nil {
		fmt.Println(fmt.Errorf("os.ReadFile: %w", err))
	}

	req := &documentaipb.ProcessRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/processors/%s", *projectID, *location, *processorID),
		Source: &documentaipb.ProcessRequest_RawDocument{
			RawDocument: &documentaipb.RawDocument{
				Content:  data,
				MimeType: *mimeType,
			},
		},
	}
	resp, err := client.ProcessDocument(ctx, req)
	if err != nil {
		fmt.Println(fmt.Errorf("processDocument: %w", err))
	}

	// Handle the results.
	document := resp.GetDocument()
	fmt.Printf("Document Text: %s", document.GetText())
}

// [END documentai_quickstart]
