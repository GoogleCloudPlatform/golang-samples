// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// dlp is an example of using the DLP API.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"golang.org/x/net/context"

	dlp "cloud.google.com/go/dlp/apiv2"
	"github.com/fatih/color"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

func minLikelihoodValues() string {
	var s []string
	for _, m := range dlppb.Likelihood_name {
		s = append(s, m)
	}
	return strings.Join(s, ", ")
}

func bytesTypeValues() string {
	var s []string
	for _, m := range dlppb.ByteContentItem_BytesType_name {
		s = append(s, m)
	}
	return strings.Join(s, ", ")
}

var (
	project                = flag.String("project", "", "Project ID (required)")
	languageCode           = flag.String("languageCode", "en-US", "Language code for infoTypes")
	infoTypesString        = flag.String("infoTypes", "PHONE_NUMBER,EMAIL_ADDRESS,CREDIT_CARD_NUMBER,US_SOCIAL_SECURITY_NUMBER", "Info types to inspect*, redactImage, createTrigger, and createInspectTemplate")
	customDictionaryString = flag.String("customDictionary", "", "Custom dictionary for inspect*")
	customRegexString      = flag.String("customRegex", "", "Custom regex for inspect*")
	minLikelihoodFlag      = flag.String("minLikelihood", "LIKELIHOOD_UNSPECIFIED", fmt.Sprintf("Minimum likelihood value for inspect*, redactImage, createTrigger, and createInspectTemplate [%v]", minLikelihoodValues()))
	bytesTypeFlag          = flag.String("bytesType", "BYTES_TYPE_UNSPECIFIED", fmt.Sprintf("Bytes type of input file for inspectFile and redactImage [%v]", bytesTypeValues()))
	maxFindings            = flag.Int("maxFindings", 0, "Number of results for inspect*, createTrigger, and createInspectTemplate (default 0 (no limit))")
	autoPopulateTimespan   = flag.Bool("autoPopulateTimespan", false, "Limit scan to new content only (default false)")
	includeQuote           = flag.Bool("includeQuote", false, "Include a quote of findings for inspect* (default false)")
)

