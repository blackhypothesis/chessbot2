package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	lc := new(Lichess)
	lc.Url = "https://lichess.org"
	env, err := getENV()
	if err != nil {
		log.Fatal("Error: ", err)
	}

	lc.Url = "https://lichess.org"
	lc.UserName = env.Login
	lc.Password = env.Password
	lc.TimeControl = "1+0"

	log.Println("Login: ", lc.UserName, "  Password: ", lc.Password)

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

	err = driver.Get(lc.Url)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	lc.PlayWithHuman(driver)
	// lc.PlayWithComputer(driver)

	for {
		lc.IsPlayWithWhite(driver)
		moveList := lc.GetMoveList(driver)
		if err != nil {
			log.Fatal(err)
		}
		playMove, err := lc.PlayMoveWithMouse(driver)
		if err != nil {
			log.Fatal(err)
		}

		for {
			lc.NewGame()
			moveList()

			if lc.IsMyTurn(lc.PlayWithWhite) && len(lc.MoveList) > 8 {
				err := lc.GetEngineBestMove()
				if err != nil {
					log.Println("Can't get best move from engine: ", err)
				} else {
					err := lc.GetTimeLeftSeconds(driver)
					if err != nil {
						log.Println("Can't get time left")
					}
					fmt.Println("movelist: ", lc.MoveList)
					fmt.Println("bestmove: ", lc.BestMove)
					playMove(lc.BestMove.String(), len(lc.MoveList), lc.TimeLeftSeconds)
				}
			}
			lc.GetGameState(driver)
			if lc.GameState != "ongoing" {
				log.Println("Game State: ", lc.GameState)
				time.Sleep(3 * time.Second)
				lc.NewOpponent(driver)
				break
			}

		}
	}
}
