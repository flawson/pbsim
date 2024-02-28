package main

import (
	"github.com/flawson/pbsim/simulator"
)

func main() {
	sim := simulator.NewSimulator()

	sim.Run(50)
	sim.PrintResults()
}
