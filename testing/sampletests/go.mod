module github.com/GoogleCloudPlatform/golang-samples/testing/sampletests

go 1.12

require (
	github.com/google/go-cmp v0.5.2
	github.com/jstemmer/go-junit-report v0.9.1
	golang.org/x/tools v0.0.0-20200904185747-39188db58858
)

// https://github.com/jstemmer/go-junit-report/issues/107
replace github.com/jstemmer/go-junit-report => github.com/tbpg/go-junit-report v0.9.2-0.20200506144438-50086c54f894
