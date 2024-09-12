package main

import (
	"log"
	"time"

	"github.com/tebeka/selenium"
)

type ChessBot struct {
	co chessOnline
}

func NewChessBot(co chessOnline) *ChessBot {
	return &ChessBot{co: co}
}

func (cb ChessBot) Run(driver selenium.WebDriver) {
	cb.co.PlayWithHuman()
	// lc.PlayWithComputer(driver)

	for {
		cb.co.IsPlayWithWhite()
		moveList := cb.co.GetMoveList()
		playMove, err := cb.co.PlayMoveWithMouse()
		if err != nil {
			log.Fatal(err)
		}

		for {
			cb.co.NewGame()
			moveList()

			if cb.co.IsMyTurn(cb.co.PlayWithWhite) && len(lc.MoveList) > 8 {
				err := lc.GetEngineBestMove()
				if err != nil {
					log.Println("Can't get best move from engine: ", err)
				} else {
					err := lc.GetTimeLeftSeconds(driver)
					if err != nil {
						log.Println("Can't get time left")
					}
					playMove(lc.BestMove.String(), len(lc.MoveList), lc.TimeLeftSeconds)
				}
			}
			lc.GetGameState(driver)
			if lc.GameState != "ongoing" {
				log.Println("Game State: ", lc.GameState)
				time.Sleep(3 * time.Second)
				lc.NewOpponent(driver)
				break
			}

		}
	}
}
