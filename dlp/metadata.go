package main

import (
	"context"
	"fmt"
	"io"
	"log"

	dlp "cloud.google.com/go/dlp/apiv2"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

func infoTypes(w io.Writer, client *dlp.Client, filter string) {
	rcr := &dlppb.ListInfoTypesRequest{
		Filter: filter,
	}
	r, err := client.ListInfoTypes(context.Background(), rcr)
	if err != nil {
		log.Fatal(err)
	}
	for _, it := range r.GetInfoTypes() {
		fmt.Fprintln(w, it.GetName())
	}
}