func main() {
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	flag.Parse()

	infoTypesList := strings.Split(*infoTypesString, ",")
	var customDictionariesList []string
	if *customDictionaryString != "" {
		customDictionariesList = []string{*customDictionaryString}
	}
	var customRegexesList []string
	if *customRegexString != "" {
		customRegexesList = []string{*customRegexString}
	}

	if *project == "" {
		fmt.Fprintf(os.Stderr, "Must provide a -project\n\n")
		flag.Usage()
		os.Exit(1)
	}

	minLikelihoodInt, ok := dlppb.Likelihood_value[*minLikelihoodFlag]
	if !ok {
		fmt.Fprintf(os.Stderr, "not a valid likelihood: %q\n", *minLikelihoodFlag)
		fmt.Fprintf(os.Stderr, "must be one of [%s]\n", minLikelihoodValues())
		os.Exit(1)
	}
	minLikelihood := dlppb.Likelihood(minLikelihoodInt)

	bytesTypeInt, ok := dlppb.ByteContentItem_BytesType_value[*bytesTypeFlag]
	if !ok {
		fmt.Fprintf(os.Stderr, "not a valid BytesType: %q\n", *minLikelihoodFlag)
		fmt.Fprintf(os.Stderr, "must be one of [%s]\n", bytesTypeValues())
		os.Exit(1)
	}
	bytesType := dlppb.ByteContentItem_BytesType(bytesTypeInt)

	switch flag.Arg(0) {
	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand: %q\n\n", flag.Arg(0))
		flag.Usage()
		os.Exit(1)
	case "inspect":
		checkNArg(1)
		inspectString(os.Stdout, client, *project, minLikelihood, int32(*maxFindings), *includeQuote, infoTypesList, customDictionariesList, customRegexesList, flag.Arg(1))
	case "inspectFile":
		checkNArg(1)
		f, err := os.Open(flag.Arg(1))
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		inspectFile(os.Stdout, client, *project, minLikelihood, int32(*maxFindings), *includeQuote, infoTypesList, customDictionariesList, customRegexesList, bytesType, f)
	case "inspectGCSFile":
		checkNArg(4)
		inspectGCSFile(os.Stdout, client, *project, minLikelihood, int32(*maxFindings), *includeQuote, infoTypesList, customDictionariesList, customRegexesList, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4))
	case "inspectDatastore":
		checkNArg(5)
		inspectDatastore(os.Stdout, client, *project, minLikelihood, int32(*maxFindings), *includeQuote, infoTypesList, customDictionariesList, customRegexesList, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4), flag.Arg(5))
	case "inspectBigquery":
		checkNArg(5)
		inspectBigquery(os.Stdout, client, *project, minLikelihood, int32(*maxFindings), *includeQuote, infoTypesList, customDictionariesList, customRegexesList, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4), flag.Arg(5))

	case "redactImage":
		checkNArg(2)
		redactImage(os.Stdout, client, *project, minLikelihood, infoTypesList, bytesType, flag.Arg(1), flag.Arg(2))

	case "infoTypes":
		checkNArg(1)
		infoTypes(os.Stdout, client, *languageCode, flag.Arg(1))

	case "mask":
		checkNArg(1)
		mask(os.Stdout, client, *project, flag.Arg(1), infoTypesList, "*", 0)
	case "dateShift":
		checkNArg(1)
		deidentifyDateShift(os.Stdout, client, *project, -2000, 2000, flag.Arg(1))
	case "fpe":
		checkNArg(4)
		deidentifyFPE(os.Stdout, client, *project, flag.Arg(1), infoTypesList, flag.Arg(2), flag.Arg(3), flag.Arg(4))
	case "reidentifyFPE":
		checkNArg(4)
		reidentifyFPE(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4))

	case "riskNumerical":
		checkNArg(6)
		riskNumerical(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4), flag.Arg(5), flag.Arg(6))
	case "riskCategorical":
		checkNArg(6)
		riskCategorical(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4), flag.Arg(5), flag.Arg(6))
	case "riskKAnonymity":
		checkNArg(6)
		riskKAnonymity(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4), flag.Arg(5), strings.Split(flag.Arg(6), ",")...)
	case "riskLDiversity":
		checkNArg(7)
		riskLDiversity(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4), flag.Arg(5), flag.Arg(6), strings.Split(flag.Arg(7), ",")...)
	case "riskKMap":
		checkNArg(7)
		riskKMap(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4), flag.Arg(5), flag.Arg(6), strings.Split(flag.Arg(7), ",")...)

	case "createTrigger":
		checkNArg(4)
		createTrigger(os.Stdout, client, *project, minLikelihood, int32(*maxFindings), flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4), *autoPopulateTimespan, 12, infoTypesList)
	case "listTriggers":
		checkNArg(0)
		listTriggers(os.Stdout, client, *project)
	case "deleteTrigger":
		checkNArg(1)
		deleteTrigger(os.Stdout, client, flag.Arg(1))

	case "createInspectTemplate":
		checkNArg(3)
		createInspectTemplate(os.Stdout, client, *project, minLikelihood, int32(*maxFindings), flag.Arg(1), flag.Arg(2), flag.Arg(3), infoTypesList)
	case "listInspectTemplates":
		checkNArg(0)
		listInspectTemplates(os.Stdout, client, *project)
	case "deleteInspectTemplate":
		checkNArg(1)
		deleteInspectTemplate(os.Stdout, client, flag.Arg(1))

	case "listJobs":
		checkNArg(2)
		listJobs(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2))
	case "deleteJob":
		checkNArg(1)
		deleteJob(os.Stdout, client, flag.Arg(1))
	}

}

