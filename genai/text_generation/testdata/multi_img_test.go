package text_generation_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/genai/tools/text_generation"
)

func TestGenerateWithMultiImg(t *testing.T) {
	
	dummyImagePath := filepath.Join(os.TempDir(), "latte.jpg")
	err := os.WriteFile(dummyImagePath, []byte("This is a dummy JPEG image."), 0644)
	if err != nil {
		t.Fatalf("Failed to create dummy latte.jpg: %v", err)
	}
	defer os.Remove(dummyImagePath)

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	err = os.Chdir("./")
	if err != nil {
		t.Fatalf("Failed to change directory to testdata: %v", err)
	}
	defer os.Chdir(originalDir)

	var buf bytes.Buffer
	err = text_generation.generateWithMultiImg(&buf)
	if err != nil {
		t.Errorf("generateWithMultiImg failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Errorf("generateWithMultiImg produced empty output")
	}

	t.Logf("Output:\n%s", buf.String())
}