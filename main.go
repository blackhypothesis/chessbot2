package main

import (
	"fmt"
	"log"
	"regexp"
	"runtime"
	"time"

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
	//playWithComputer(driver)
	playWithHuman("1+0", driver)
	time.Sleep(1 * time.Second)

	for {
		time.Sleep(2 * time.Second)
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

		// get closure function to play moves
		playMove, err := playMoveWithMouse(driver, is_white_orientation)
		if err != nil {
			log.Fatal(err)
		}
		// get closure function to get move list
		moveList := getMoveList(driver)
		if err != nil {
			log.Fatal(err)
		}

		for {
			move_list := moveList()
			game := chess.NewGame()

			if isMyTurn(move_list, is_white_orientation) && len(move_list) > 20 {
				move, err := getEngineBestMove(game, eng, move_list)
				if err != nil {
					log.Println("Can't get best move from engine: ", err)
				} else {
					playMove(move.String())
				}
			}
			game_state := getGameState(driver)
			if game_state != "ongoing" {
				log.Println("Game State: ", game_state)
				time.Sleep(4 * time.Second)
				newOpponent(driver)
				break
			}
		}
	}
}

func TimeTrack(start time.Time) {
	elapsed := time.Since(start)

	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(1)

	// Retrieve a function object this functions parent.
	funcObj := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path).
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

	log.Println(fmt.Sprintf("%s took %s", name, elapsed))
}
