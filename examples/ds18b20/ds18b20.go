package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/onewire"
)

func main() {
	d := onewire.NewOneWire(machine.D2)
	for i := 0; i < 10; i++ {
		println(d.SearchAddress())
		time.Sleep(time.Second * 1)
	}
}