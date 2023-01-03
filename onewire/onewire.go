package onewire

import (
	"encoding/binary"
	"errors"
	"machine"
	"time"
)

const (
	resetPulse = time.Microsecond * 480
)

type Device interface {
	Tx(out []byte) error
	SearchAddress() (uint8, error)
}

type device struct {
	pin machine.Pin
	address uint64
}

func NewOneWire(SensorPin machine.Pin) Device {
	return &device{
		pin: SensorPin,
	}
}

func (d *device) reset(p machine.Pin) error {
	// I/O SIGNALING datasheet page 13 of 17
	// p.Configure(machine.PinConfig{Mode: machine.PinOutput})
	p.Configure(machine.PinConfig{Mode: machine.PinOutput})
	p.Low()

	// Bus master drives line low for 480 µs
	time.Sleep(resetPulse)
	p.Configure(machine.PinConfig{Mode: machine.PinInput})

	// ds18b20 pulls high after reset pulse
	if !p.Get() {
		return errors.New("Pin wasn't pulled high before reset time period")
	}

	// TODO monitor time window below and break once met
	// Ds18b20 waits 15-60 µs then pulls low
	time.Sleep(time.Microsecond * 60)
	if p.Get() {
		return errors.New("Onewire device didn't pull low after 60 µs")
	}

	// TODO monitor time window below and break once met
	// presence pulse is 60-240 µs
	time.Sleep(time.Microsecond * 240)
	
	if !p.Get() {
		return errors.New("onewire reset failed, pin wasn't pulled back to high after presence pulse")
	}
	return nil
}

// TODO Finish search address
func (d *device) SearchAddress() (uint8, error) {
	err := d.Tx([]byte{readRom})
	if err != nil {
		return 0, err
	}
	return d.Rx(), nil
}

func (d *device) Rx() uint8 {
	_, m := d.pin.PortMaskSet()
	return m
}

func (d *device) Tx(out []byte) error {
	ww := make([]byte, 9, len(out)+9)
	ww[0] = matchRom
	err := d.reset(d.pin)
	if err != nil {
		return err
	}
	binary.LittleEndian.PutUint64(ww, d.address)
	b := binary.LittleEndian.Uint64(ww)
	p, _ := d.pin.PortMaskSet()
	p.Set(uint8(b))
	return nil
}
