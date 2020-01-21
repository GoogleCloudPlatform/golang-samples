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

// Package objects contains samples for creation, listing, deleting objects
// and runs other similar operations on them by using the Google Storage API.
// More documentation is available at
// https://cloud.google.com/storage/docs/json_api/v1/.
package objects

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		fmt.Fprintf(os.Stderr, "GOOGLE_CLOUD_PROJECT environment variable must be set.\n")
		os.Exit(1)
	}
	var o string
	flag.StringVar(&o, "o", "", "source object; in the format of <bucket:object>")
	flag.Parse()

	names := strings.Split(o, ":")
	if len(names) < 2 {
		usage("missing -o flag")
	}
	bucket, object := names[0], names[1]

	if len(os.Args) < 3 {
		usage("missing subcommand")
	}

	switch os.Args[2] {
	case "write":
		if err := write(bucket, object); err != nil {
			log.Fatalf("Cannot write object: %v", err)
		}
	case "read":
		data, err := read(bucket, object)
		if err != nil {
			log.Fatalf("Cannot read object: %v", err)
		}
		fmt.Printf("Object contents: %s\n", data)
	case "metadata":
		attrs, err := attrs(bucket, object)
		if err != nil {
			log.Fatalf("Cannot get object metadata: %v", err)
		}
		fmt.Printf("Object metadata: %v\n", attrs)
	case "makepublic":
		if err := makePublic(bucket, object); err != nil {
			log.Fatalf("Cannot to make object public: %v", err)
		}
	case "delete":
		if err := delete(bucket, object); err != nil {
			log.Fatalf("Cannot to delete object: %v", err)
		}
	}
}

// TODO(jbd): Add test for downloadUsingRequesterPays.

const helptext = `usage: objects -o=bucket:name [subcommand] <args...>

subcommands:
	- write
	- read
	- metadata
	- makepublic
	- delete
`

func usage(msg string) {
	if msg != "" {
		fmt.Fprintln(os.Stderr, msg)
	}
	fmt.Fprintln(os.Stderr, helptext)
	os.Exit(2)
}
