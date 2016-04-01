// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package aeintegrate facilitates end-to-end testing against the production Google App Engine.
//
// This is a specialized tool that might be used in addition to unit tests. It
// calls the `gcloud app` command directly, which is in preview/beta, so expect
// any tests using this package to spontaneously break.
//
// aedeploy (go get google.golang.org/appengine/cmd/aedeploy) and gcloud
// (https://cloud.google.com/sdk) must be installed. You must be authorized via
// the gcloud command-line tool (`gcloud auth login`) and the project ID must be
// set via `gcloud config set project`.
//
// You may specify the locations of aedeploy and/or gcloud via the AEDEPLOY_BIN
// and GCLOUD_BIN environment variables, respectively.
//
// Sample usage with `go test`:
//
// 	package myapp
//
// 	import (
// 		"testing"
// 		"google.golang.org/appengine/aeintegrate"
// 	)
//
// 	func TestApp(t *testing.T) {
// 		t.Parallel()
// 		app := aeintegrate.App{Name: "A", Dir: "app"},
// 		if err := app.Deploy(); err != nil {
// 			t.Fatalf("could not deploy app: %v", err)
// 		}
// 		defer app.Cleanup()
// 		resp, err := app.Get("/")
// 		...
// 	}
package aeintegrate

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	appengine "google.golang.org/api/appengine/v1beta4"

	"gopkg.in/yaml.v2"
)

// runID is an identifier that changes between runs.
var runID = time.Now().Format("20060102-150405")

// App describes an App Engine application.
type App struct {
	// Name is an ID, used for logging and to generate a unique version to this run.
	Name string

	// The root directory containing the app's source code.
	Dir string

	// The configuration (app.yaml) file, relative to Dir. Defaults to "app.yaml".
	AppYaml string

	// The project to deploy to.
	ProjectID string

	// Additional runtime environment variable overrides for the app.
	// NOTE: does not yet work in gcloud.
	Env map[string]string

	deployed bool // Whether the app has been deployed.

	module string // The Module ID (read from the app.yaml)

	adminService *appengine.Service // Used during clean up to delete the deployed version.
}

// Deployed reports whether the application has been deployed.
func (p *App) Deployed() bool {
	return p.deployed
}

// Get issues a GET request against the base URL of the deployed application.
func (p *App) Get(path string) (*http.Response, error) {
	if !p.deployed {
		return nil, errors.New("Get called before Deploy")
	}
	url, _ := p.URL(path)
	return http.Get(url)
}

// URL prepends the deployed application's base URL to the given path.
// Returns an error if the application has not been deployed.
func (p *App) URL(path string) (string, error) {
	if !p.deployed {
		return "", errors.New("URL called before Deploy")
	}
	return fmt.Sprintf("https://%s-dot-%s-dot-%s.appspot.com%s", p.version(), p.module, p.ProjectID, path), nil
}

// version returns the version that the app will be deployed to.
func (p *App) validate() error {
	if p.ProjectID == "" {
		return errors.New("Project ID missing")
	}
	return nil
}

// version returns the version that the app will be deployed to.
func (p *App) version() string {
	return p.Name + "-" + runID
}

// Deploy deploys the application to App Engine. If the deployment fails, it tries to clean up the failed deployment.
func (p *App) Deploy() error {
	// Don't deploy unless we're certain everything is ready for deployment
	// (i.e. admin client is authenticated and authorized)
	if err := p.validate(); err != nil {
		return err
	}
	if _, err := p.Module(); err != nil {
		return fmt.Errorf("could not get module id: %v", err)
	}
	if err := p.setupAdmin(); err != nil {
		return fmt.Errorf("could not setup admin service: %v", err)
	}

	log.Printf("(%s) Deploying...", p.Name)

	out, err := p.deployCmd().CombinedOutput()
	// TODO: add a flag for verbose output (e.g. when running with binary created with `go test -c`)
	if err != nil {
		log.Printf("(%s) Output from deploy:", p.Name)
		os.Stderr.Write(out)
		// Try to clean up resources.
		p.Cleanup()
	} else if err == nil {
		p.deployed = true
		log.Printf("(%s) Deploy successful.", p.Name)
	}
	return err
}

// appyaml returns the path of the config file.
func (p *App) appyaml() string {
	appyaml := p.AppYaml
	if appyaml == "" {
		appyaml = "app.yaml"
	}
	return appyaml
}

func (p *App) deployCmd() *exec.Cmd {
	gcloudBin := os.Getenv("GCLOUD_BIN")
	if gcloudBin == "" {
		gcloudBin = "gcloud"
	}
	aedeploy := os.Getenv("AEDEPLOY_BIN")
	if aedeploy == "" {
		aedeploy = "aedeploy"
	}
	// NOTE: if the "preview" and/or "app" modules are not available, and this is
	// run in parallel, such as from AppGroup.Deploy, gcloud will attempt to
	// install those components multiple times and will eventually fail on IO.
	args := []string{gcloudBin,
		"--quiet",
		"preview", "app", "deploy", p.appyaml(),
		"--project", p.ProjectID,
		"--version", p.version(),
		"--no-promote",
	}
	cmd := exec.Command(aedeploy, args...)
	cmd.Dir = p.Dir

	if len(p.Env) != 0 {
		cmd.Env = os.Environ()
		keys := make([]string, 0)
		for k, v := range p.Env {
			cmd.Env = append(cmd.Env, k+"="+v)
			keys = append(keys, k)
		}
		args = append(args, "--env-vars", strings.Join(keys, ","))
	}
	return cmd
}

// Module returns the Module ID, which is read from the app.yaml file.
func (p *App) Module() (string, error) {
	if p.module != "" {
		return p.module, nil
	}

	b, err := ioutil.ReadFile(filepath.Join(p.Dir, p.appyaml()))
	if err != nil {
		return "", err
	}

	s := &struct {
		Module string `yaml:"module"`
	}{}

	if err := yaml.Unmarshal(b, s); err != nil {
		return "", err
	}

	if s.Module == "" {
		s.Module = "default"
	}

	p.module = s.Module
	return p.module, err
}

// setupAdmin populates p.adminService and checks that the user is authenticated and project ID is valid.
func (p *App) setupAdmin() error {
	c, err := google.DefaultClient(context.Background(), appengine.CloudPlatformScope)
	if err != nil {
		return err
	}

	if p.adminService, err = appengine.New(c); err != nil {
		return err
	}

	if err := p.validate(); err != nil {
		return err
	}

	// Check that the user is authenticated, etc.
	_, err = p.adminService.Apps.Get(p.ProjectID).Do()
	return err
}

// Cleanup deletes the created version from App Engine.
func (p *App) Cleanup() error {
	// NOTE: don't check whether p.deployed is set.
	// We may want to attempt to clean up if deployment failed.
	// However, we require adminService to be set up, which happens during Deploy().
	if p.adminService == nil {
		return errors.New("Cleanup called before Deploy")
	}

	if err := p.validate(); err != nil {
		return err
	}

	log.Printf("(%s) Cleaning up.", p.Name)

	var err error
	for try := 0; try < 10; try++ {
		_, err = p.adminService.Apps.Modules.Versions.Delete(p.ProjectID, p.module, p.version()).Do()
		if err == nil {
			log.Printf("(%s) Succesfully cleaned up.", p.Name)
			break
		}
		time.Sleep(time.Second)
	}
	return err
}

func gcloudBin() string {
	bin := os.Getenv("GCLOUD_BIN")
	if bin == "" {
		return "gcloud"
	}
	return bin
}
