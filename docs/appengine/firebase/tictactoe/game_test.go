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

package tictactoe

import "testing"

func TestCheckWin(t *testing.T) {
	tests := []struct {
		board  string
		winner string
	}{
		{
			board: "" +
				"X X" +
				"XOX" +
				"OOO",
			winner: "O",
		},
		{
			board: "" +
				"X O" +
				"XOX" +
				"XOO",
			winner: "X",
		},
		{
			board: "" +
				"X O" +
				" OX" +
				"XOO",
			winner: "",
		},
		{
			board:  "         ",
			winner: "",
		},
	}

	for _, tt := range tests {
		g := &Game{Board: tt.board}
		winner, _ := g.CheckWin()
		if winner != tt.winner {
			t.Errorf("CheckWin(%s): got %q, want %q", tt.board, winner, tt.winner)
		}
	}
}
