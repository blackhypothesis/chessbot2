package main

import (
	"fmt"
	"log"
	"time"

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

	game_2_1, err := driver.FindElement(selenium.ByXPATH, "/html/body/div/main/div[2]/div[2]/div[2]/div[1]")
	if err != nil {
		log.Fatal("Error: ", err)
	}

	game_2_1.Click()

	board, err := GetBoard(driver)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	fmt.Println(board)

	time.Sleep(2 * time.Second)
	cg_board, err := driver.FindElement(selenium.ByTagName, "cg-board")
	if err != nil {
		log.Fatal("Error: ", err)
	}

	for x := 1500; x < 1600; x = x + 100 {
		for y := 1500; y < 1600; y = y + 100 {
			cg_board.MoveTo(x, y)
			cg_board.Click()
			fmt.Println(x, y, "clicked")
			time.Sleep(4 * time.Second)
		}
	}
	cg_board.Click()

	time.Sleep(20 * time.Second)

}
