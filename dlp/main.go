/*
Copyright 2018 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	l dlppb.Likelihood
}

func (m *minLikelihoodFlag) String() string {
	return fmt.Sprint(m.l)
}

func (m *minLikelihoodFlag) Set(s string) error {
	l, ok := dlppb.Likelihood_value[s]
	if !ok {
		return fmt.Errorf("not a valid likelihood: %q", s)
	}
	m.l = dlppb.Likelihood(l)
	return nil
}

func minLikelihoodValues() string {
	var s []string
	for _, m := range dlppb.Likelihood_name {
		s = append(s, m)
	}
	return strings.Join(s, ", ")
}

type bytesTypeFlag struct {
	bt dlppb.ByteContentItem_BytesType
}

func (f *bytesTypeFlag) String() string {
	return fmt.Sprint(f.bt)
}

func (f *bytesTypeFlag) Set(s string) error {
	b, ok := dlppb.ByteContentItem_BytesType_value[s]
	if !ok {
		return fmt.Errorf("not a valid BytesType: %q", s)
	}
	f.bt = dlppb.ByteContentItem_BytesType(b)
	return nil
}

func bytesTypeValues() string {
	var s []string
	for _, m := range dlppb.ByteContentItem_BytesType_name {
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

	project := flag.String("project", "", "GCloud project ID (required)")
	languageCode := flag.String("languageCode", "en-US", "Language code for infoTypes")
	maxFindings := flag.Int("maxFindings", 0, "Number of results for inspect (default 0 (no limit))")
	includeQuote := flag.Bool("includeQuote", false, "Include a quote of findings for inspect (default false)")
	infoTypesString := flag.String("infoTypes", "PHONE_NUMBER,EMAIL_ADDRESS,CREDIT_CARD_NUMBER,US_SOCIAL_SECURITY_NUMBER", "Info types to inspect or redact")

	var minLikelihood minLikelihoodFlag
	flag.Var(&minLikelihood, "minLikelihood", fmt.Sprintf("Minimum likelihood value [%v] (default %v)", minLikelihoodValues(), dlppb.Likelihood_name[0]))

	var bytesType bytesTypeFlag
	flag.Var(&bytesType, "bytesType", fmt.Sprintf("Bytes type of input file [%v] (default %v)", bytesTypeValues(), dlppb.ByteContentItem_BytesType_name[0]))
	flag.Parse()

	infoTypesList := strings.Split(*infoTypesString, ",")

	if *project == "" {
		flag.Usage()
		os.Exit(1)
	}

	switch flag.Arg(0) {
	case "inspect":
		inspect(os.Stdout, client, minLikelihood.l, int32(*maxFindings), *includeQuote, infoTypesList, *project, flag.Arg(1))
	case "inspectFile":
		inspectFile(os.Stdout, client, minLikelihood.l, int32(*maxFindings), *includeQuote, infoTypesList, *project, bytesType.bt, flag.Arg(1))
	case "inspectGCSFile":
		inspectGCSFile(os.Stdout, client, minLikelihood.l, int32(*maxFindings), *includeQuote, infoTypesList, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4))
	case "inspectDatastore":
		inspectGCSFile(os.Stdout, client, minLikelihood.l, int32(*maxFindings), *includeQuote, infoTypesList, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4))
	case "inspectBigquery":
		inspectBigquery(os.Stdout, client, minLikelihood.l, int32(*maxFindings), *includeQuote, infoTypesList, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4), flag.Arg(5))
	case "redactImage":
		redactImage(os.Stdout, client, minLikelihood.l, infoTypesList, *project, bytesType.bt, flag.Arg(1), flag.Arg(2))
	case "infoTypes":
		infoTypes(os.Stdout, client, *languageCode, flag.Arg(1))
	case "mask":
		mask(os.Stdout, client, *project, flag.Arg(1), "*", 0)
	case "dateShift":
		deidentifyDateShift(os.Stdout, client, *project, -2000, 2000, flag.Arg(1))
	case "fpe":
		deidentifyFPE(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4))
	case "reidentifyFPE":
		reidentifyFPE(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4))
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
	case "createTrigger":
		createTrigger(os.Stdout, client, minLikelihood.l, int32(*maxFindings), *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4), 12, infoTypesList)
	case "listTriggers":
		listTriggers(os.Stdout, client, *project)
	case "deleteTrigger":
		deleteTrigger(os.Stdout, client, flag.Arg(1))
	case "createInspectTemplate":
		createInspectTemplate(os.Stdout, client, minLikelihood.l, int32(*maxFindings), *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), infoTypesList)
	case "listInspectTemplates":
		listInspectTemplates(os.Stdout, client, *project)
	case "deleteInspectTemplate":
		deleteInspectTemplate(os.Stdout, client, flag.Arg(1))
	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand: %q\n\nUsage: %s CMD\n", flag.Arg(0), os.Args[0])
		os.Exit(1)
	}
}
