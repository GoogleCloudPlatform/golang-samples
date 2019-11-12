package cloudrunci

import (
	"strings"
	"testing"
)

// TestServiceValidateErrors checks for errors in the Service definition.
func TestServiceValidateErrors(t *testing.T) {
	service := Service{Name: "my-serivce"}
	if err := service.validate(); err == nil {
		t.Errorf("service.validate: expected error 'Project ID missing', got success")
	}

	service.ProjectID = "my-project"
	if err := service.validate(); err == nil {
		t.Errorf("service.validate: expected error 'Platform configuration missing', got success")
	}
}

// TestServiceStateErrors checks that a service in the wrong state will be blocked from the requested operation.
func TestServiceStateErrors(t *testing.T) {
	service := NewService("my-serivce", "my-project")

	want := "Request called before Deploy"
	if _, err := service.Request("GET", "/"); !strings.Contains(err.Error(), want) {
		t.Errorf("service.Request: error expected '%s', got %s", want, err.Error())
	}

	want = "NewRequest called before Deploy"
	if _, err := service.NewRequest("GET", "/"); !strings.Contains(err.Error(), want) {
		t.Errorf("service.NewRequest: error expected '%s', got %s", want, err.Error())
	}

	want = "URL called before Deploy"
	if _, err := service.URL("/"); !strings.Contains(err.Error(), want) {
		t.Errorf("service.URL: error expected '%s', got %s", want, err.Error())
	}

	want = "container image already built"
	service.built = true
	if err := service.Build(); !strings.Contains(err.Error(), want) {
		t.Errorf("service.Build: error expected '%s', got %s", want, err.Error())
	}
}