func sortedKeys(m map[string]string) []string {
	var l []string
	for k := range m {
		l = append(l, k)
	}
	sort.Strings(l)
	return l
}

func sortedMapKeys(m map[string]map[string]string) []string {
	var l []string
	for k := range m {
		l = append(l, k)
	}
	sort.Strings(l)
	return l
}

var subcommands = map[string]map[string]string{
	"inspect": {
		"inspect":          "<string>",
		"inspectFile":      "<filename>",
		"inspectGCSFile":   "<pubSubTopic> <pubSubSub> <bucketName> <fileName> ",
		"inspectDatastore": "<pubSubTopic> <pubSubSub> <dataProject> <namespaceID> <kind>",
		"inspectBigquery":  "<pubSubTopic> <pubSubSub> <dataProject> <datasetID> <tableID>",
	},
	"redact": {
		"redactImage": "<inputPath> <outputPath>",
	},
	"metadata": {
		"infoTypes": "<filter>",
	},
	"deidentify": {
		"mask":          "<string>",
		"dateShift":     "<string>",
		"fpe":           "<string> <wrappedKeyFileName> <cryptoKeyname> <surrogateInfoType>",
		"reidentifyFPE": "<string> <wrappedKeyFileName> <cryptoKeyname> <surrogateInfoType>",
	},
	"risk": {
		"riskNumerical":   "<dataProject> <pubSubTopic> <pubSubSub> <datasetID> <tableID> <columnName>",
		"riskCategorical": "<dataProject> <pubSubTopic> <pubSubSub> <datasetID> <tableID> <columnName>",
		"riskKAnonymity":  "<dataProject> <pubSubTopic> <pubSubSub> <datasetID> <tableID> <column,names>",
		"riskLDiversity":  "<dataProject> <pubSubTopic> <pubSubSub> <datasetID> <tableID> <sensitiveAttribute> <column,names>",
		"riskKMap":        "<dataProject> <pubSubTopic> <pubSubSub> <datasetID> <tableID> <region> <column,names>",
	},
	"triggers": {
		"createTrigger": "<triggerID> <displayName> <description> <bucketName>",
		"listTriggers":  "",
		"deleteTrigger": "<fullTriggerID>",
	},
	"templates": {
		"createInspectTemplate": "<templateID> <displayName> <description>",
		"listInspectTemplates":  "",
		"deleteInspectTemplate": "<fullTemplateID>",
	},
	"jobs": {
		"listJobs":  "<filter> <jobType>",
		"deleteJob": "<jobID>",
	},
}

func init() {
	bold := color.New(color.Bold).FprintfFunc()
	blue := color.New(color.FgHiBlue, color.Bold).FprintfFunc()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -project <project> subcommand [args]\n\n", os.Args[0])
		bold(os.Stderr, "Subcommands:\n")
		for _, c := range sortedMapKeys(subcommands) {
			blue(os.Stderr, "  %s\n", c)
			for _, s := range sortedKeys(subcommands[c]) {
				fmt.Fprintf(os.Stderr, "    %s -project <project> [options] ", os.Args[0])
				bold(os.Stderr, "%s ", s)
				fmt.Fprintf(os.Stderr, "%s\n", subcommands[c][s])
			}
		}
		bold(os.Stderr, "\n\nOptions:\n")
		flag.PrintDefaults()
	}
}

func findSubcommandArgs(s string) string {
	for _, v := range subcommands {
		if a, ok := v[s]; ok {
			return a
		}
	}
	return ""
}

// checkNArg ensures there are n arguments after the subcommand.
func checkNArg(n int) {
	if flag.NArg()-1 != n {
		fmt.Fprintf(os.Stderr, "Error:  found %d args, expected %d\n", flag.NArg()-1, n)
		fmt.Fprintf(os.Stderr, "Usage: %s -project %s [options] %s %s\n\n", os.Args[0], flag.Lookup("project").Value, flag.Arg(0), findSubcommandArgs(flag.Arg(0)))
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
}
