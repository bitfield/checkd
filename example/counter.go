package main

import (
	"time"

	"github.com/bitfield/checkd"
)

func init() {
	checkd.Every(time.Minute, checkCounter)
}

func checkCounter() {
	checkd.Counter("check_calls_total", "Number of times the checker has been called").Inc()
}
