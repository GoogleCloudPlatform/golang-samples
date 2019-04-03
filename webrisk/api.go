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
package webrisk

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	pb "github.com/GoogleCloudPlatform/golang-samples/webrisk/internal/webrisk_proto"
	"github.com/golang/protobuf/proto"
)

const (
	findHashPath    = "/$rpc/google.cloud.webrisk.v1beta1.WebRiskServiceV1Beta1/SearchHashes"
	fetchUpdatePath = "/$rpc/google.cloud.webrisk.v1beta1.WebRiskServiceV1Beta1/ComputeThreatListDiff"
)

// The api interface specifies wrappers around the Web Risk API.
type api interface {
	ListUpdate(ctx context.Context, req *pb.ComputeThreatListDiffRequest) (*pb.ComputeThreatListDiffResponse, error)
	HashLookup(ctx context.Context, req *pb.SearchHashesRequest) (*pb.SearchHashesResponse, error)
}

// netAPI is an api object that talks to the server over HTTP.
type netAPI struct {
	client *http.Client
	url    *url.URL
}

// newNetAPI creates a new netAPI object pointed at the provided root URL.
// For every request, it will use the provided API key.
// If a proxy URL is given, it will be used in place of the default $HTTP_PROXY.
// If the protocol is not specified in root, then this defaults to using HTTPS.
func newNetAPI(root string, key string, proxy string) (*netAPI, error) {
	if !strings.Contains(root, "://") {
		root = "https://" + root
	}
	u, err := url.Parse(root)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{}

	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			return nil, err
		}
		httpClient = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	}

	q := u.Query()
	q.Set("key", key)
	q.Set("alt", "proto")
	u.RawQuery = q.Encode()
	return &netAPI{url: u, client: httpClient}, nil
}

// doRequests performs a POST to requestPath. It uses the marshaled form of req
// as the request body payload, and automatically unmarshals the response body
// payload as resp.
func (a *netAPI) doRequest(ctx context.Context, requestPath string, req proto.Message, resp proto.Message) error {
	p, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	u := *a.url // Make a copy of URL
	u.Path = requestPath
	httpReq, err := http.NewRequest("POST", u.String(), bytes.NewReader(p))
	httpReq.Header.Add("Content-Type", "application/x-protobuf")
	httpReq = httpReq.WithContext(ctx)
	httpResp, err := a.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode != 200 {
		return fmt.Errorf("webrisk: unexpected server response code: %d", httpResp.StatusCode)
	}
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}
	return proto.Unmarshal(body, resp)
}

// ListUpdate issues a ComputeThreatListDiff API call and returns the response.
func (a *netAPI) ListUpdate(ctx context.Context, req *pb.ComputeThreatListDiffRequest) (*pb.ComputeThreatListDiffResponse, error) {
	resp := new(pb.ComputeThreatListDiffResponse)
	return resp, a.doRequest(ctx, fetchUpdatePath, req, resp)
}

// HashLookup issues a SearchHashes API call and returns the response.
func (a *netAPI) HashLookup(ctx context.Context, req *pb.SearchHashesRequest) (*pb.SearchHashesResponse, error) {
	resp := new(pb.SearchHashesResponse)
	return resp, a.doRequest(ctx, findHashPath, req, resp)
}
