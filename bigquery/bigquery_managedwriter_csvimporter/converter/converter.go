// Copyright 2022 Google LLC
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

// Package converter provides a simple example of validating and converting
// CSV data into an appropriate protocol buffer form suitable by the BigQuery
// storage write API.
package converter

import (
	"fmt"

	"cloud.google.com/go/bigquery/storage/managedwriter/adapt"
	storagepb "google.golang.org/genproto/googleapis/cloud/bigquery/storage/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

var msgRoot = "root"

// CSVConverter is a simplistic example of a converter that can rewrite CSV data
// into a protocol buffer message compatible with a given table's schema.
type CSVConverter struct {
	schema        *storagepb.TableSchema
	msgDescriptor protoreflect.MessageDescriptor
}

// NewCSVConverter instantiates a new converter.
func NewCSVConverter(schema *storagepb.TableSchema) (*CSVConverter, error) {
	if schema == nil {
		return nil, fmt.Errorf("no input schema provided")
	}
	desc, err := adapt.StorageSchemaToProto2Descriptor(schema, msgRoot)
	if err != nil {
		return nil, err
	}
	md, ok := desc.(protoreflect.MessageDescriptor)
	if !ok {
		return nil, fmt.Errorf("couldn't convert descriptor (type %T) to message descriptor.", desc)
	}

	return &CSVConverter{
		schema:        schema,
		msgDescriptor: md,
	}, nil

}

// ValidateColumns ensures that the CSV headers can be mapped to the
// table by comparing the schema.  It also identifies the case where
// the CSV data does not include fields that the table schema requires.
func (conv *CSVConverter) Validate(colNames []string) error {

	verifiedCols := make(map[string]bool)
	for _, col := range colNames {
		found := false
		for _, f := range conv.schema.GetFields() {
			if f.GetName() == col {
				if f.GetType() == storagepb.TableFieldSchema_STRING {
					verifiedCols[col] = true
					found = true
				} else {
					return fmt.Errorf("table column %s has incompatible type %s", col, f.GetType().String())
				}
			}
		}
		if !found {
			return fmt.Errorf("data column %s not present in table schema", col)
		}
	}
	// final check: ensure no required fields that aren't part of csv schema
	// exist.
	for _, f := range conv.schema.GetFields() {
		if f.GetMode() == storagepb.TableFieldSchema_REQUIRED {
			name := f.GetName()
			if _, ok := verifiedCols[name]; !ok {
				return fmt.Errorf("table column %s is required, but not present in provided headers", name)
			}
		}
	}
	return nil
}

// Convert serializes the row data into a protobuf message.
func (conv *CSVConverter) Convert(data map[string]string) ([]byte, error) {
	msg := dynamicpb.NewMessage(conv.msgDescriptor)

	for name, val := range data {
		fd := msg.Descriptor().Fields().ByName(protoreflect.Name(name))
		if fd == nil {
			return nil, fmt.Errorf("message didn't have expected field %s", name)
		}
		msg.Set(fd, protoreflect.ValueOf(val))
	}

	return proto.Marshal(msg)
}

// ProtoSchema produces the protocol buffer schema
func (conv *CSVConverter) ProtoSchema() (*descriptorpb.DescriptorProto, error) {
	return adapt.NormalizeDescriptor(conv.msgDescriptor)
}
