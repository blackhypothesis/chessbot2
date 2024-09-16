package main

type chessOnline interface {
	ConnectToSite() error
	ServiceStop()
	SignIn() error
	PlayWithHuman() error
	PlayWithComputer() error
	NewGame()
	IsPlayWithWhite()
	UpdateMoveList() func()
	IsMyTurn(bool) bool
	CalculateEngineBestMove() error
	PrintSearchResults()
	CalculateTimeLeftSeconds() error
	PlayMoveWithMouse() (func(string, int, [2]int), error)
	GetGameState() string
	NewOpponent() error
	GetPlayWithWhite() bool
	GetMoveList() []string
	GetBestMove() string
	GetTimeLeftSeconds() [2]int
}
