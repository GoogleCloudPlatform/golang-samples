// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
