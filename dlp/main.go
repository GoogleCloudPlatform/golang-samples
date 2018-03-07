// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// dlp is an example of using the DLP API.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	dlp "cloud.google.com/go/dlp/apiv2"
)

func main() {
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	project := flag.String("project", "", "GCloud project ID")
	flag.Parse()

	if *project == "" {
		flag.Usage()
		os.Exit(1)
	}

	switch flag.Arg(0) {
	case "inspect":
		inspect(os.Stdout, client, *project, flag.Arg(1))
	case "redact":
		redact(os.Stdout, client, *project, flag.Arg(1))
	case "infoTypes":
		infoTypes(os.Stdout, client, flag.Arg(1))
	case "mask":
		mask(os.Stdout, client, *project, flag.Arg(1))
	case "deidfpe":
		deidentifyFPE(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3))
	case "riskNumerical":
		// For example:
		// dlp -project my-project riskNumerical bigquery-public-data risk-topic risk-sub nhtsa_traffic_fatalities accident_2015 state_number
		riskNumerical(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4), flag.Arg(5), flag.Arg(6))
	case "riskCategorical":
		// For example:
		// dlp -project my-project riskCategorical bigquery-public-data risk-topic risk-sub nhtsa_traffic_fatalities accident_2015 state_number
		riskCategorical(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4), flag.Arg(5), flag.Arg(6))
	case "riskKAnonymity":
		// For example:
		// dlp -project my-project riskKAnonymity bigquery-public-data risk-topic risk-sub nhtsa_traffic_fatalities accident_2015 state_number,county
		riskKAnonymity(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4), flag.Arg(5), strings.Split(flag.Arg(6), ",")...)
	case "riskLDiversity":
		// For example:
		// dlp -project my-project riskLDiversity bigquery-public-data risk-topic risk-sub nhtsa_traffic_fatalities accident_2015 city state_number,county
		riskLDiversity(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4), flag.Arg(5), flag.Arg(6), strings.Split(flag.Arg(7), ",")...)
	case "riskKMap":
		// For example:
		// dlp -project my-project kMap bigquery-public-data san_francisco bikeshare_trips risk-topic risk-sub US_ZIP_5 zip_code
		riskKMap(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4), flag.Arg(5), flag.Arg(6), strings.Split(flag.Arg(7), ",")...)
	default:
		fmt.Fprintf(os.Stderr, `Usage: %s CMD "string"\n`, os.Args[0])
		os.Exit(1)
	}
}
