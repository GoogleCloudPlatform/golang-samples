// Copyright 2019 Google LLC
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

// Command gimmeproj provides access to a pool of projects.
//
// The metadata about the project pool is stored in Cloud Datastore in a meta-project.
// Projects are leased for a certain duration, and automatically returned to the pool when the lease expires.
// Projects should be returned before the lease expires.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	ds "cloud.google.com/go/datastore"
)

var (
	metaProject = flag.String("project", "", "Meta-project that manages the pool.")
	format      = flag.String("output", "", "Output format for selected operations. Options include: list")
	datastore   *ds.Client

	version   = "dev"
	buildDate = "unknown"
)

type Pool struct {
	Projects []Project
}

type Project struct {
	ID          string
	LeaseExpiry time.Time
}

func (p *Pool) Get(projID string) (*Project, bool) {
	for i := range p.Projects {
		proj := &p.Projects[i]
		if proj.ID == projID {
			return proj, true
		}
	}
	return nil, false
}

func (p *Pool) Add(proj string) (ok bool) {
	if _, ok := p.Get(proj); ok {
		return false
	}
	p.Projects = append(p.Projects, Project{ID: proj})
	return true
}

func (p *Pool) Lease(d time.Duration) (*Project, bool) {
	if len(p.Projects) == 0 {
		return nil, false
	}

	oldest := &p.Projects[0]
	for i := range p.Projects {
		proj := &p.Projects[i]
		if proj.LeaseExpiry.Before(oldest.LeaseExpiry) {
			oldest = proj
		}
	}
	if !oldest.Expired() {
		return nil, false
	}
	oldest.LeaseExpiry = time.Now().Add(d)
	return oldest, true
}

func (p *Project) Expired() bool {
	return time.Now().After(p.LeaseExpiry)
}

func main() {
	flag.Parse()
	if err := submain(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
}

func submain() error {
	ctx := context.Background()

	usage := errors.New(`
Usage:
	gimmeproj -project=[meta project ID] command
	gimmeproj -project=[meta project ID] -output=list status

Commands:
	lease [duration]    Leases a project for a given duration. Prints the project ID to stdout.
	done [project ID]   Returns a project to the pool.
	version             Prints the version of gimmeproj.

Administrative commands:
	pool-add [project ID]       Adds a project to the pool.
	pool-rm  [project ID]       Removes a project from the pool.
	status                      Displays the current status of the meta project. Respects -output.
`)

	if flag.Arg(0) == "version" {
		fmt.Printf("gimmeproj %s; built at %s\n", version, buildDate)
		return nil
	}

	if *metaProject == "" {
		fmt.Fprintln(os.Stderr, "-project flag is required.")
		return usage
	}

	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "Missing command.")
		return usage
	}

	var err error
	datastore, err = ds.NewClient(ctx, *metaProject)
	if err != nil {
		return fmt.Errorf("datastore.NewClient: %v", err)
	}

	switch flag.Arg(0) {
	case "help":
		fmt.Fprintln(os.Stderr, usage.Error())
		return nil
	case "lease":
		return lease(ctx, flag.Arg(1))
	case "pool-add":
		return addToPool(ctx, flag.Arg(1))
	case "pool-rm":
		return removeFromPool(ctx, flag.Arg(1))
	case "status":
		return status(ctx)
	case "done":
		return done(ctx, flag.Arg(1))
	}
	fmt.Fprintln(os.Stderr, "Unknown command.")
	return usage
}

// withPool runs the given function in a transaction, saving the state of the pool if the function returns with a non-nil error.
func withPool(ctx context.Context, f func(pool *Pool) error) error {
	_, err := datastore.RunInTransaction(ctx, func(tx *ds.Transaction) error {
		key := ds.NameKey("Pool", "pool", nil)
		var pool Pool
		if err := tx.Get(key, &pool); err != nil {
			if err == ds.ErrNoSuchEntity {
				if _, err := tx.Put(key, &pool); err != nil {
					return fmt.Errorf("Initial Pool.Put: %v", err)
				}
			} else {
				return fmt.Errorf("Pool.Get: %v", err)
			}
		}
		if err := f(&pool); err != nil {
			return err
		}
		_, err := tx.Put(key, &pool)
		if err != nil {
			return fmt.Errorf("Pool.Put: %v", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("datastore: %v", err)
	}
	return nil
}

func lease(ctx context.Context, duration string) error {
	if duration == "" {
		return errors.New("must provide a duration (e.g. 10m). See https://golang.org/pkg/time/#ParseDuration")
	}
	d, err := time.ParseDuration(duration)
	if err != nil {
		return fmt.Errorf("Could not parse duration: %v", err)
	}

	var proj *Project
	err = withPool(ctx, func(pool *Pool) error {
		var ok bool
		proj, ok = pool.Lease(d)
		if !ok {
			return errors.New("Could not find a free project. Try again soon.")
		}
		return nil
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Leased! %s is yours for %s.\n", proj.ID, d)
	fmt.Print(proj.ID)
	return nil
}

func done(ctx context.Context, projectID string) error {
	if projectID == "" {
		return errors.New("must provide project id")
	}
	err := withPool(ctx, func(pool *Pool) error {
		proj, ok := pool.Get(projectID)
		if !ok {
			return fmt.Errorf("Could not find project %s in project pool.", projectID)
		}
		proj.LeaseExpiry = time.Now().Add(-10 * time.Second)
		return nil
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Returned %s to the pool.\n", projectID)
	return nil
}

func status(ctx context.Context) error {
	return withPool(ctx, func(pool *Pool) error {
		if *format == "" {
			fmt.Printf("%-8s %s\n", "LEASE", "PROJECT")
		}
		for _, proj := range pool.Projects {
			exp := ""
			if !proj.Expired() {
				secs := proj.LeaseExpiry.Sub(time.Now()) / time.Second * time.Second
				exp = secs.String()
			}
			switch *format {
			case "":
				fmt.Printf("%-8s %s\n", exp, proj.ID)
			case "list":
				fmt.Printf("%s\n", proj.ID)
			default:
				return errors.New("output may be '', 'list'")
			}
		}
		return nil
	})
}

func addToPool(ctx context.Context, proj string) error {
	if proj == "" {
		return errors.New("must provide project id")
	}
	return withPool(ctx, func(pool *Pool) error {
		if !pool.Add(proj) {
			return fmt.Errorf("%s already in pool", proj)
		}
		return nil
	})
}

func removeFromPool(ctx context.Context, projectID string) error {
	if projectID == "" {
		return errors.New("must provide project id")
	}
	return withPool(ctx, func(pool *Pool) error {
		if _, ok := pool.Get(projectID); !ok {
			return fmt.Errorf("%s not in pool", projectID)
		}

		projs := make([]Project, 0)
		for _, proj := range pool.Projects {
			if proj.ID == projectID {
				continue
			}
			projs = append(projs, proj)
		}
		pool.Projects = projs
		return nil
	})
}
