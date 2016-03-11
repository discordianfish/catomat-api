package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/discordianfish/catomat/stepper"
	rpio "github.com/stianeikeland/go-rpio"
)

const (
	dispenseRoot = "/api/dispense/"
)

var (
	pinA       = flag.Int("pa", 18, "Pin for coil A")
	pinB       = flag.Int("pb", 23, "Pin for coil B")
	pinC       = flag.Int("pc", 24, "Pin for coil C")
	pinD       = flag.Int("pd", 25, "Pin for coil D")
	tick       = flag.Duration("t", 1*time.Millisecond, "Time to wait between steps")
	deg        = flag.Float64("m", 35, "Degree to move to open dispenser")
	duration   = flag.Duration("d", 100*time.Millisecond, "Time to keep dispensor slot open")
	listenAddr = flag.String("l", ":8080", "Address to listen on")

	motor *stepper.Motor
	lock  sync.Mutex
)

func initMotor() {
	if err := rpio.Open(); err != nil {
		log.Fatal(err)
	}
	m, err := stepper.NewMotor(*tick, *pinA, *pinB, *pinC, *pinD)
	if err != nil {
		log.Fatal(err)
	}
	motor = m
}

func main() {
	flag.Parse()
	initMotor()
	http.HandleFunc(dispenseRoot, dispenserHandler)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		log.Printf("Received %v signal. Waiting for inflight requests and cleaning up", <-c)
		lock.Lock() // Aquiring a lock to make sure no request is active
		rpio.Close()
		os.Exit(0)
	}()

	log.Fatal(http.ListenAndServe(*listenAddr, nil))
	rpio.Close()
}

func dispenserHandler(w http.ResponseWriter, r *http.Request) {
	slot := r.URL.Path[len(dispenseRoot):]
	lock.Lock()
	defer lock.Unlock()
	switch slot {
	case "a":
		motor.Clockwise(*deg)
		time.Sleep(*duration)
		motor.CounterClockwise(*deg)
	case "b":
		motor.CounterClockwise(*deg)
		time.Sleep(*duration)
		motor.Clockwise(*deg)
	default:
		http.Error(w, "Unknown dispenser slot", 404)
	}
}
