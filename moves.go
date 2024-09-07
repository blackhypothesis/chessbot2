package main

import (
	"strconv"
	"strings"

	"github.com/go-vgo/robotgo"
	"github.com/tebeka/selenium"
)

func getMoveList(driver selenium.WebDriver) func() []string {
	move_list := []string{}
	last_move_list_len := 0

	return func() []string {
		move_list_container, err := driver.FindElements(selenium.ByTagName, "kwdb")
		if err != nil {
			return move_list
		}
		move_list_container_len := len(move_list_container)
		number_new_moves := move_list_container_len - last_move_list_len

		if number_new_moves > 0 {
			for move_index := number_new_moves; move_index > 0; move_index-- {
				move_element := move_list_container[len(move_list_container)-move_index]
				move, err := move_element.Text()
				if err != nil {
					return move_list
				}
				move_list = append(move_list, move)
			}
			last_move_list_len = len(move_list)
		}
		return move_list
	}
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

func isMyTurn(move_list []string, is_white_orientation bool) bool {
	// we play with white
	if len(move_list)%2 == 0 && is_white_orientation {
		return true
	}
	// we play with black
	if len(move_list)%2 == 1 && !is_white_orientation {
		return true
	}
	return false
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

	robotgo.MouseSleep = 100

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
		piece_start := new(selenium.Point)
		piece_end := new(selenium.Point)
		m := strings.Split(move, "")
		piece_start.X = getCoordinate(m[0])
		piece_start.Y, _ = strconv.Atoi(m[1])
		piece_end.X = getCoordinate(m[2])
		piece_end.Y, _ = strconv.Atoi(m[3])
		piece_start.Y--
		piece_end.Y--
		if !is_white_orientation {
			piece_start.X = 7 - piece_start.X
			piece_start.Y = 7 - piece_start.Y
			piece_end.X = 7 - piece_end.X
			piece_end.Y = 7 - piece_end.Y
		}

		y_offset := 170
		location_start := selenium.Point{
			X: board_location.X + piece_start.X*field_size.Width + field_size.Width/2,
			Y: y_offset + board_location.Y + (7-piece_start.Y)*field_size.Height + field_size.Height/2,
		}
		location_end := selenium.Point{
			X: board_location.X + piece_end.X*field_size.Width + field_size.Width/2,
			Y: y_offset + board_location.Y + (7-piece_end.Y)*field_size.Height + field_size.Height/2,
		}

		robotgo.Move(location_start.X, location_start.Y)
		robotgo.Click("left")
		robotgo.Move(location_end.X, location_end.Y)
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
