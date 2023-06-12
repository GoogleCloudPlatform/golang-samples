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

// Command cleaneversions deletes App Engine versions for a given project, service and/or version ID filter.
//
//	Usage of cleanaeversions:
//	  -async
//	      Don't wait for successful deletion.
//	  -filter regexp
//	      Filter regexp for version IDs. If empty, attempts to clean all versions.
//	  -n  Dry run.
//	  -project Project ID
//	      Project ID to clean.
//	  -service Service/module ID
//	      Service/module ID to clean. If omitted, cleans all services.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/oauth2/google"

	appengine "google.golang.org/api/appengine/v1"
)

var (
	proj    = flag.String("project", "", "`Project ID` to clean.")
	service = flag.String("service", "", "`Service/module ID` to clean. If omitted, cleans all services.")
	filter  = flag.String("filter", "", "Filter `regexp` for version IDs. If empty, attempts to clean all versions.")
	async   = flag.Bool("async", false, "Don't wait for successful deletion.")
	dryRun  = flag.Bool("n", false, "Dry run.")
)

var gae *appengine.APIService

type pendingDelete struct {
	service string
	version string
	op      *appengine.Operation
}

func main() {
	flag.Parse()
	if *proj == "" {
		fmt.Fprintln(os.Stderr, "-project flag is required")
		flag.Usage()
		os.Exit(2)
	}

	filterRE, err := regexp.Compile(*filter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Filter is not a valid regexp: %v", err)
		os.Exit(2)
	}
	_ = filterRE

	ctx := context.Background()
	hc, err := google.DefaultClient(ctx, appengine.CloudPlatformScope)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create DefaultClient: %v", err)
		os.Exit(1)
	}
	gae, err = appengine.New(hc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create App Engine service: %v", err)
		os.Exit(1)
	}

	var services []string
	if *service != "" {
		services = append(services, *service)
	} else {
		if err := gae.Apps.Services.List(*proj).Pages(ctx, func(lsr *appengine.ListServicesResponse) error {
			for _, s := range lsr.Services {
				services = append(services, s.Id)
			}
			return nil
		}); err != nil {
			fmt.Fprintf(os.Stderr, "Could not list App Engine services: %v", err)
			os.Exit(1)
		}
	}

	var pending []pendingDelete

	for _, service := range services {
		if err := gae.Apps.Services.Versions.List(*proj, service).Pages(ctx, func(lvr *appengine.ListVersionsResponse) error {
			for _, v := range lvr.Versions {
				if !filterRE.MatchString(v.Id) {
					continue
				}

				log.Printf("Deleting %s/%s", service, v.Id)
				if *dryRun {
					continue
				}

				op, err := gae.Apps.Services.Versions.Delete(*proj, service, v.Id).Do()
				if err != nil {
					log.Printf("Could not delete version %s/%s: %v\n", service, v.Id, err)
				} else {
					pending = append(pending, pendingDelete{service: service, version: v.Id, op: op})
				}
			}
			return nil
		}); err != nil {
			fmt.Fprintf(os.Stderr, "Could not list versions for %q: %v\n", service, err)
			os.Exit(1)
		}
	}

	if *async {
		log.Printf("Not waiting for operations to complete. Exiting.")
		os.Exit(0)
	}

	log.Printf("Waiting for operations to complete.")

	var failed int64
	var wg sync.WaitGroup
	wg.Add(len(pending))
	for _, pd := range pending {
		pd := pd
		go func() {
			if err := waitForCompletion(pd); err != nil {
				log.Printf("FAILED %v/%v/%v: %v", *proj, pd.service, pd.version, err)
				atomic.AddInt64(&failed, 1)
			} else {
				log.Printf("Deleted %v/%v/%v", *proj, pd.service, pd.version)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if failed != 0 {
		log.Printf("FAILED (%d)", failed)
		os.Exit(1)
	}
}

func waitForCompletion(pd pendingDelete) error {
	parts := strings.Split(pd.op.Name, "/")
	id := parts[len(parts)-1]
	for {
		op, err := gae.Apps.Operations.Get(*proj, id).Do()
		if err != nil {
			return err
		}
		if !op.Done {
			// 5 to 10 second sleep.
			time.Sleep(time.Duration(5+rand.Float64()*5) * time.Second)
			continue
		}
		if op.Error == nil {
			return nil
		}
		return fmt.Errorf("%s (code %d)", op.Error.Message, op.Error.Code)
	}
}
