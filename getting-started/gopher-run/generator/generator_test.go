// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package generator

import "testing"

func TestGenerateBackground(t *testing.T) {
	objects := GenerateBackground(0, 500, 16)
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
