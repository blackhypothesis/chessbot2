package main

import (
	"fmt"
	"log"
	"time"

	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
)

func getEngineBestMove(game *chess.Game, eng *uci.Engine, move_list []string) *chess.Move {
	// defer TimeTrack(time.Now())
	for _, move := range move_list {
		if err := game.MoveStr(move); err != nil {
			log.Println("Loading moves: ", err)
		}
	}

	// setoption name Threads value 8
	cmdThreads := uci.CmdSetOption{
		Name:  "Threads",
		Value: "8",
	}

	depth := 22
	if len(move_list) > 60 {
		depth = 18
	}
	cmdPos := uci.CmdPosition{Position: game.Position()}
	cmdGo := uci.CmdGo{
		Depth:    depth,
		MoveTime: 1 * time.Second,
	}

	fmt.Println("cmdPos: ", cmdPos)
	if err := eng.Run(cmdThreads, cmdPos, cmdGo); err != nil {
		panic(err)
	}
	search_resultes := eng.SearchResults()
	move := search_resultes.BestMove

	pv_len := len(search_resultes.Info.PV)
	if pv_len > 10 {
		pv_len = 10
	}
	log.Println("Best Move:     ", move)
	log.Println("Info: Depth:   ", search_resultes.Info.Depth)
	log.Println("Info: Score:   ", search_resultes.Info.Score.CP)
	log.Println("Info: PV:      ", search_resultes.Info.PV[:pv_len])
	log.Println("Info: NPS:     ", search_resultes.Info.NPS)
	log.Println("Info: Time:    ", search_resultes.Info.Time)
	log.Println("-----------------------------------------------------------")

	return move
}
