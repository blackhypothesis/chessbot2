package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-vgo/robotgo"
	"github.com/tebeka/selenium"
)

func getMoveList(driver selenium.WebDriver) ([]string, error) {
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

func isWhiteOrientation(driver selenium.WebDriver) (bool, error) {
	board_coords, err := driver.FindElement(selenium.ByTagName, "coords")
	if err != nil {
		return false, err
	}
	board_coords_class, err := board_coords.GetAttribute("class")
	if err != nil {
		return false, err
	}
	if board_coords_class == "ranks black" {
		return false, nil
	} else {
		return true, nil
	}
}

func isMyTurn(move_list []string, is_white_orientation bool, driver selenium.WebDriver) bool {

	// we play with white
	if len(move_list)%2 == 0 && is_white_orientation {
		return true
	}
	// we play with black
	if len(move_list)%2 == 1 && !is_white_orientation {
		return true
	}

	return false

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

func playMoveWithMouse(driver selenium.WebDriver, is_white_orientation bool) (func(move string), error) {
	cg_board, err := driver.FindElement(selenium.ByTagName, "cg-board")

	robotgo.MouseSleep = 200

	// get board size
	board_size := new(selenium.Size)
	if err != nil {
		return nil, err
	}

	board_size, err = cg_board.Size()
	if err != nil {
		return nil, err
	}

	// get board location
	board_location := new(selenium.Point)
	if err != nil {
		return nil, err
	}

	board_location, err = cg_board.Location()
	if err != nil {
		return nil, err
	}

	field_size := new(selenium.Size)
	field_size.Width = board_size.Width / 8
	field_size.Height = board_size.Height / 8

	return func(move string) {

		mouse_click_start := new(selenium.Point)
		mouse_click_end := new(selenium.Point)
		fmt.Println(board_location, board_size, field_size)

		m := strings.Split(move, "")

		mouse_click_start.X = getCoordinate(m[0])
		mouse_click_start.Y, _ = strconv.Atoi(m[1])
		mouse_click_end.X = getCoordinate(m[2])
		mouse_click_end.Y, _ = strconv.Atoi(m[3])

		mouse_click_start.Y--
		mouse_click_end.Y--

		if !is_white_orientation {
			mouse_click_start.X = 7 - mouse_click_start.X
			mouse_click_start.Y = 7 - mouse_click_start.Y
			mouse_click_end.X = 7 - mouse_click_end.X
			mouse_click_end.Y = 7 - mouse_click_end.Y
		}

		y_offset := 170

		robotgo.Move(board_location.X+mouse_click_start.X*field_size.Width+field_size.Width/2, y_offset+board_location.Y+(7-mouse_click_start.Y)*field_size.Height+field_size.Height/2)
		robotgo.Click("left")
		robotgo.Move(board_location.X+mouse_click_end.X*field_size.Width+field_size.Width/2, y_offset+board_location.Y+(7-mouse_click_end.Y)*field_size.Height+field_size.Height/2)
		robotgo.Click("left")

	}, nil
}

func getCoordinate(x string) int {
	switch x {
	case "a":
		return 0
	case "b":
		return 1
	case "c":
		return 2
	case "d":
		return 3
	case "e":
		return 4
	case "f":
		return 5
	case "g":
		return 6
	case "h":
		return 7
	}
	return 0
}
