package util

import (
	"time"
)

var MaxRetries = 3

type Func func(count int) (noError bool)

func Retry(fn Func) bool {
	attempt := 1
	for {
		noError := fn(attempt)
		if noError == true {
			break
		}
		attempt++
		time.Sleep(10 * time.Second) // total waiting time is 10 + call function(blocked) time
		if attempt > MaxRetries {
			return false
		}
	}

	return true
}
