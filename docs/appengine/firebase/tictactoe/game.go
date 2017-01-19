// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package tictactoe

import (
	"errors"

	"google.golang.org/appengine/datastore"
)

type Game struct {
	K            *datastore.Key `json:"-" datastore:"__key__"`
	UserX        string         `json:"userX"`
	UserO        string         `json:"userO"`
	Board        string         `json:"board"`
	MoveX        bool           `json:"moveX"`
	Winner       string         `json:"winner"`
	WinningBoard string         `json:"winningBoard"`
}

func NewGame() *Game {
	g := Game{}
	g.Board = "         " // 9 spaces
	g.MoveX = true
	return &g
}

// CheckWin returns "X" or "O", depending on who won. It will be empty if the game was a draw.
func (g *Game) CheckWin() (winner string, gameOver bool) {
	// Check horizontal/vertical
	for i := 0; i < 3; i++ {
		if g.Board[i+0] != ' ' && g.Board[i+0] == g.Board[i+3] && g.Board[i+3] == g.Board[i+6] {
			return string(g.Board[i+0]), true
		}
		j := i * 3
		if g.Board[j+0] != ' ' && g.Board[j+0] == g.Board[j+1] && g.Board[j+1] == g.Board[j+2] {
			return string(g.Board[j+0]), true
		}
	}
	// Check diagonals
	if g.Board[0] != ' ' && g.Board[0] == g.Board[4] && g.Board[4] == g.Board[8] {
		return string(g.Board[0]), true
	}
	if g.Board[2] != ' ' && g.Board[2] == g.Board[4] && g.Board[4] == g.Board[6] {
		return string(g.Board[2]), true
	}

	// Check draw
	for _, c := range g.Board {
		if c == ' ' {
			return "", false
		}
	}

	return "Draw", true
}

// MoveAt plays a move at the specified index.
// Input is assumed to be valid.
func (g *Game) MoveAt(index int) error {
	if g.Board[index] != ' ' {
		return errors.New("Not an empty space")
	}
	player := "X"
	if !g.MoveX {
		player = "O"
	}
	g.Board = g.Board[0:index] + player + g.Board[index+1:len(g.Board)]
	return nil
}
