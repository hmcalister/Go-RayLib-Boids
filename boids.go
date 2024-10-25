package main

import (
	"image/color"
	"log/slog"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/exp/rand"
)

const (
	BOID_INIT_VELOCITY_MAX_MAGNITUDE float32 = 1.0
)

var (
	BOID_COLOR color.RGBA = rl.White
)

type Boid struct {
	Position rl.Vector2
	Velocity rl.Vector2
}

type BoidManager struct {
	Boids           []*Boid
	randomGenerator *rand.Rand
}

// Create a new BoidManager, which in turn makes a number of new Boids.
// randomSeed is used to initialize the Boids. If passed randomSeed is exactly 0 then a new seed is generated based on the timestamp.
// Boids are initialized randomly with both velocity and position.
func NewBoidManager(config Config) BoidManager {
	slog.Debug("start BoidManager initialization")

	manager := BoidManager{}
	if config.RandomSeed == 0 {
		config.RandomSeed = uint64(time.Now().UnixMicro())
		slog.Info("randomSeed set based on timestamp", "randomSeed", config.RandomSeed)
	}
	manager.randomGenerator = rand.New(rand.NewSource(config.RandomSeed))

	manager.Boids = make([]*Boid, config.NumBoids)
	for i := range config.NumBoids {
		manager.Boids[i] = &Boid{
			Position: rl.NewVector2(
				manager.randomGenerator.Float32()*float32(config.WindowWidth),
				manager.randomGenerator.Float32()*float32(config.WindowHeight),
			),
			Velocity: rl.NewVector2(
				(2*manager.randomGenerator.Float32()-1)*BOID_INIT_VELOCITY_MAX_MAGNITUDE,
				(2*manager.randomGenerator.Float32()-1)*BOID_INIT_VELOCITY_MAX_MAGNITUDE,
			),
		}

		slog.Debug("boid initialized", "boidIndex", i, "boid", manager.Boids[i])
	}

	return manager
}
