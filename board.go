package main

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

type Board struct {
	orientation    string
	castling_right string
	fen            string
	active_color   string
	move_list      []string
	board          [8][8]Piece
	board_size_x   int
	board_size_y   int
}

type Piece struct {
	Piece string
	Color string
}

func GetBoard(driver selenium.WebDriver) (Board, error) {
	board := Board{}
	time.Sleep(1 * time.Second)
	cg_container, err := driver.FindElement(selenium.ByTagName, "cg-container")
	if err != nil {
		return Board{}, err
	}
	cg_board, err := driver.FindElement(selenium.ByTagName, "cg-board")
	if err != nil {
		return Board{}, err
	}

	// get board orientation
	board_coords, err := driver.FindElement(selenium.ByTagName, "coords")
	if err != nil {
		return Board{}, err
	}
	board_coords_class, err := board_coords.GetAttribute("class")
	if err != nil {
		return Board{}, err
	}
	if board_coords_class == "ranks black" {
		board.orientation = "black"
	} else {
		board.orientation = "white"
	}

	board_size_string, err := cg_container.GetAttribute("style")
	if err != nil {
		return Board{}, err
	}
	pattern := regexp.MustCompile(`width: (?P<x_size>\d+)px; height: (?P<y_size>\d+)px;`)
	board_size := pattern.FindStringSubmatch(board_size_string)
	board_x := board_size[1]
	board_y := board_size[2]
	board_size_x, err := strconv.Atoi(board_x)
	if err != nil {
		return Board{}, err
	}
	board_size_y, err := strconv.Atoi(board_y)
	if err != nil {
		return Board{}, err
	}
	board.board_size_x = board_size_x
	board.board_size_y = board_size_y
	field_size_x := board_size_x / 8
	field_size_y := board_size_y / 8

	pieces, err := cg_board.FindElements(selenium.ByTagName, "piece")
	if err != nil {
		return Board{}, err
	}

	for _, piece := range pieces {
		var x, y int
		var pt, color string

		piece_type_string, err := piece.GetAttribute("class")
		if err != nil {
			return Board{}, err
		}

		pattern := regexp.MustCompile(`(?P<color>\w+) (?P<type>\w+)`)
		piece_type := pattern.FindStringSubmatch(piece_type_string)
		color = piece_type[1]
		pt = piece_type[2]

		piece_coordinates_string, err := piece.GetAttribute("style")
		if err != nil {
			return Board{}, err
		}

		pattern = regexp.MustCompile(`transform: translate\((?P<x_c>\d+)px, (?P<y_c>\d+)px\);`)
		piece_coordinates := pattern.FindStringSubmatch(piece_coordinates_string)
		if len(piece_coordinates) == 0 {
			pattern = regexp.MustCompile(`transform: translate\((?P<x_c>\d+)px\);`)
			piece_coordinates = pattern.FindStringSubmatch(piece_coordinates_string)
			if len(piece_coordinates) == 0 {
				x = 0
				y = 0
			} else {
				x, err = strconv.Atoi(piece_coordinates[1])
				if err != nil {
					return Board{}, err
				}
				y = 0
			}
		} else {
			x, err = strconv.Atoi(piece_coordinates[1])
			if err != nil {
				return Board{}, err
			}
			y, err = strconv.Atoi(piece_coordinates[2])
			if err != nil {
				return Board{}, err
			}
			x = x / field_size_x
			y = 7 - y/field_size_y

			if board.orientation == "black" {
				x = 7 - x
				y = 7 - y
			}
			board.board[x][y] = Piece{Color: color, Piece: pt}
		}
	}

	move_list, err := GetMoveList(driver)
	if err != nil {
		return Board{}, err
	}

	active_color := "w"
	if len(move_list)%2 != 0 {
		active_color = "b"
	}
	board.active_color = active_color
	board.move_list = move_list

	fen, castling_right := GetFEN(board)
	board.fen = fen
	board.castling_right = castling_right

	return board, nil
}

func GetFEN(board Board) (string, string) {
	fen := ""
	space := 0
	piece_letter := ""
	for y := 7; y >= 0; y-- {
		for x := 0; x <= 7; x++ {
			switch board.board[x][y].Piece {
			case "pawn":
				piece_letter = "p"
			case "knight":
				piece_letter = "n"
			case "bishop":
				piece_letter = "b"
			case "rook":
				piece_letter = "r"
			case "queen":
				piece_letter = "q"
			case "king":
				piece_letter = "k"
			default:
				piece_letter = ""
			}

			if board.board[x][y].Color == "white" {
				piece_letter = strings.ToUpper(piece_letter)
			}
			if piece_letter == "" {
				space++
			} else {
				if space > 0 {
					fen = fen + strconv.Itoa(space)
					space = 0
				}
				fen = fen + piece_letter
			}
		}
		if space > 0 {
			fen = fen + strconv.Itoa(space)
			space = 0
		}
		if y > 0 {
			fen = fen + "/"
		}
	}

	// calculate castling rights
	castling_right := ""
	if board.board[4][0].Piece == "king" && board.board[4][0].Color == "white" {
		if board.board[7][0].Piece == "rook" && board.board[7][0].Color == "white" {
			castling_right = castling_right + "K"
		}
		if board.board[0][0].Piece == "rook" && board.board[0][0].Color == "white" {
			castling_right = castling_right + "Q"
		}
	}
	if board.board[4][7].Piece == "king" && board.board[4][7].Color == "black" {
		if board.board[7][7].Piece == "rook" && board.board[7][7].Color == "black" {
			castling_right = castling_right + "k"
		}
		if board.board[0][7].Piece == "rook" && board.board[0][7].Color == "black" {
			castling_right = castling_right + "q"
		}
	}
	if len(castling_right) == 0 {
		castling_right = "-"
	}
	fen = fen + " " + board.active_color + " " + castling_right + " - 0 1"

	return fen, castling_right
}
