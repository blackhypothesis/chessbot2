package main

type chessOnline interface {
	// ConnectToSite() error
	SignIn() error
	PlayWithHuman() error
	PlayWithComputer() error
	NewGame()
	IsPlayWithWhite()
	GetMoveList() func() []string
	IsMyTurn(bool) bool
	GetEngineBestMove() error
	GetTimeLeftSeconds() error
	PlayMoveWithMouse() (func(string, int, [2]int), error)
	GetGameState()
	NewOpponent() error
}
