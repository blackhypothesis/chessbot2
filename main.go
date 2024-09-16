package main

import (
	"log"

	"github.com/blackhypothesis/chessbot2/lichess"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	lc, err := lichess.New()
	if err != nil {
		log.Fatal(err)
	}

	chess_bot := NewChessBot(lc)
	chess_bot.Run()
}
