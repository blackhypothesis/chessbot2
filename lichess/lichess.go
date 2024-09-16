package lichess

import (
	"errors"
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

type Lichess struct {
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

func New() (*Lichess, error) {
	env, err := getENV()
	if err != nil {
		log.Fatal("Error: ", err)
	}
	return &Lichess{
		Url:         "https://lichess.org",
		UserName:    env.Login,
		Password:    env.Password,
		TimeControl: "1+0",
	}, nil
}

/*
 * START implementation of interface chessOnline
 */
func (lc *Lichess) ConnectToSite() error {
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

	log.Println("Connecting to: ", lc.Url)
	err = driver.Get(lc.Url)
	if err != nil {
		return err
	}

	lc.Service = service
	lc.Driver = driver

	return nil
}

func (lc *Lichess) ServiceStop() {
	lc.Service.Stop()
}

func (lc *Lichess) SignIn() error {
	log.Println("Login: ", lc.UserName, "  Password: ", lc.Password)
	sign_in, err := lc.Driver.FindElement(selenium.ByClassName, "signin")
	if err != nil {
		return err
	}
	sign_in.Click()

	err = lc.Driver.WaitWithTimeout(func(driver selenium.WebDriver) (bool, error) {
		form_username, _ := driver.FindElement(selenium.ByCSSSelector, ".post:nth-child(60)")
		if form_username != nil {
			return form_username.IsDisplayed()
		}
		return false, nil
	}, 2*time.Second)
	if err != nil {
		log.Fatal("Error:", err)
	}

	form_username, err := lc.Driver.FindElement(selenium.ByID, "form3-username")
	if err != nil {
		return err
	}
	form_password, err := lc.Driver.FindElement(selenium.ByID, "form3-password")
	if err != nil {
		return err
	}
	button_submit, err := lc.Driver.FindElement(selenium.ByClassName, "submit")
	if err != nil {
		return err
	}
	log.Println("Login with: ", lc.UserName, lc.Password)
	form_username.Clear()
	form_username.SendKeys(lc.UserName)
	time.Sleep(1 * time.Second)

	form_password.Clear()
	form_password.SendKeys(lc.Password)
	time.Sleep(1 * time.Second)
	button_submit.Click()

	return nil
}

func (lc *Lichess) PlayWithHuman() error {
	err := lc.Driver.WaitWithTimeout(func(driver selenium.WebDriver) (bool, error) {
		time_selectors, _ := lc.Driver.FindElements(selenium.ByClassName, "clock")
		if time_selectors != nil {
			return time_selectors[0].IsDisplayed()
		}
		return false, nil
	}, 2*time.Second)
	if err != nil {
		return err
	}
	time_selectors, err := lc.Driver.FindElements(selenium.ByClassName, "clock")
	if err != nil {
		return err
	}
	switch lc.TimeControl {
	case "1+0":
		time_selectors[0].Click()
	case "2+1":
		time_selectors[1].Click()
	case "3+0":
		time_selectors[2].Click()
	case "3+2":
		time_selectors[3].Click()
	default:
		return errors.New("timecontrol does not exist")
	}
	return nil
}

func (lc *Lichess) PlayWithComputer() error {
	// Button [PLAY WITH COMPUTER]
	button, err := lc.Driver.FindElement(selenium.ByClassName, "config_ai")
	if err != nil {
		return err
	}
	button.Click()
	time.Sleep(500 * time.Millisecond)

	// Button Strength [8]
	level, err := lc.Driver.FindElement(selenium.ByXPATH, "/html/body/div/main/div[1]/dialog/div[2]/div/div/div[3]/div[1]/group/div[8]/label")
	if err != nil {
		return err
	}
	level.Click()
	time.Sleep(500 * time.Millisecond)

	// Dropdown Time Control
	tc, err := lc.Driver.FindElement(selenium.ByID, "sf_timeMode")
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
	bw, err := lc.Driver.FindElement(selenium.ByXPATH, "/html/body/div/main/div[1]/dialog/div[2]/div/div/div[4]/button[2]")
	if err != nil {
		return err
	}
	bw.Click()
	return nil
}

func (lc *Lichess) NewGame() {
	lc.Game = chess.NewGame()
}

func (lc *Lichess) IsPlayWithWhite() {
	board_coords_class := ""
	for {
		err := lc.Driver.WaitWithTimeout(func(driver selenium.WebDriver) (bool, error) {
			board_coords, _ := driver.FindElement(selenium.ByTagName, "coords")
			if board_coords != nil {
				return board_coords.IsDisplayed()
			}
			return false, nil
		}, 2*time.Second)

		if err != nil {
			log.Println("IsPlayWithWhite: can't get board orientation")
		}

		board_coords, err := lc.Driver.FindElement(selenium.ByTagName, "coords")
		if err != nil {
			continue
		}
		board_coords_class, err = board_coords.GetAttribute("class")
		if err != nil {
			continue
		}
		break
	}
	if board_coords_class == "ranks black" {
		lc.PlayWithWhite = false
	} else {
		lc.PlayWithWhite = true
	}
}

func (lc *Lichess) UpdateMoveList() func() {
	defer TimeTrack(time.Now())
	move_list := []string{}
	last_move_list_len := 0

	return func() {
		move_list_container, err := lc.Driver.FindElements(selenium.ByTagName, "kwdb")
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
		lc.MoveList = move_list
	}
}

func (lc *Lichess) IsMyTurn(is_white_orientation bool) bool {
	// we play with white
	if len(lc.MoveList)%2 == 0 && is_white_orientation {
		return true
	}
	// we play with black
	if len(lc.MoveList)%2 == 1 && !is_white_orientation {
		return true
	}
	return false
}

func (lc *Lichess) CalculateEngineBestMove() error {
	// defer TimeTrack(time.Now())

	lc.Game = chess.NewGame()
	eng, err := uci.New("stockfish")
	if err != nil {
		return err
	}
	// initialize uci with new game
	if err := eng.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
		return err
	}
	defer eng.Close()

	for _, move := range lc.MoveList {
		if err := lc.Game.MoveStr(move); err != nil {
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
	if len(lc.MoveList) > 60 {
		depth = 16
	}
	cmdPos := uci.CmdPosition{Position: lc.Game.Position()}
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

	lc.BestMove = move
	return nil
}

func (lc *Lichess) CalculateTimeLeftSeconds() error {
	time_left, err := lc.Driver.FindElements(selenium.ByClassName, "time")
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

		lc.TimeLeftSeconds = [2]int{time_self_secs, time_opponent_secs}
	}
	return nil
}

func (lc *Lichess) PlayMoveWithMouse() (func(move string, len_move_list int, time_left_seconds [2]int), error) {
	defer TimeTrack(time.Now())
	cg_board, err := lc.Driver.FindElement(selenium.ByTagName, "cg-board")
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
		if !lc.PlayWithWhite {
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
		// calculate the field to click, to promote the pawn to the desired piece
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

func (lc *Lichess) GetGameState() string {
	game_state, err := lc.Driver.FindElement(selenium.ByClassName, "result")
	if err != nil {
		lc.GameState = "ongoing"
		return lc.GameState
	}
	state, err := game_state.Text()
	if err != nil {
		lc.GameState = "unknown"
		return lc.GameState
	}
	lc.GameState = state
	return lc.GameState
}

func (lc *Lichess) NewOpponent() error {
	new_opponent, err := lc.Driver.FindElement(selenium.ByXPATH, `//*[@id="main-wrap"]/main/div[1]/div[5]/div/a[1]`) // New opponent
	if err != nil {
		return err
	}
	new_opponent.Click()
	return nil
}

// getter functions
func (lc *Lichess) GetPlayWithWhite() bool {
	return lc.PlayWithWhite
}
func (lc *Lichess) GetMoveList() []string {
	return lc.MoveList
}
func (lc *Lichess) GetBestMove() string {
	return lc.BestMove.String()
}
func (lc *Lichess) GetTimeLeftSeconds() [2]int {
	return lc.TimeLeftSeconds
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
