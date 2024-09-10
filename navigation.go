package main

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/tebeka/selenium"
)

func signIn(username string, password string, driver selenium.WebDriver) error {
	sign_in, err := driver.FindElement(selenium.ByClassName, "signin")
	if err != nil {
		return err
	}
	sign_in.Click()

	time.Sleep(1 * time.Second)

	form_username, err := driver.FindElement(selenium.ByID, "form3-username")
	if err != nil {
		return err
	}
	form_password, err := driver.FindElement(selenium.ByID, "form3-password")
	if err != nil {
		return err
	}
	button_submit, err := driver.FindElement(selenium.ByClassName, "submit")
	if err != nil {
		return err
	}
	log.Println("Login with: ", username, password)
	form_username.Clear()
	form_username.SendKeys(username)
	time.Sleep(1 * time.Second)

	form_password.Clear()
	form_password.SendKeys(password)
	time.Sleep(1 * time.Second)
	button_submit.Click()

	return nil
}

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

type envVAR struct {
	Login    string
	Password string
}

func getENV() (envVAR, error) {
	err := godotenv.Load()
	if err != nil {
		return envVAR{}, err
	}
	login := os.Getenv("LOGIN")
	if login == "" {
		return envVAR{}, errors.New("LOGIN is not found in the ENV")
	}
	password := os.Getenv("PASSWORD")
	if password == "" {
		return envVAR{}, errors.New("PASSWORD is not found in the ENV")
	}

	return envVAR{Login: login, Password: password}, nil
}
