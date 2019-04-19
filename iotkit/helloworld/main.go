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
