package snippets

// [START aiplatform_gemini_get_started]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"cloud.google.com/go/vertexai/genai"
)

var projectId = "PROJECT_ID"
var region = "us-central1"

func tryGemini(w io.Writer, projectId string, region string, modelName string) error {

	client, err := genai.NewClient(context.TODO(), projectId, region)
	gemini := client.GenerativeModel("gemini-pro-vision")

	img := genai.FileData{
		MIMEType: "image/jpeg",
		FileURI:  "https://storage.googleapis.com/generativeai-downloads/images/scones.jpg",
	}
	prompt := genai.Text("What is in this image?")
	resp, err := gemini.GenerateContent(context.Background(), img, prompt)
	if err != nil {
		return fmt.Errorf("error generating content: %w", err)
	}
	rb, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Fprintln(w, string(rb))
	return nil
}

// [END aiplatform_gemini_get_started]
