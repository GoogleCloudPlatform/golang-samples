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

// Package generator returns procedurally generated parts of levels for Gopher Run.
package generator

import (
	"fmt"
	"math/rand"
	"testing"
)

type RequestData struct {
	Xmin  float64
	Xmax  float64
	Speed float64
}

// Vector3 is 3-value vector.
type Vector3 struct {
	X float64
	Y float64
	Z float64
}

// Transform is information about the GameObject (corresponds to Unity Transform).
type Transform struct {
	position   Vector3
	localScale Vector3
}

// GameObject corresponds to a Unity GameObject.
type GameObject struct {
	name      string
	transform Transform
}

func TestGenerateBackground(t *testing.T) {
	objects := GenerateBackground(Vector3{0, 0, 0}, Vector3{500, 0, 0}, 16)
	for _, obj := range objects {
		if obj.transform.position.X < 0 {
			t.Errorf("GenerateBackground object got %v, want %v", "negative x", "nonnegative x")
		}
		if obj.transform.position.X > 500 {
			t.Errorf("GenerateBackground object got %v, want %v", "out of range x", "x <= 500")
		}
		if obj.transform.position.Y < 0 {
			t.Errorf("GenerateBackground object got %v, want %v", "negative y", "nonnegative y")
		}
	}
}

func randRange(i float64, j float64) float64 {
	return rand.Float64()*(j-i) + i
}

func (o GameObject) ToString() string {
	return fmt.Sprintf("%v %v %v %v %v %v %v", o.name, o.transform.position.X, o.transform.position.Y, o.transform.position.Z, o.transform.localScale.X, o.transform.localScale.Y, o.transform.localScale.Z)
}

// GenerateBackground determines positions for background objects.
func GenerateBackground(start Vector3, end Vector3, speed float64) []GameObject {
	objects := []GameObject{}
	for n := start.X; n < end.X; n += 15 {
		for m := 0; m < 3; m++ {
			cscale := randRange(0.2, 0.6)
			cx := randRange(n, n+10)
			cy := randRange(10, 25)
			cz := 15 + randRange(-5, 5)
			objects = append(objects, GameObject{"cloud", Transform{Vector3{cx, cy, cz}, Vector3{cscale, cscale, cscale}}})
		}
		scale := randRange(1.5, 2.5)
		x := randRange(n, n+10)
		y := 5.0
		z := 15 + randRange(-5, 5)
		objects = append(objects, GameObject{"hill", Transform{Vector3{x, y, z}, Vector3{scale, scale, scale}}})
	}
	return objects
}
