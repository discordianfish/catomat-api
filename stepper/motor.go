package stepper

import (
	"errors"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

var (
	defaultCycle = [][]int{
		{0},
		{0, 1},
		{1},
		{1, 2},
		{2},
		{2, 3},
		{3},
		{3, 0},
	}
	ErrNotEnoughPins = errors.New("Not enough pins for cycle")
)

type Motor struct {
	cycle          [][]int
	pause          time.Duration
	cyclesPerRound int
}

func NewMotor(pause time.Duration, pins ...int) (*Motor, error) {
	for _, p := range pins {
		pin := rpio.Pin(p)
		pin.Output()
	}

	cycle := make([][]int, len(defaultCycle))
	for i, state := range defaultCycle {
		cycle[i] = make([]int, len(state))
		for j, pi := range state {
			if pi >= len(pins) {
				return nil, ErrNotEnoughPins
			}
			cycle[i][j] = pins[pi]
		}
	}
	return &Motor{
		cycle:          cycle,
		pause:          pause,
		cyclesPerRound: 512,
	}, nil
}

func (m *Motor) Clockwise(dec float64) {
	m.Move(int(dec/(360/float64(m.cyclesPerRound))), false)
}

func (m *Motor) CounterClockwise(dec float64) {
	m.Move(int(dec/(360/float64(m.cyclesPerRound))), true)
}

func (m *Motor) Move(cycles int, ccw bool) {
	cycle := make([][]int, len(m.cycle))
	copy(cycle, m.cycle)
	if ccw {
		cycle = reverse(cycle)
	}

	for i := 0; i <= cycles; i++ {
		for _, state := range cycle {
			for _, p := range state {
				pin := rpio.Pin(p)
				pin.High()
			}
			time.Sleep(m.pause)
			for _, p := range state {
				pin := rpio.Pin(p)
				pin.Low()
			}
		}
	}
}

func reverse(in [][]int) [][]int {
	ln := len(in)
	for i := 0; i <= ln/2; i++ {
		in[ln-i-1], in[i] = in[i], in[ln-i-1]
	}
	return in
}
