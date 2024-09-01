package main

import "github.com/tebeka/selenium"

func GetMoveList(driver selenium.WebDriver) ([]string, error) {
	move_list := []string{}

	move_list_container, err := driver.FindElements(selenium.ByTagName, "kwdb")
	if err != nil {
		return move_list, err
	}

	for _, move := range move_list_container {
		move_str, err := move.Text()
		if err != nil {
			return move_list, err
		}
		move_list = append(move_list, move_str)
	}
	return move_list, nil
}
