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

// Command wrlookup is a tool for looking up URLs via the command-line.
//
// The tool reads one URL per line from STDIN and checks every URL against
// the Web Risk API. The "Safe" or "Unsafe" verdict is printed to STDOUT.
// If an error occurred, debug information may be printed to STDERR.
//
// To build the tool:
//	$ go get github.com/google/webrisk/cmd/wrlookup
//
// Example usage:
//	$ wrlookup -apikey $APIKEY
//	https://google.com
//	Safe URL: https://google.com
//	http://bad1url.org
//	Unsafe URL: [{bad1url.org {MALWARE ANY_PLATFORM URL}}]
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/GoogleCloudPlatform/golang-samples/webrisk"
	"os"
)

var (
	apiKeyFlag    = flag.String("apikey", "", "specify your Web Risk API key")
	databaseFlag  = flag.String("db", "", "path to the Web Risk database. By default persistent storage is disabled (not recommended).")
	serverURLFlag = flag.String("server", webrisk.DefaultServerURL, "Web Risk API server address.")
	proxyFlag     = flag.String("proxy", "", "proxy to use to connect to the HTTP server")
)

const usage = `wrlookup: command-line tool to lookup URLs with Web Risk.

Tool reads one URL per line from STDIN and checks every URL against the
Web Risk API. The Safe or Unsafe verdict is printed to STDOUT. If an error
occurred, debug information may be printed to STDERR.

Exit codes (bitwise OR of following codes):
  0  if and only if all URLs were looked up and are safe.
  1  if at least one URL is not safe.
  2  if at least one URL lookup failed.
  4  if the input was invalid.

Usage: %s -apikey=$APIKEY

`

const (
	codeSafe = (1 << iota) / 2 // Sequence of 0, 1, 2, 4, 8, etc...
	codeUnsafe
	codeFailed
	codeInvalid
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if *apiKeyFlag == "" {
		fmt.Fprintln(os.Stderr, "No -apikey specified")
		os.Exit(codeInvalid)
	}
	sb, err := webrisk.NewWebriskClient(webrisk.Config{
		APIKey:    *apiKeyFlag,
		DBPath:    *databaseFlag,
		Logger:    os.Stderr,
		ServerURL: *serverURLFlag,
		ProxyURL:  *proxyFlag,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to initialize Web Risk client: ", err)
		os.Exit(codeInvalid)
	}

	scanner := bufio.NewScanner(os.Stdin)
	code := codeSafe
	for scanner.Scan() {
		url := scanner.Text()
		threats, err := sb.LookupURLs([]string{url})
		if err != nil {
			fmt.Fprintln(os.Stdout, "Unknown URL:", url)
			fmt.Fprintln(os.Stderr, "Lookup error:", err)
			code |= codeFailed
		} else if len(threats[0]) == 0 {
			fmt.Fprintln(os.Stdout, "Safe URL:", url)
		} else {
			fmt.Fprintln(os.Stdout, "Unsafe URL:", threats[0])
			code |= codeUnsafe
		}
	}
	if scanner.Err() != nil {
		fmt.Fprintln(os.Stderr, "Unable to read input:", scanner.Err())
		code |= codeInvalid
	}
	os.Exit(code)
}
