// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var port = flag.Int("port", 8080, "port to listen on")

func main() {
	flag.Parse()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
