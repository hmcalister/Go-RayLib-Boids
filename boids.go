package main

import (
	"log/slog"
	"time"

	"golang.org/x/exp/rand"
)

type CartesianVector2 struct {
	X float64
	Y float64
}

type Boid struct {
	Position CartesianVector2
	Velocity CartesianVector2
}

type BoidManager struct {
	Boids           []*Boid
	randomGenerator *rand.Rand
}

// Create a new BoidManager, which in turn makes a number of new Boids.
// randomSeed is used to initialize the Boids. If passed randomSeed is exactly 0 then a new seed is generated based on the timestamp.
// Boids are initialized randomly with both velocity and position.
func NewBoidManager(numBoids int, randomSeed uint64) BoidManager {
	slog.Debug("start BoidManager initialization")

	manager := BoidManager{}
	if randomSeed == 0 {
		randomSeed = uint64(time.Now().UnixMicro())
		slog.Info("randomSeed set based on timestamp", "randomSeed", randomSeed)
	}
}
