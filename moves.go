package main

import (
	"github.com/tebeka/selenium"
)

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

func IsMyTurn(board *Board, driver selenium.WebDriver) (bool, error) {

	// we play with white
	if len(board.move_list)%2 == 0 && board.orientation == "white" {
		return true, nil
	}
	// we play with black
	if len(board.move_list)%2 == 1 && board.orientation == "black" {
		return true, nil
	}

	return false, nil

	// wait for the element "your turn" for about 60 seconds
	// this does only work, when there is no time control, which is usualy not the case

	// err := driver.WaitWithTimeout(func(driver selenium.WebDriver) (bool, error) {
	// 	yourTurn, _ := driver.FindElement(selenium.ByXPATH, `//*[@id="main-wrap"]/main/div[1]/div[8]/div`)
	//
	// 	if yourTurn != nil {
	// 		return yourTurn.IsDisplayed()
	// 	}
	// 	return false, nil
	// }, 60*time.Second)
	//
	// if err != nil {
	// 	return false, err
	// }
	//
	// yourTurn, err := driver.FindElement(selenium.ByXPATH, `//*[@id="main-wrap"]/main/div[1]/div[8]/div`)
	// if err != nil {
	// 	return false, err
	// }
	//
	// yt, err := yourTurn.Text()
	// if err != nil {
	// 	return false, nil
	// }
	// fmt.Println("yourturn: ", yt)
	// return true, nil

}

func getGameState(driver selenium.WebDriver) string {
	game_state, err := driver.FindElement(selenium.ByClassName, "result")
	if err != nil {
		return "ongoing"
	}
	state, err := game_state.Text()
	if err != nil {
		return "unknown"
	}
	return state
}
