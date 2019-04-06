package main

import (
	"time"

	"github.com/bitfield/checkd"
)

func init() {
	checkd.Register(checkTime, 10*time.Second)
}

func checkTime() {
	checkd.Gauge("unix_time_seconds", "Current Unix time").Set(float64(time.Now().UnixNano()))
}
