// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Program helloworld is a tiny program that blinks "Hello gopher!" on a
// Grove LCD RGB Backlight display connected to a BeagleBone.
package main

import (
	"log"
	"time"

	"golang.org/x/exp/io/i2c"

	"github.com/GoogleCloudPlatform/golang-samples/iotkit/helloworld/display"
)

func main() {
	d, err := display.Open(&i2c.Devfs{Dev: "/dev/i2c-2"})
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	// Set the backlight to Go blue.
	if err := d.SetRGB(0, 128, 64); err != nil {
		log.Fatal(err)
	}

	// Blink "Hello gopher!" on the display.
	for {
		if err := d.SetText("Hello gopher!"); err != nil {
			log.Fatal(err)
		}
		time.Sleep(200 * time.Millisecond)
		if err := d.Clear(); err != nil {
			log.Fatal(err)
		}
		time.Sleep(200 * time.Millisecond)
	}
}
