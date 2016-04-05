package e2e

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/aeintegrate"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestHelloWorld(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	helloworld := &aeintegrate.App{
		Name:      "hw",
		Dir:       tc.Path("docs", "managed_vms", "helloworld"),
		ProjectID: tc.ProjectID,
	}

	bodyShouldContain(t, helloworld, "/", "Hello world!")
}

func TestDatastore(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	datastore := &aeintegrate.App{
		Name:      "ds",
		Dir:       tc.Path("docs", "managed_vms", "datastore"),
		ProjectID: tc.ProjectID,
		Env: map[string]string{
			"GCLOUD_DATASET_ID": tc.ProjectID,
		},
	}

	bodyShouldContain(t, datastore, "/", "Succesfully stored")
}

func TestMemcache(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	memcache := &aeintegrate.App{
		Name:      "mem",
		Dir:       tc.Path("docs", "managed_vms", "memcache"),
		ProjectID: tc.ProjectID,
	}

	bodyShouldContain(t, memcache, "/", "Count")
}

func bodyShouldContain(t *testing.T, p *aeintegrate.App, path, shouldContain string) {
	if p.Deployed() {
		t.Fatalf("[%s] expected non-deployed app", p.Name)
	}

	if err := p.Deploy(); err != nil {
		t.Fatalf("could not deploy %s: %v", p.Name, err)
	}

	defer p.Cleanup()

	resp, err := p.Get(path)
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read body: %v", err)
	}

	if !strings.Contains(string(b), shouldContain) {
		t.Fatalf("wanted to contain %q, but got body: %q", shouldContain, string(b))
	}
}
