package main

import (
	"time"

	"github.com/bitfield/checkd"
)

func init() {
	checkd.Register(checkCounter, 1*time.Minute)
}

func checkCounter() {
	checkd.Counter("check_calls_total", "Number of times the checker has been called").Inc()
}
