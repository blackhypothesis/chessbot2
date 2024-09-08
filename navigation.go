package main

import (
	"time"

	"github.com/tebeka/selenium"
)

func playWithHuman(time_setting string, driver selenium.WebDriver) error {
	time_settings, err := driver.FindElements(selenium.ByClassName, "clock")
	if err != nil {
		return err
	}
	switch time_setting {
	case "1+0":
		time_settings[0].Click()
	case "2+1":
		time_settings[1].Click()
	case "3+0":
		time_settings[2].Click()
	case "3+2":
		time_settings[3].Click()
	}
	return nil
}

func playWithComputer(driver selenium.WebDriver) error {
	// Button [PLAY WITH COMPUTER]
	button, err := driver.FindElement(selenium.ByClassName, "config_ai")
	if err != nil {
		return err
	}
	button.Click()
	time.Sleep(500 * time.Millisecond)

	// Button Strength [8]
	level, err := driver.FindElement(selenium.ByXPATH, "/html/body/div/main/div[1]/dialog/div[2]/div/div/div[3]/div[1]/group/div[8]/label")
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

func giveMoreTime(driver selenium.WebDriver) error {
	more_time, err := driver.FindElement(selenium.ByClassName, "moretime")
	if err != nil {
		return err
	}
	more_time.Click()
	return nil
}

func newOpponent(driver selenium.WebDriver) error {
	new_opponent, err := driver.FindElement(selenium.ByXPATH, `//*[@id="main-wrap"]/main/div[1]/div[5]/div/a[1]`) // New opponent
	if err != nil {
		return err
	}
	new_opponent.Click()
	return nil
}
