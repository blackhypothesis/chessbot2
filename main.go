package main

import (
	"log"

	"github.com/blackhypothesis/chessbot2/chesscom"
	"github.com/blackhypothesis/chessbot2/lichess"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	lc, err := lichess.New("1+0")
	if err != nil {
		log.Fatal(err)
	}
	cc, err := chesscom.New("2+1")
	if err != nil {
		log.Fatal(err)
	}

	li_chess := true

	if li_chess {
		chess_bot := NewChessBot(lc)
		chess_bot.Run()
	} else {
		chess_bot := NewChessBot(cc)
		chess_bot.TestRun()
	}

}
