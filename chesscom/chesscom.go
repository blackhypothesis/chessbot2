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
	NumberOfGames   int
	SearchResults   uci.SearchResults
	TimeLeftSeconds [2]int
	GameState       string
	Service         *selenium.Service
	Driver          selenium.WebDriver
}

func New(time_control string) (*Chesscom, error) {
	env, err := getENV()
	if err != nil {
		log.Fatal("Error: ", err)
	}
	return &Chesscom{
		Url:         "https://chess.com",
		UserName:    env.Login,
		Password:    env.Password,
		TimeControl: time_control,
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

// not jet implemented
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
	// time selector button, click on it reveales a drop down menu with time controls
	err = cc.Driver.WaitWithTimeout(func(driver selenium.WebDriver) (bool, error) {
		time_selector_button, _ := cc.Driver.FindElement(selenium.ByClassName, "selector-button-button")
		if time_selector_button != nil {
			return time_selector_button.IsDisplayed()
		}
		return false, nil
	}, 5*time.Second)
	if err != nil {
		return err
	}
	time_selector_button, err := cc.Driver.FindElement(selenium.ByClassName, "selector-button-button")
	if err != nil {
		return err
	}
	time_selector_button.Click()

	// buttons in dropdown menu with time controls
	i := 0
	switch cc.TimeControl {
	case "1+0":
		i = 0
	case "1+1":
		i = 1
	case "2+1":
		i = 2
	case "3+0":
		i = 3
	case "3+2":
		i = 4
	case "5+0":
		i = 5
	default:
		return errors.New("timecontrol does not exist")
	}

	err = cc.Driver.WaitWithTimeout(func(driver selenium.WebDriver) (bool, error) {
		time_selectors, _ := cc.Driver.FindElements(selenium.ByClassName, "time-selector-button-button")
		if time_selectors != nil {
			return time_selectors[i].IsDisplayed()
		}
		return false, nil
	}, 5*time.Second)
	if err != nil {
		return err
	}
	time_selectors, err := cc.Driver.FindElements(selenium.ByClassName, "time-selector-button-button")
	if err != nil {
		return err
	}
	time_selectors[i].Click()

	// large button [Play]
	err = cc.Driver.WaitWithTimeout(func(driver selenium.WebDriver) (bool, error) {
		button_play, _ := cc.Driver.FindElement(selenium.ByClassName, "cc-button-xx-large")
		if button_play != nil {
			return button_play.IsDisplayed()
		}
		return false, nil
	}, 5*time.Second)
	if err != nil {
		return err
	}
	button_play, err := cc.Driver.FindElement(selenium.ByClassName, "cc-button-xx-large")
	if err != nil {
		return err
	}
	button_play.Click()

	// if "Play as a Guest" appears, click on the webelement
	err = cc.Driver.WaitWithTimeout(func(driver selenium.WebDriver) (bool, error) {
		button_guest, _ := cc.Driver.FindElement(selenium.ByID, "guest-button")
		if button_guest != nil {
			return button_guest.IsDisplayed()
		}
		return false, nil
	}, 5*time.Second)
	if err != nil {
		fmt.Println("'Play as a Guest' text does not appear")
	}
	button_guest, err := cc.Driver.FindElement(selenium.ByID, "guest-button")
	if err != nil {
		return err
	}
	button_guest.Click()

	return nil
}

// not jet implemented
func (cc *Chesscom) PlayWithComputer() error {
	return nil
}

func (cc *Chesscom) NewGame() {
	cc.Game = chess.NewGame()
}

func (cc *Chesscom) IsPlayWithWhite() {
	for {
		// get clocks
		// since there was no webelement to hint if to play with white or black:
		// get the time from both players. if the seconds differ from 0, then the clock has started.
		// this means, that the player with this clock (bottom or top) has white.
		err := cc.Driver.WaitWithTimeout(func(driver selenium.WebDriver) (bool, error) {
			clocks, _ := cc.Driver.FindElements(selenium.ByClassName, "clock-time-monospace")
			if clocks != nil {
				return clocks[0].IsDisplayed()
			}
			return false, nil
		}, 1500*time.Millisecond)
		if err != nil {
			continue
		}
		clocks, err := cc.Driver.FindElements(selenium.ByClassName, "clock-time-monospace")
		if err != nil {
			continue
		}
		time_opponent, _ := clocks[0].Text()
		time_self, _ := clocks[1].Text()
		time_opponent_minutes_seconds := strings.Split(strings.Replace(time_opponent, "\n", "", -1), ":")
		time_self_minutes_seconds := strings.Split(strings.Replace(time_self, "\n", "", -1), ":")

		time_opponent_seconds, _ := strconv.Atoi(time_opponent_minutes_seconds[1])
		time_self_seconds, _ := strconv.Atoi(time_self_minutes_seconds[1])

		if time_opponent_seconds != 0 {
			cc.PlayWithWhite = false
			break
		}
		if time_self_seconds != 0 {
			cc.PlayWithWhite = true
			break
		}
	}
}

func (cc *Chesscom) UpdateMoveList() func() {
	defer TimeTrack(time.Now())

	return func() {
		move_list_container, err := cc.Driver.FindElements(selenium.ByClassName, "main-line-row")
		if err != nil {
			log.Println("Can't get move list container")
		}
		move_list := []string{}
		figurine := ""
		for _, move := range move_list_container {
			move_white_container, err := move.FindElement(selenium.ByClassName, "white-move")
			if err != nil {
				continue
			}
			move_white, err := move_white_container.FindElement(selenium.ByClassName, "node-highlight-content")
			if err != nil {
				continue
			}
			move_white_text, err := move_white.Text()
			if err != nil {
				continue
			}
			move_figurine_white, err := move_white.FindElement(selenium.ByClassName, "icon-font-chess")
			if err != nil {
				figurine = ""
			} else {
				figurine, err = move_figurine_white.GetAttribute("data-figurine")
				if err != nil {
					figurine = ""
				}
			}
			move_list = append(move_list, figurine+move_white_text)

			move_black_container, err := move.FindElement(selenium.ByClassName, "black-move")
			if err != nil {
				break
			}
			move_black, err := move_black_container.FindElement(selenium.ByClassName, "node-highlight-content")
			if err != nil {
				break
			}
			move_black_text, err := move_black.Text()
			if err != nil {
				break
			}
			move_figurine_black, err := move_black.FindElement(selenium.ByClassName, "icon-font-chess")
			if err != nil {
				figurine = ""
			} else {
				figurine, err = move_figurine_black.GetAttribute("data-figurine")
				if err != nil {
					figurine = ""
				}
			}
			move_list = append(move_list, figurine+move_black_text)

			log.Println("Movelist: ", move_list)
		}
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

	// setoption name Threads value 4
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
	cc.SearchResults = eng.SearchResults()

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

func (cc *Chesscom) PrintSearchResults() {
	pv_len := len(cc.SearchResults.Info.PV)
	if len(cc.SearchResults.Info.PV) > 14 {
		pv_len = 14
	}
	log.Println("Best Move:                 ", cc.SearchResults.BestMove.String())
	log.Println("Info: PV:                  ", cc.SearchResults.Info.PV[:pv_len])
	log.Println("Info: Depth / selective:   ", cc.SearchResults.Info.Depth, " / ", cc.SearchResults.Info.Seldepth)
	log.Println("Info: Score / Mate in:     ", cc.SearchResults.Info.Score.CP, " / ", cc.SearchResults.Info.Score.Mate)
	log.Println("Info: NPS / Nodes:         ", cc.SearchResults.Info.NPS, " / ", cc.SearchResults.Info.Nodes)
	log.Println("Info: Time:                ", cc.SearchResults.Info.Time)
	log.Println("---------------------------------------------------------------------------------------------------------")
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

func (cc *Chesscom) SaveGame() {
	if len(cc.MoveList) == 0 {
		return
	}
	cc.NumberOfGames++
	cc.Game = chess.NewGame()
	if cc.PlayWithWhite {
		cc.Game.AddTagPair("White", "Stockfish 17")
		cc.Game.AddTagPair("Black", "Anonymous")
	} else {
		cc.Game.AddTagPair("White", "Anonymous")
		cc.Game.AddTagPair("Black", "Stockfish 17")
	}
	cc.Game.AddTagPair("Result", cc.GameState)
	cc.Game.AddTagPair("Date", time.Now().Format("2006-01-02 15:04:05"))
	cc.Game.AddTagPair("Round", strconv.Itoa(cc.NumberOfGames))
	cc.Game.AddTagPair("TimeControl", cc.TimeControl)
	for _, move := range cc.MoveList {
		if err := cc.Game.MoveStr(move); err != nil {
			log.Println("Loading moves: ", err)
		}
	}
	f, err := os.OpenFile("chessbot2.pgn", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err = f.WriteString(fmt.Sprintf("%s\n\n\n", cc.Game)); err != nil {
		panic(err)
	}

	fmt.Println(cc.Game)
}

// getter functions
func (cc *Chesscom) GetPlayWithWhite() bool {
	return cc.PlayWithWhite
}
func (cc *Chesscom) GetMoveList() []string {
	return cc.MoveList
}
func (cc *Chesscom) GetBestMove() string {
	return cc.SearchResults.BestMove.String()
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
