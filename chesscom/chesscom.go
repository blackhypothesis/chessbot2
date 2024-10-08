package chesscom

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/joho/godotenv"
	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type envVAR struct {
	Login    string
	Password string
}

type Chesscom struct {
	Url             string
	UserName        string
	Password        string
	TimeControl     string
	MoveList        []string
	PlayWithWhite   bool
	Game            *chess.Game
	BestMove        *chess.Move
	TimeLeftSeconds [2]int
	GameState       string
	Service         *selenium.Service
	Driver          selenium.WebDriver
}

func New() (*Chesscom, error) {
	env, err := getENV()
	if err != nil {
		log.Fatal("Error: ", err)
	}
	return &Chesscom{
		Url:         "https://chess.com",
		UserName:    env.Login,
		Password:    env.Password,
		TimeControl: "1+0",
	}, nil
}

/*
 * START implementation of interface chessOnline
 */
func (cc *Chesscom) ConnectToSite() error {
	service, err := selenium.NewChromeDriverService("./chromedriver-linux64/chromedriver", 4444)
	if err != nil {
		return err
	}

	//configure browser options
	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		"--headless-new", // comment out this line for testing
	}})

	// create a new remote client with the specified options
	driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		return err
	}

	// maximize the current window to avoid responsive rendering
	err = driver.MaximizeWindow("")
	if err != nil {
		return err
	}

	log.Println("Connecting to: ", cc.Url)
	err = driver.Get(cc.Url)
	if err != nil {
		return err
	}

	cc.Service = service
	cc.Driver = driver

	return nil
}

func (cc *Chesscom) ServiceStop() {
	cc.Service.Stop()
}

func (cc *Chesscom) SignIn() error {
	return nil
}

func (cc *Chesscom) PlayWithHuman() error {
	time.Sleep(1 * time.Second)
	play_url := fmt.Sprintf("%s/play/online", cc.Url)
	log.Println("Connecting to: ", play_url)
	err := cc.Driver.Get(play_url)
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	time_selector_button, err := cc.Driver.FindElement(selenium.ByClassName, "selector-button-button")
	if err != nil {
		return err
	}
	time_selector_button.Click()
	time.Sleep(500 * time.Millisecond)

	time_selectors, err := cc.Driver.FindElements(selenium.ByClassName, "time-selector-button-button")
	if err != nil {
		return err
	}
	switch cc.TimeControl {
	case "1+0":
		time_selectors[0].Click()
	case "1+1":
		time_selectors[1].Click()
	case "2+1":
		time_selectors[2].Click()
	case "3+0":
		time_selectors[3].Click()
	case "3+2":
		time_selectors[4].Click()
	case "5+0":
		time_selectors[5].Click()
	default:
		return errors.New("timecontrol does not exist")
	}
	time.Sleep(500 * time.Millisecond)
	button_play, err := cc.Driver.FindElement(selenium.ByClassName, "cc-button-xx-large")
	if err != nil {
		return err
	}
	button_play.Click()
	time.Sleep(500 * time.Millisecond)

	button_guest, err := cc.Driver.FindElement(selenium.ByID, "guest-button")
	if err != nil {
		return err
	}
	button_guest.Click()

	return nil
}

func (cc *Chesscom) PlayWithComputer() error {
	return nil
}

func (cc *Chesscom) NewGame() {
	cc.Game = chess.NewGame()
}

func (cc *Chesscom) IsPlayWithWhite() {
	coordinate := ""
	for {
		log.Println("trying to get board orientation, ...")
		coordinates_light, err := cc.Driver.FindElements(selenium.ByClassName, "coordinate-light")
		if err != nil {
			continue
		}
		coordinate, err = coordinates_light[0].Text()
		if err != nil {
			continue
		}
		break
	}
	fmt.Printf("board coordinate upper left: %s\n", coordinate)
	if coordinate == "1" {
		cc.PlayWithWhite = false
	} else if coordinate == "8" {
		cc.PlayWithWhite = true
	}
}

