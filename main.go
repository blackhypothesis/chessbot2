package main

import (
	"log"

	"github.com/blackhypothesis/chessbot2/chesscom"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	/*
		lc, err := lichess.New("2+1")
		if err != nil {
			log.Fatal(err)
		}
	*/
	cc, err := chesscom.New("1+1")
	if err != nil {
		log.Fatal(err)
	}

	chess_bot := NewChessBot(cc)
	chess_bot.TestRun()
}
