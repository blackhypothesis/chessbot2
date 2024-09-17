package main

import (
	"log"
	"time"
)

type ChessBot struct {
	co chessOnline
}

func NewChessBot(co chessOnline) *ChessBot {
	return &ChessBot{co: co}
}

func (cb ChessBot) Run() {
	err := cb.co.ConnectToSite()
	if err != nil {
		log.Fatal(err)
	}
	defer cb.co.ServiceStop()

	err = cb.co.PlayWithHuman()
	if err != nil {
		log.Fatal(err)
	}

	for {
		cb.co.IsPlayWithWhite()
		log.Printf("IsPlayWithWhite: %v", cb.co.GetPlayWithWhite())

		// get closure functions
		updateMoveList := cb.co.UpdateMoveList()
		playMove, err := cb.co.PlayMoveWithMouse()

		if err != nil {
			log.Fatal(err)
		}

		for {
			cb.co.NewGame()
			updateMoveList()

			if cb.co.IsMyTurn(cb.co.GetPlayWithWhite()) {
				err := cb.co.CalculateEngineBestMove()
				if err != nil {
					log.Println("Can't get best move from engine: ", err)
				} else {
					err := cb.co.CalculateTimeLeftSeconds()
					if err != nil {
						log.Println("Can't get time left")
					}
					if len(cb.co.GetMoveList()) > -1 {
						playMove(cb.co.GetBestMove(), len(cb.co.GetMoveList()), cb.co.GetTimeLeftSeconds())
						cb.co.PrintSearchResults()
					}
				}
			}
			cb.co.GetGameState()
			if cb.co.GetGameState() != "ongoing" {
				log.Println("Game State: ", cb.co.GetGameState())
				updateMoveList()
				cb.co.SaveGame()
				time.Sleep(3 * time.Second)
				cb.co.NewOpponent()
				break
			}
		}
	}
}
