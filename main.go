package main

import (
	"log"

	"github.com/blackhypothesis/chessbot2/lichess"
)

func main() {

	lc, err := lichess.New()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(lc)

	chess_bot := NewChessBot(lc)

	chess_bot.Run()
}
