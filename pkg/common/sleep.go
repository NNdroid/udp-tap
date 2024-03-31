package common

import (
	"math"
	"time"
)

// SleepBeforeConnect How much time to sleep on trying to connect to decoys to prevent overwhelming them
func SleepBeforeConnect(attempt int) (waitTime <-chan time.Time) {
	if attempt >= 2 { // return nil for first 2 attempts
		waitTime = time.After(time.Second *
			time.Duration(math.Pow(3, float64(attempt-1))))
	}
	return
}
