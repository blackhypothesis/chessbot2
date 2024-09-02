package main

import (
	"time"

	"github.com/tebeka/selenium"
)

func playWithComputer(driver selenium.WebDriver) error {
	button, err := driver.FindElement(selenium.ByXPATH, "/html/body/div/main/div[1]/div[1]/button[3]")
	if err != nil {
		return err
	}
	button.Click()
	time.Sleep(500 * time.Millisecond)

	level, err := driver.FindElement(selenium.ByXPATH, "/html/body/div/main/div[1]/dialog/div[2]/div/div/div[3]/div[1]/group/div[4]/label")
	if err != nil {
		return err
	}
	level.Click()
	time.Sleep(500 * time.Millisecond)

	bw, err := driver.FindElement(selenium.ByXPATH, "/html/body/div/main/div[1]/dialog/div[2]/div/div/div[4]/button[2]")
	if err != nil {
		return err
	}
	bw.Click()
	return nil
}
