package main

import (
	"flag"
	"log"
	"time"

	"github.com/discordianfish/catomat/stepper"
	rpio "github.com/stianeikeland/go-rpio"
)

const (
	dispenseRoot = "/dispense/"
)

var (
	pinA = flag.Int("pa", 18, "Pin for coil A")
	pinB = flag.Int("pb", 23, "Pin for coil B")
	pinC = flag.Int("pc", 24, "Pin for coil C")
	pinD = flag.Int("pd", 25, "Pin for coil D")
	tick = flag.Duration("t", 5*time.Millisecond, "Time to wait between steps")
	deg  = flag.Float64("m", 22, "Degree to move")
	ccw  = flag.Bool("ccw", false, "Move counter-clockwise")
)

func main() {
	flag.Parse()
	if err := rpio.Open(); err != nil {
		log.Fatal(err)
	}
	m, err := stepper.NewMotor(*tick, *pinA, *pinB, *pinC, *pinD)
	if err != nil {
		log.Fatal(err)
	}
	if *ccw {
		m.CounterClockwise(*deg)
	} else {
		m.Clockwise(*deg)
	}
	rpio.Close()
}
