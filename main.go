package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func main() {
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

	robotgo.MouseSleep = 200

	robotgo.Move(1978, 1569)
	robotgo.Click("left")
	robotgo.Move(1978, 1138)
	robotgo.Click("left")

	for {
		board, err := GetBoard(driver)
		if err != nil {
			log.Fatal("Error: ", err)
		}

		fmt.Println("Location:   ", board.location)
		fmt.Println("Size:       ", board.size)
		fmt.Println("Field size: ", board.field_size)

		fmt.Println("Active:     ", board.active_color)
		fmt.Println("FEN:        ", board.fen)
		fmt.Println("Moves:      ", board.move_list)
		fmt.Println("")

		x, y := robotgo.Location()
		fmt.Println("Mouse: ", x, y)

		is_my_turn, err := IsMyTurn(driver)
		if err != nil {
			log.Println("Error: ", err)
		}
		fmt.Println("is my turn: ", is_my_turn)

		time.Sleep(1 * time.Second)

	}
}
