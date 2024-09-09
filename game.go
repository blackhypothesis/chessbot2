package main

import (
	"strconv"
	"strings"

	"github.com/tebeka/selenium"
)

func getTimeLeftSeconds(driver selenium.WebDriver) ([2]int, error) {
	time_left, err := driver.FindElements(selenium.ByClassName, "time")
	if err != nil {
		return [2]int{0, 0}, err
	}
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

	return [2]int{time_self_secs, time_opponent_secs}, nil
}
