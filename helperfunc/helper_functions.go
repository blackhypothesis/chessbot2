package helperfunc

import (
	"errors"
	"log"
	"math/rand/v2"
	"os"
	"regexp"
	"runtime"
	"time"

	"github.com/joho/godotenv"
)

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

func waitToPlayMove(len_move_list int, time_left_seconds [2]int) {
	// artificially wait for some time to get higher standarddeviation of move time usage
	min_wait_seconds := 0.5
	max_wait_seconds := 10.0

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
	time.Sleep(time.Duration(wait_seconds) * time.Second)

}

// for performance mesurement
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
