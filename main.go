package main

import (
	"log"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	eng, err := uci.New("stockfish")
	if err != nil {
		panic(err)
	}
	// initialize uci with new game
	if err := eng.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
		panic(err)
	}
	defer eng.Close()

	service, err := selenium.NewChromeDriverService("./chromedriver-linux64/chromedriver", 4444)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	defer service.Stop()

	//configure browser options
	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		"--headless-new", // comment out this line for testing
	}})

	// create a new remote client with the specified options
	driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		log.Fatal("Error: ", err)
	}

	// maximize the current window to avoid responsive rendering
	err = driver.MaximizeWindow("")
	if err != nil {
		log.Fatal("Error: ", err)
	}

	err = driver.Get("https:lichess.org/")
	if err != nil {
		log.Fatal("Error: ", err)
	}

	time.Sleep(2 * time.Second)
	// playWithComputer(driver)
	time.Sleep(1 * time.Second)

	// robotgo.MouseSleep = 200
	// robotgo.Move(1978, 1569)
	// robotgo.Click("left")
	// robotgo.Move(1978, 1138)
	// robotgo.Click("left")

	is_white_orientation := true

	for {
		is_white_orientation, err = isWhiteOrientation(driver)
		if err != nil {
			log.Println("IsWhiteOrientation Error: ", err)
			log.Println("Will retry to get orientation, ...")
		} else {
			break
		}
	}

	playMove, err := playMoveWithMouse(driver, is_white_orientation)

	if err != nil {
		log.Fatal(err)
	}

	for {
		log.Println("Get move list ...")
		move_list, err := getMoveList(driver)
		if err != nil {
			log.Println(err)
		}

		log.Println("Movelist processed.")
		game := chess.NewGame()

		if isMyTurn(move_list, is_white_orientation, driver) {
			for _, move := range move_list {
				if err := game.MoveStr(move); err != nil {
					log.Println("Loading moves: ", err)
				}
			}
			cmdPos := uci.CmdPosition{Position: game.Position()}
			cmdGo := uci.CmdGo{MoveTime: time.Second / 1}
			if err := eng.Run(cmdPos, cmdGo); err != nil {
				panic(err)
			}
			move := eng.SearchResults().BestMove
			playMove(move.String())
			log.Println("ENGINE: best move: ", move)

			log.Println("CHESS: ", game.Position().Board())
		}
		x, y := robotgo.Location()
		log.Println("Mouse location: ", x, y)
	}
}
