// Copyright 2020 Google LLC
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

// Package shakesapp defines a server which can be queried to determined
// many times a string appears in the works of Shakespeare, and a client
// which can be used to send load to that server.
package shakesapp

import (
	"context"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// server is an implementation of the server for ShakespeareService (defined
// in shakesapp.proto).
type server struct{}

// NewServer returns an implementation of the server for ShakespeareService
// (defined in shakesapp.proto).
func NewServer() ShakespeareServiceServer {
	return &server{}
}

const bucketName = "dataflow-samples"
const bucketPrefix = "shakespeare/"

// GetMatchCount implements a server for ShakespeareService.
func (s *server) GetMatchCount(ctx context.Context, req *ShakespeareRequest) (*ShakespeareResponse, error) {
	resp := &ShakespeareResponse{}
	texts, err := readFiles(ctx, bucketName, bucketPrefix)
	if err != nil {
		return resp, fmt.Errorf("fails to read files: %s", err)
	}
	for _, text := range texts {
		for _, line := range strings.Split(text, "\n") {
			line, query := strings.ToLower(line), strings.ToLower(req.Query)
			// TODO: Compiling and matching a regular expression on every request
			// might be too expensive? Consider optimizing.
			isMatch, err := regexp.MatchString(query, line)
			if err != nil {
				return resp, err
			}
			if isMatch {
				resp.MatchCount++
			}
		}
	}
	return resp, nil
}

// readFiles reads the content of files within the specified bucket with the
// specified prefix path in parallel and returns their content. It fails if
// operations to find or read any of the files fails.
func readFiles(ctx context.Context, bucketName, prefix string) ([]string, error) {
	type resp struct {
		s   string
		err error
	}

	client, err := storage.NewClient(ctx, option.WithoutAuthentication())
	if err != nil {
		return []string{}, fmt.Errorf("failed to create storage client: %s", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)

	var paths []string
	it := bucket.Objects(ctx, &storage.Query{Prefix: bucketPrefix})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []string{}, fmt.Errorf("failed to iterate over files in %s starting with %s: %w", bucketName, prefix, err)
		}
		if attrs.Name != "" {
			paths = append(paths, attrs.Name)
		}
	}

	resps := make(chan resp)
	for _, path := range paths {
		go func(path string) {
			obj := bucket.Object(path)
			r, err := obj.NewReader(ctx)
			if err != nil {
				resps <- resp{"", err}
			}
			defer r.Close()
			data, err := ioutil.ReadAll(r)
			resps <- resp{string(data), err}
		}(path)
	}
	ret := make([]string, len(paths))
	for i := 0; i < len(paths); i++ {
		r := <-resps
		if r.err != nil {
			err = r.err
		}
		ret[i] = r.s
	}
	return ret, err
}
