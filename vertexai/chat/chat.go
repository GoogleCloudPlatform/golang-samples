package chat

// [START aiplatform_gemini_multiturn_chat]
import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/vertexai/genai"
)

var projectId = "PROJECT_ID"
var region = "us-central1"
var modelName = "gemini-pro-vision"

func makeChatRequests(projectId string, region string, modelName string) error {
	client, err := genai.NewClient(context.TODO(), projectId, region)
	if err != nil {
		return fmt.Errorf("error creating client: %v", err)
	}
	defer client.Close()

	gemini := client.GenerativeModel(modelName)
	chat := gemini.StartChat()

	r, err := chat.SendMessage(
		context.TODO(),
		genai.Text("Hello"))
	if err != nil {
		return err
	}
	rb, _ := json.MarshalIndent(r, "", "  ")
	fmt.Println(string(rb))

	r, err = chat.SendMessage(
		context.TODO(),
		genai.Text("What are all the colors in a rainbow?"))
	if err != nil {
		return err
	}
	rb, _ = json.MarshalIndent(r, "", "  ")
	fmt.Println(string(rb))

	r, err = chat.SendMessage(
		context.TODO(),
		genai.Text("Why does it appear when it rains?"))
	if err != nil {
		return err
	}
	rb, _ = json.MarshalIndent(r, "", "  ")
	fmt.Println(string(rb))

	return nil
}

// [END aiplatform_gemini_multiturn_chat]
