package functest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"google.golang.org/api/idtoken"
)

// URL assemvbles the URL based on the Cloud Functions URL convention.
func (f *CloudFunction) URL() *url.URL {
	if f.url == nil {
		uStr := fmt.Sprintf("https://%s-%s.cloudfunctions.net/%s", f.Region, f.ProjectID, f.DeployName())
		u, err := url.Parse(uStr)
		if err != nil {
			log.Fatal("url.Parse: %v", err)
		}
		f.url = u
	}
	return f.url
}

// HTTPClient provides an HTTP request with built-in authentication support.
func (f *CloudFunction) HTTPClient() (*http.Client, error) {
	if !f.deployed {
		f.log("[WARNING] function not deployed")
	}
	if f.client == nil {
		url := f.URL()
		ctx := context.Background()
		client, err := idtoken.NewClient(ctx, url.String())
		if err != nil {
			return nil, fmt.Errorf("idtoken.NewClient: %w", err)
		}
		client.Transport = loggingTransport{client.Transport}
		f.client = client
	}

	return f.client, nil
}

// HTTP Transport wrapper that logs before and after requests.
type loggingTransport struct {
	base http.RoundTripper
}

func (t loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Printf("http: request ===> %s %s", req.Method, req.URL)
	// The ID token is added to a cloned request object.
	// It is not accessible from functest.
	resp, err := t.base.RoundTrip(req)
	log.Printf("http: response <== %d %s", resp.StatusCode, req.URL)
	return resp, err
}
