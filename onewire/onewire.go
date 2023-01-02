package onewire

import (
	"machine"
	"time"
)

const (
	retries = 125
	resetWait = time.Microsecond * 480
)

func reset(p machine.Pin) bool {
	// datasheet page 14 of 27
	p.Low()
	n := time.Now()
	s := time.Now().Add(resetWait)
	
	for {
	if n.Before(s) {

	}
	}
	time.Sleep(resetWait)

	return false
}
