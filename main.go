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
	time.Sleep(10 * time.Second)

	board, err := GetBoard(driver)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	fmt.Println(board)

	move_list, err := GetMoveList(driver)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	fmt.Println(move_list)

	time.Sleep(2 * time.Second)

}
