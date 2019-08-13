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

// Command gimme-acc provides access to a pool of projects.
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
	ServiceAccounts []ServiceAccount
}

type ServiceAccount struct {
	Email       string
	LeaseExpiry time.Time
}

func (p *Pool) Get(email string) (*ServiceAccount, bool) {
	for i := range p.ServiceAccounts {
		acc := &p.ServiceAccounts[i]
		if acc.Email == email {
			return acc, true
		}
	}
	return nil, false
}

func (p *Pool) Add(email string) (ok bool) {
	if _, ok := p.Get(email); ok {
		return false
	}
	p.ServiceAccounts = append(p.ServiceAccounts, ServiceAccount{Email: email})
	return true
}

func (p *Pool) Lease(d time.Duration) (*ServiceAccount, bool) {
	if len(p.ServiceAccounts) == 0 {
		return nil, false
	}

	oldest := &p.ServiceAccounts[0]
	for i := range p.ServiceAccounts {
		acc := &p.ServiceAccounts[i]
		if acc.LeaseExpiry.Before(oldest.LeaseExpiry) {
			oldest = acc
		}
	}
	if !oldest.Expired() {
		return nil, false
	}
	oldest.LeaseExpiry = time.Now().Add(d)
	return oldest, true
}

func (p *ServiceAccount) Expired() bool {
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
	gimme-acc -project=[meta project ID] command
	gimme-acc -project=[meta project ID] -output=list status

Commands:
	lease [duration]    Leases a service account for a given duration. Prints the service account email to stdout.
	done [service account email]   Returns a service account to the pool.
	version             Prints the version of gimme-acc.

Administrative commands:
	pool-add [service account email]       Adds a service account to the pool.
	pool-rm  [service account email]       Removes a service account from the pool.
	status                      Displays the current status of the meta project. Respects -output.
`)

	if flag.Arg(0) == "version" {
		fmt.Printf("gimme-acc %s; built at %s\n", version, buildDate)
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
		key := ds.NameKey("Pool", "acc-pool", nil)
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

	var acc *ServiceAccount
	err = withPool(ctx, func(pool *Pool) error {
		var ok bool
		acc, ok = pool.Lease(d)
		if !ok {
			return errors.New("Could not find a free service account. Try again soon.")
		}
		return nil
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Leased! %s is yours for %s.\n", acc.Email, d)
	fmt.Print(acc.Email)
	return nil
}

func done(ctx context.Context, accEmail string) error {
	if accEmail == "" {
		return errors.New("must provide service account email")
	}
	err := withPool(ctx, func(pool *Pool) error {
		acc, ok := pool.Get(accEmail)
		if !ok {
			return fmt.Errorf("Could not find service account %s in pool.", accEmail)
		}
		acc.LeaseExpiry = time.Now().Add(-10 * time.Second)
		return nil
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Returned %s to the pool.\n", accEmail)
	return nil
}

func status(ctx context.Context) error {
	return withPool(ctx, func(pool *Pool) error {
		if *format == "" {
			fmt.Printf("%-8s %s\n", "LEASE", "SERVICE ACCOUNT")
		}
		for _, acc := range pool.ServiceAccounts {
			exp := ""
			if !acc.Expired() {
				secs := acc.LeaseExpiry.Sub(time.Now()) / time.Second * time.Second
				exp = secs.String()
			}
			switch *format {
			case "":
				fmt.Printf("%-8s %s\n", exp, acc.Email)
			case "list":
				fmt.Printf("%s\n", acc.Email)
			default:
				return errors.New("output may be '', 'list'")
			}
		}
		return nil
	})
}

func addToPool(ctx context.Context, email string) error {
	if email == "" {
		return errors.New("must provide service account email")
	}
	return withPool(ctx, func(pool *Pool) error {
		if !pool.Add(email) {
			return fmt.Errorf("%s already in pool", email)
		}
		return nil
	})
}

func removeFromPool(ctx context.Context, email string) error {
	if email == "" {
		return errors.New("must provide service account email")
	}
	return withPool(ctx, func(pool *Pool) error {
		if _, ok := pool.Get(email); !ok {
			return fmt.Errorf("%s not in pool", email)
		}

		accs := make([]ServiceAccount, 0)
		for _, acc := range pool.ServiceAccounts {
			if acc.Email == email {
				continue
			}
			accs = append(accs, acc)
		}
		pool.ServiceAccounts = accs
		return nil
	})
}
