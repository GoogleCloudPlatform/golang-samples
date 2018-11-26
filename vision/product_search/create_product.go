// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

// [START imports]
import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	
	"context"

	vision "cloud.google.com/go/vision/apiv1"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)
// [END imports]

// [START vision_product_search_create_product]
func createProduct(w io.Writer, projectId string, location string, productId string, productDisplayName string, productCategory string) error {
	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return err
	}

	product := &visionpb.Product{}
	product.DisplayName = productDisplayName
	product.ProductCategory = productCategory

	req := &visionpb.CreateProductRequest{}
	req.Parent = fmt.Sprintf("projects/%s/locations/%s", projectId, location)
	req.ProductId = productId
	req.Product = product

	resp, err := c.CreateProduct(ctx, req)
	if err != nil {
		return err
	}

	fmt.Println("Product name:", resp.Name)

	return nil
}
// [END vision_product_search_create_product]

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <project-id> <location> <product-id> <product-display-name> <product-category>\n", filepath.Base(os.Args[0]))
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 5 {
		flag.Usage()
		os.Exit(1)
	}

	err := createProduct(os.Stdout, args[0], args[1], args[2], args[3], args[4])

	if err != nil {
		fmt.Println("Error:", err)
	}
}
