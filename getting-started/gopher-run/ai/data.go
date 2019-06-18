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

//This package moves the Firestore play data into data/playdata.csv
package data

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func main() {
	fmt.Println("hello")
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("firestore.NewClient: %v", err)
	}
	defer client.Close()
	f, err := os.Create("data/playdata.csv")
	if err != nil {
		log.Fatalf("os.Create: %v", err)
	}
	defer f.Close()

	iter := client.Collection("patterns").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed iteration %v", err)
		}
		pos := ""
		vals := []string{"h", "x0", "y0", "x1", "y1", "x2", "y2", "x3", "y3"}
		for _, v := range vals {
			fl, ok := doc.Data()[v].(float64)
			if ok {
				pos += "," + strconv.FormatFloat(fl, 'f', 6, 64)
			} else {
				fmt.Println(doc.Data()[v])
				continue
			}
		}
		a, ok := doc.Data()["act"].(string)
		if ok {
			f.WriteString(a + pos + "\n")
		}
	}
	f.Sync()
}
