// Copyright 2024 Google LLC
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

// Mock objects for shared project

package snippets

import (
	"context"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	gax "github.com/googleapis/gax-go/v2"
)

type ClientInterface interface {
	Close() error
	Delete(context.Context, *computepb.DeleteReservationRequest, ...gax.CallOption) (*compute.Operation, error)
	Insert(context.Context, *computepb.InsertReservationRequest, ...gax.CallOption) (*compute.Operation, error)
}

type ReservationsClient struct{}

func (client ReservationsClient) Close() error {
	return nil
}

func (client ReservationsClient) Insert(context.Context, *computepb.InsertReservationRequest, ...gax.CallOption) (*compute.Operation, error) {
	return nil, nil
}

func (client ReservationsClient) Delete(context.Context, *computepb.DeleteReservationRequest, ...gax.CallOption) (*compute.Operation, error) {
	return nil, nil
}