func (cc *Chesscom) UpdateMoveList() func() {
	defer TimeTrack(time.Now())
	move_list := []string{}
	last_move_list_len := 0

	return func() {
		move_list_container, err := cc.Driver.FindElements(selenium.ByTagName, "kwdb")
		if err != nil {
			log.Println("Can't get move list container")
		}
		move_list_container_len := len(move_list_container)
		number_new_moves := move_list_container_len - last_move_list_len
		for move_index := number_new_moves; move_index > 0; move_index-- {
			move_element := move_list_container[len(move_list_container)-move_index]
			move, err := move_element.Text()
			if err != nil {
				log.Println("Can't decode move: move index: ", move_index)
			}
			move_list = append(move_list, move)
		}
		last_move_list_len = len(move_list)
		cc.MoveList = move_list
	}
}

func (cc *Chesscom) IsMyTurn(is_white_orientation bool) bool {
	// we play with white
	if len(cc.MoveList)%2 == 0 && is_white_orientation {
		return true
	}
	// we play with black
	if len(cc.MoveList)%2 == 1 && !is_white_orientation {
		return true
	}
	return false
}

func (cc *Chesscom) CalculateEngineBestMove() error {
	// defer TimeTrack(time.Now())

	cc.Game = chess.NewGame()
	eng, err := uci.New("stockfish")
	if err != nil {
		return err
	}
	// initialize uci with new game
	if err := eng.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
		return err
	}
	defer eng.Close()

	for _, move := range cc.MoveList {
		if err := cc.Game.MoveStr(move); err != nil {
			log.Println("Loading moves: ", err)
		}
	}

	// setoption name Threads value 8
	cmdThreads := uci.CmdSetOption{
		Name:  "Threads",
		Value: "4",
	}

	cmdSkill := uci.CmdSetOption{
		Name:  "Skill Level",
		Value: "20",
	}

	depth := 21
	if len(cc.MoveList) > 60 {
		depth = 16
	}
	cmdPos := uci.CmdPosition{Position: cc.Game.Position()}
	cmdGo := uci.CmdGo{
		Depth:    depth,
		MoveTime: 1000 * time.Millisecond,
	}

	if err := eng.Run(cmdThreads, cmdSkill, cmdPos, cmdGo); err != nil {
		return err
	}
	search_resultes := eng.SearchResults()
	move := search_resultes.BestMove

	pv_len := len(search_resultes.Info.PV)
	if pv_len > 14 {
		pv_len = 14
	}
	log.Println("Best Move:                 ", move)
	log.Println("Info: Depth / selective:   ", search_resultes.Info.Depth, " / ", search_resultes.Info.Seldepth)
	log.Println("Info: Score / Mate in:     ", search_resultes.Info.Score.CP, " / ", search_resultes.Info.Score.Mate)
	log.Println("Info: PV:                  ", search_resultes.Info.PV[:pv_len])
	log.Println("Info: NPS / Nodes:         ", search_resultes.Info.NPS, " / ", search_resultes.Info.Nodes)
	log.Println("Info: Time:                ", search_resultes.Info.Time)
	log.Println("---------------------------------------------------------------------------------------------------------")

	cc.BestMove = move
	return nil
}

func (cc *Chesscom) CalculateTimeLeftSeconds() error {
	time_left, err := cc.Driver.FindElements(selenium.ByClassName, "time")
	if err != nil {
		return err
	}
	// sometimes it crashe, because of:
	//   panic: runtime error: index out of range [0] with length 0
	// therefore check if len is 2
	if len(time_left) == 2 {
		time_opponent, _ := time_left[0].Text()
		time_self, _ := time_left[1].Text()
		time_opponent_minutes_seconds := strings.Split(strings.Replace(time_opponent, "\n", "", -1), ":")
		time_self_minutes_seconds := strings.Split(strings.Replace(time_self, "\n", "", -1), ":")

		time_opponent_minutes, _ := strconv.Atoi(time_opponent_minutes_seconds[0])
		time_opponent_seconds, _ := strconv.Atoi(time_opponent_minutes_seconds[1])
		time_self_minutes, _ := strconv.Atoi(time_self_minutes_seconds[0])
		time_self_seconds, _ := strconv.Atoi(time_self_minutes_seconds[1])
		time_opponent_secs := 60*time_opponent_minutes + time_opponent_seconds
		time_self_secs := 60*time_self_minutes + time_self_seconds

		cc.TimeLeftSeconds = [2]int{time_self_secs, time_opponent_secs}
	}
	return nil
}

