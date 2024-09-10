package main

import (
	"fmt"
	"log"
	"regexp"
	"runtime"
	"time"

	"github.com/notnil/chess"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func main() {
	env, err := getENV()
	if err != nil {
		log.Fatal("Error: ", err)
	}
	log.Println("Login: ", env.Login, "  Password: ", env.Password)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

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
	// err = signIn(env.Login, env.Password, driver)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	time.Sleep(2 * time.Second)
	playWithHuman("1+0", driver)

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

			if isMyTurn(move_list, is_white_orientation) {
				move, err := getEngineBestMove(game, move_list)
				if err != nil {
					log.Println("Can't get best move from engine: ", err)
				} else {
					time_left_seconds, err := getTimeLeftSeconds(driver)
					if err != nil {
						log.Println("Can't get time left")
					}
					playMove(move.String(), len(move_list), time_left_seconds)
				}
				if err != nil {
					log.Println("Can't get time left")
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
