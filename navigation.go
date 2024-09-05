package main

import (
	"time"

	"github.com/tebeka/selenium"
)

func playWithComputer(driver selenium.WebDriver) error {
	// Button [PLAY WITH COMPUTER]
	button, err := driver.FindElement(selenium.ByClassName, "config_ai")
	if err != nil {
		return err
	}
	button.Click()
	time.Sleep(500 * time.Millisecond)

	// Button Strength [4]
	level, err := driver.FindElement(selenium.ByXPATH, "/html/body/div/main/div[1]/dialog/div[2]/div/div/div[3]/div[1]/group/div[4]/label")
	if err != nil {
		return err
	}
	level.Click()
	time.Sleep(500 * time.Millisecond)

	// Dropdown Time Control
	tc, _ := driver.FindElement(selenium.ByID, "sf_timeMode")
	if err != nil {
		return err
	}
	tcv, err := tc.FindElements(selenium.ByTagName, "option")
	if err != nil {
		return err
	}
	// Real Time (first option in dropdown)
	tcv[0].Click()
	time.Sleep(500 * time.Millisecond)

	// Button [white/black]
	bw, err := driver.FindElement(selenium.ByXPATH, "/html/body/div/main/div[1]/dialog/div[2]/div/div/div[4]/button[2]")
	if err != nil {
		return err
	}
	bw.Click()
	return nil
}
