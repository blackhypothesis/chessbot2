package main

import (
	"log"
	"time"

	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
)

func playBestMove(game *chess.Game, eng *uci.Engine, move_list []string) *chess.Move {
	defer TimeTrack(time.Now())
	for _, move := range move_list {
		if err := game.MoveStr(move); err != nil {
			log.Println("Loading moves: ", err)
		}
	}
	cmdPos := uci.CmdPosition{Position: game.Position()}
	// cmdGo := uci.CmdGo{MoveTime: time.Second / 4}
	cmdGo := uci.CmdGo{
		Depth:    18,
		MoveTime: 1 * time.Second,
	}
	if err := eng.Run(cmdPos, cmdGo); err != nil {
		panic(err)
	}
	search_resultes := eng.SearchResults()
	move := search_resultes.BestMove
	return move
}
