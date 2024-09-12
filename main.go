package main

import (
	"log"

	"github.com/blackhypothesis/chessbot2/chesscom"
	"github.com/blackhypothesis/chessbot2/lichess"
)

func main() {
	lc, err := lichess.New()
	if err != nil {
		log.Fatal(err)
	}

	cc, err := chesscom.New()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(lc)

	chess_bot := NewChessBot(lc)
	chess_bot_2 := NewChessBot(cc)

	chess_bot.Run()
	chess_bot_2.Run()
}
