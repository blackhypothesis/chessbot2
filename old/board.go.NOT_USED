package main

import (
	"log"
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
	size           selenium.Size
	location       selenium.Point
	field_size     selenium.Size
	cg_board       selenium.WebElement
}

type Piece struct {
	Piece string
	Color string
}

func GetBoard(driver selenium.WebDriver) (Board, error) {
	board := Board{}
	time.Sleep(1 * time.Second)

	log.Println("get chess board")
	cg_board, err := driver.FindElement(selenium.ByTagName, "cg-board")
	if err != nil {
		return Board{}, err
	}
	board.cg_board = cg_board

	// get board orientation
	log.Println("get board orientation")
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

	// get board size
	log.Println("get board size")
	size, err := cg_board.Size()
	if err != nil {
		return Board{}, err
	}
	board.size.Height = size.Height
	board.size.Width = size.Width

	board.field_size.Width = size.Width / 8
	board.field_size.Height = size.Height / 8

	// get board location
	log.Println("get board location")
	location, err := cg_board.Location()
	if err != nil {
		return Board{}, err
	}
	board.location.X = location.X
	board.location.Y = location.Y

	// get pieces
	log.Println("get pieces")
	pieces, err := cg_board.FindElements(selenium.ByTagName, "piece")
	if err != nil {
		return Board{}, err
	}

	for _, piece := range pieces {
		var x, y int
		var pt, color string

		// the class contains the color and the type of piece
		piece_type_string, err := piece.GetAttribute("class")
		if err != nil {
			return Board{}, err
		}

		pattern := regexp.MustCompile(`(?P<color>\w+) (?P<type>\w+)`)
		piece_type := pattern.FindStringSubmatch(piece_type_string)
		color = piece_type[1]
		pt = piece_type[2]

		// the style attributes contains the coordinates of the piece
		piece_coordinates_string, err := piece.GetAttribute("style")
		if err != nil {
			return Board{}, err
		}

		// extract color and piece type from the class attribute
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
			x = x / board.field_size.Height
			y = 7 - y/board.field_size.Width

			if board.orientation == "black" {
				x = 7 - x
				y = 7 - y
			}
			board.board[x][y] = Piece{Color: color, Piece: pt}
		}
	}

	// get move list
	log.Println("get move list")

	moveList := getMoveList(driver)
	if err != nil {
		log.Fatal(err)
	}
	move_list := moveList()

	// get active color
	active_color := "w"
	if len(move_list)%2 != 0 {
		active_color = "b"
	}
	board.active_color = active_color
	board.move_list = move_list

	// get fen and castling rights
	log.Println("calculate FEN and castling rights")
	fen, castling_right := GetFEN(board)
	board.fen = fen
	board.castling_right = castling_right
	log.Println("all board informatoin collected")
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
