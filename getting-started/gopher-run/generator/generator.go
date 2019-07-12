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

// Package generator returns procedurally generated level for Gopher Run
package generator

import (
	"fmt"
	"math/rand"
)

type RequestData struct {
	Xmin  float64
	Xmax  float64
	Speed float64
}

// Vector3 is a position vector
type Vector3 struct {
	X float64
	Y float64
	Z float64
}

func V3Add(u, v Vector3) Vector3 {
	return Vector3{u.X + v.X, u.Y + v.Y, u.Z + v.Z}
}

type Transform struct {
	position   Vector3
	localScale Vector3
}

type GameObject struct {
	name      string
	transform Transform
}

var Speed float64
var xmax float64
var objects []GameObject

const ymax = 100

func randRange(i float64, j float64) float64 {
	return rand.Float64()*(j-i) + i
}

func GetObjects() []GameObject {
	return objects
}
func (o GameObject) ToString() string {
	return fmt.Sprintf("%v %v %v %v %v %v %v", o.name, o.transform.position.X, o.transform.position.Y, o.transform.position.Z, o.transform.localScale.X, o.transform.localScale.Y, o.transform.localScale.Z)
}
func GenerateBackground(start Vector3, end Vector3) {
	objects = []GameObject{}
	for n := start.X; n < end.X; n += 15 {
		for m := 0; m < 3; m++ {
			cscale := randRange(0.2, 0.6)
			cx := randRange(n, n+10)
			cy := randRange(10, 25)
			cz := 15 + randRange(-5, 5)
			cloud(Vector3{cx, cy, cz}, cscale)
		}
		scale := randRange(1.5, 2.5)
		x := randRange(n, n+10)
		y := 5.0
		z := 15 + randRange(-5, 5)
		hill(Vector3{x, y, z}, scale)
	}
}

func cloud(pos Vector3, scale float64) Vector3 {
	instantiate("cloud", pos, scale)
	return pos
}

func hill(pos Vector3, scale float64) Vector3 {
	instantiate("hill", pos, scale)
	return pos
}

func instantiate(s string, pos Vector3, scale float64) GameObject {
	o := GameObject{s, Transform{pos, Vector3{scale, scale, scale}}}
	objects = append(objects, o)
	return o
}
