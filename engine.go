package main

import (
	"log"
	"time"

	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
)

func getEngineBestMove(game *chess.Game, move_list []string) (*chess.Move, error) {
	// defer TimeTrack(time.Now())

	eng, err := uci.New("stockfish")
	if err != nil {
		return nil, err
	}
	// initialize uci with new game
	if err := eng.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
		return nil, err
	}
	defer eng.Close()

	for _, move := range move_list {
		if err := game.MoveStr(move); err != nil {
			log.Println("Loading moves: ", err)
		}
	}

	// setoption name Threads value 8
	cmdThreads := uci.CmdSetOption{
		Name:  "Threads",
		Value: "4",
	}

	cmdSkill := uci.CmdSetOption{
		Name:  "Skill Level",
		Value: "20",
	}

	depth := 21
	if len(move_list) > 60 {
		depth = 16
	}
	cmdPos := uci.CmdPosition{Position: game.Position()}
	cmdGo := uci.CmdGo{
		Depth:    depth,
		MoveTime: 1000 * time.Millisecond,
	}

	if err := eng.Run(cmdThreads, cmdSkill, cmdPos, cmdGo); err != nil {
		return nil, err
	}
	search_resultes := eng.SearchResults()
	move := search_resultes.BestMove

	pv_len := len(search_resultes.Info.PV)
	if pv_len > 14 {
		pv_len = 14
	}
	log.Println("Best Move:                 ", move)
	log.Println("Info: Depth / selective:   ", search_resultes.Info.Depth, " / ", search_resultes.Info.Seldepth)
	log.Println("Info: Score / Mate in:     ", search_resultes.Info.Score.CP, " / ", search_resultes.Info.Score.Mate)
	log.Println("Info: PV:                  ", search_resultes.Info.PV[:pv_len])
	log.Println("Info: NPS / Nodes:         ", search_resultes.Info.NPS, " / ", search_resultes.Info.Nodes)
	log.Println("Info: Time:                ", search_resultes.Info.Time)
	log.Println("---------------------------------------------------------------------------------------------------------")

	return move, nil
}
