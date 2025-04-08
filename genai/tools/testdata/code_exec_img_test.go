package tools_test 

import (
    "bytes"
    "testing"

    "github.com/GoogleCloudPlatform/golang-samples/genai/tools"
)

func TestGenerateWithCodeExecAndImg(t *testing.T) {
    var buf bytes.Buffer
    err := tools.generateWithCodeExecAndImg(&buf)
    if err != nil {
        t.Errorf("generateWithCodeExecAndImg failed: %v", err)
    }
}