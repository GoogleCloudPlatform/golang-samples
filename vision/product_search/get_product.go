// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

// [START imports]
import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	
	"context"

	vision "cloud.google.com/go/vision/apiv1"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)
// [END imports]

// [START vision_product_search_get_product]
func getProduct(project string, location string, productId string) error {
	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return err
	}

	req := &visionpb.GetProductRequest{}
	req.Name = fmt.Sprintf("projects/%s/locations/%s/products/%s", project, location, productId)

	resp, err := c.GetProduct(ctx, req)
	if err != nil {
		return err
	}

	fmt.Println("Product name:", resp.Name)
	fmt.Println("Product display name:", resp.DisplayName)
	fmt.Println("Product category:", resp.ProductCategory)
	fmt.Println("Product labels:", resp.ProductLabels)

	return nil
}
// [END vision_product_search_get_product]

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <project-id> <location> <product-id>\n", filepath.Base(os.Args[0]))
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 3 {
		flag.Usage()
		os.Exit(1)
	}

	err := getProduct(args[0], args[1], args[2])

	if err != nil {
		fmt.Println("Error:", err)
	}
}
