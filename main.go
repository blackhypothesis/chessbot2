package main

import (
	"log"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func main() {
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
	playWithComputer(driver)
	time.Sleep(1 * time.Second)

	// robotgo.MouseSleep = 200
	// robotgo.Move(1978, 1569)
	// robotgo.Click("left")
	// robotgo.Move(1978, 1138)
	// robotgo.Click("left")

	for {
		board := Board{}
		for {
			board, err = GetBoard(driver)
			if err != nil {
				log.Println("Board Error: ", err)
				log.Println("Will retry to get board information, ...")
				time.Sleep(10 * time.Second)
			} else {
				break
			}
		}

		log.Println("Orientation: ", board.orientation)
		log.Println("Location:    ", board.location)
		log.Println("Size:        ", board.size)
		log.Println("Field size:  ", board.field_size)

		log.Println("Active:      ", board.active_color)
		log.Println("FEN:         ", board.fen)
		log.Println("Moves:       ", board.move_list)

		x, y := robotgo.Location()
		log.Println("Mouse: ", x, y)

		windowTitle := robotgo.GetTitle()
		windowPID := robotgo.GetPid()
		x, y, w, h := robotgo.GetBounds(windowPID)
		log.Println("Window Title: ", windowTitle)
		log.Println("Window PID  : ", windowPID)
		log.Println("x: ", x, " y: ", y, " h: ", h, " w: ", w)

		game_state := getGameState(driver)
		log.Println("Game state: ", game_state)
		is_my_turn, err := IsMyTurn(&board, driver)
		if err != nil {
			log.Println("Error: ", err)
		}
		log.Println("###### is my turn: ", is_my_turn)
		log.Println("------------------------------------------------------------------------------------------")

		time.Sleep(100 * time.Millisecond)

	}
}
