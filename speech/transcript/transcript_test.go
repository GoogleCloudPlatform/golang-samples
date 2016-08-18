// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	speech "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

func TestRecognize(t *testing.T) {
	ctx := context.Background()

	conn, err := transport.DialGRPC(ctx,
		option.WithEndpoint("speech.googleapis.com:443"),
		option.WithScopes("https://www.googleapis.com/auth/cloud-platform"),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	c := speech.NewSpeechClient(conn)

	data, err := ioutil.ReadFile("./quit.raw")
	if err != nil {
		t.Fatal(err)
	}

	rresp, err := recognize(ctx, c, &data)
	if err != nil {
		t.Fatal(err)
	}
	if len(rresp.Responses) < 1 {
		t.Fatal("want recognize responses; got none")
	}

	resp := rresp.Responses[0]
	if resp.Error != nil {
		t.Fatalf("got error from response: %v", err)
	}
	if len(resp.Results) < 1 {
		t.Fatal("want results; got none")
	}
	if len(resp.Results[0].Alternatives) < 1 {
		t.Fatal("want alternatives; got none")
	}
	if got, want := resp.Results[0].Alternatives[0].Transcript, "quit"; got != want {
		t.Errorf("want transcript: %q; got %q", want, got)
	}
}
