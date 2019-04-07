package main

import (
	"time"

	"github.com/bitfield/checkd"
)

func init() {
	checkd.Every(10*time.Second, func() {
		checkd.Gauge("unix_time_seconds", "Current Unix time").SetToCurrentTime()
	})
}
