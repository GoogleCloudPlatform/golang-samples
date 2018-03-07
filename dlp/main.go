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
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

type minLikelihoodFlag struct {
	likelihood dlppb.Likelihood
}

func (m *minLikelihoodFlag) String() string {
	return fmt.Sprint(m.likelihood)
}

func (m *minLikelihoodFlag) Set(s string) error {
	l, ok := dlppb.Likelihood_value[s]
	if !ok {
		return fmt.Errorf("not a valid likelihood: %q", s)
	}
	m.likelihood = dlppb.Likelihood(l)
	return nil
}

func minLikelihoodValues() string {
	var s []string
	for _, m := range dlppb.Likelihood_name {
		s = append(s, m)
	}
	return strings.Join(s, ", ")
}

func main() {
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	project := flag.String("project", "", "GCloud project ID")
	languageCode := flag.String("languageCode", "en-US", "Language code (for infoTypes)")
	maxFindings := flag.Int("maxFindings", 0, "Number of results (for inspect)")
	includeQuote := flag.Bool("includeQuote", false, "Include a quote of findings (for inspect)")

	var minLikelihood minLikelihoodFlag
	flag.Var(&minLikelihood, "minLikelihood", fmt.Sprintf("Minimum likelihood value (%v).", minLikelihoodValues()))
	flag.Parse()

	if *project == "" {
		flag.Usage()
		os.Exit(1)
	}

	switch flag.Arg(0) {
	case "inspect":
		inspect(os.Stdout, client, minLikelihood.likelihood, int32(*maxFindings), *includeQuote, *project, flag.Arg(1))
	case "redact":
		redact(os.Stdout, client, minLikelihood.likelihood, *project, flag.Arg(1))
	case "infoTypes":
		infoTypes(os.Stdout, client, *languageCode, flag.Arg(1))
	case "mask":
		mask(os.Stdout, client, *project, flag.Arg(1))
	case "fpe":
		deidentifyFPE(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3))
	case "riskNumerical":
		riskNumerical(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4), flag.Arg(5), flag.Arg(6))
	case "riskCategorical":
		riskCategorical(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4), flag.Arg(5), flag.Arg(6))
	case "riskKAnonymity":
		riskKAnonymity(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4), flag.Arg(5), strings.Split(flag.Arg(6), ",")...)
	case "riskLDiversity":
		riskLDiversity(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4), flag.Arg(5), flag.Arg(6), strings.Split(flag.Arg(7), ",")...)
	case "riskKMap":
		riskKMap(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4), flag.Arg(5), flag.Arg(6), strings.Split(flag.Arg(7), ",")...)
	default:
		fmt.Fprintf(os.Stderr, `Usage: %s CMD "string"\n`, os.Args[0])
		os.Exit(1)
	}
}