func (cc *Chesscom) PlayMoveWithMouse() (func(move string, len_move_list int, time_left_seconds [2]int), error) {
	defer TimeTrack(time.Now())
	cg_board, err := cc.Driver.FindElement(selenium.ByTagName, "cg-board")
	if err != nil {
		return nil, err
	}
	robotgo.MouseSleep = 50

	// get board size
	board_size := new(selenium.Size)
	board_size, err = cg_board.Size()
	if err != nil {
		return nil, err
	}
	// get board location
	board_location := new(selenium.Point)
	board_location, err = cg_board.Location()
	if err != nil {
		return nil, err
	}

	field_size := new(selenium.Size)
	field_size.Width = board_size.Width / 8
	field_size.Height = board_size.Height / 8

	min_wait_seconds := 0.5
	max_wait_seconds := 10.0

	return func(move string, len_move_list int, time_left_seconds [2]int) {
		piece_start := new(selenium.Point)
		piece_end := new(selenium.Point)
		m := strings.Split(move, "")
		piece_start.X = getCoordinate(m[0])
		piece_start.Y, _ = strconv.Atoi(m[1])
		piece_end.X = getCoordinate(m[2])
		piece_end.Y, _ = strconv.Atoi(m[3])
		piece_start.Y--
		piece_end.Y--
		if !cc.PlayWithWhite {
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

		// artificially wait for some time to get higher standarddeviation of move time usage
		max_wait_secs := 0.0
		// do not spend too much time in the first 6 moves
		if len_move_list < 12 {
			max_wait_secs = 2.0
		} else {
			// play faster when the time left to play is lower
			// keep 15 seconds as reserve
			max_wait_secs = float64(time_left_seconds[0]-15) / 10
			if max_wait_secs < 0.8 {
				max_wait_secs = 0.8
			}
		}
		if max_wait_secs > max_wait_seconds {
			max_wait_secs = max_wait_seconds
		}
		wait_seconds := min_wait_seconds + rand.Float64()*(max_wait_secs-min_wait_seconds)
		log.Printf("Waiting for %f seconds ...\n", wait_seconds)
		// time.Sleep(time.Duration(wait_seconds) * time.Second)
		log.Printf("Play move: %s\n", move)
		robotgo.Move(location_start.X, location_start.Y)
		robotgo.Click("left")
		robotgo.Move(location_end.X, location_end.Y)
		robotgo.Click("left")

		// piece promotion
		// cacculate the field to click, to promote the pawn to the desired piece
		if len(m) == 5 {
			promotion_click_square := new(selenium.Point)
			promotion_click_square.X = piece_end.X

			switch m[4] {
			case "q":
				promotion_click_square.Y = 7
			case "n":
				promotion_click_square.Y = 6
			case "r":
				promotion_click_square.Y = 5
			case "b":
				promotion_click_square.Y = 4
			}
			location := selenium.Point{
				X: board_location.X + promotion_click_square.X*field_size.Width + field_size.Width/2,
				Y: y_offset + board_location.Y + (7-promotion_click_square.Y)*field_size.Height + field_size.Height/2,
			}

			robotgo.Move(location.X, location.Y)
			robotgo.Click("left")
		}
	}, nil
}

func (cc *Chesscom) GetGameState() string {
	game_state, err := cc.Driver.FindElement(selenium.ByClassName, "result")
	if err != nil {
		cc.GameState = "ongoing"
	}
	state, err := game_state.Text()
	if err != nil {
		cc.GameState = "unknown"
	}
	cc.GameState = state
	return cc.GameState
}

func (cc *Chesscom) NewOpponent() error {
	new_opponent, err := cc.Driver.FindElement(selenium.ByXPATH, `//*[@id="main-wrap"]/main/div[1]/div[5]/div/a[1]`) // New opponent
	if err != nil {
		return err
	}
	new_opponent.Click()
	return nil
}

// getter functions
func (cc *Chesscom) GetPlayWithWhite() bool {
	return cc.PlayWithWhite
}
func (cc *Chesscom) GetMoveList() []string {
	return cc.MoveList
}
func (cc *Chesscom) GetBestMove() string {
	return cc.BestMove.String()
}
func (cc *Chesscom) GetTimeLeftSeconds() [2]int {
	return cc.TimeLeftSeconds
}

/*
 * END implementation of interface chessOnline
 */

/*
 * helper functions
 */
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

func TimeTrack(start time.Time) {
	elapsed := time.Since(start)

	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(1)

	// Retrieve a function object this functions parent.
	funcObj := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path).
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

	log.Printf("%s took %s", name, elapsed)
}
