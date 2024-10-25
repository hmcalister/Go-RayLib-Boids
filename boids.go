package main

import (
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

