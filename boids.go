package main

import (
	"image/color"
	"log/slog"
	"sync"
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

func (b *Boid) DrawBoid(sideLength, angle float32) {
	side := rl.Vector2Scale(rl.Vector2Normalize(b.Velocity), sideLength)
	v2 := rl.Vector2Add(rl.Vector2Rotate(side, -angle), b.Position)
	v3 := rl.Vector2Add(rl.Vector2Rotate(side, angle), b.Position)
	rl.DrawTriangle(rl.Vector2Add(b.Position, side), v2, v3, BOID_COLOR)

	// rl.DrawCircle(int32(b.Position.X), int32(b.Position.Y), 3, rl.Blue)
}

type BoidManager struct {
	Boids           []Boid
	randomGenerator *rand.Rand
	config          Config
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

	manager.Boids = make([]Boid, config.NumBoids)
	for i := range config.NumBoids {
		manager.Boids[i] = Boid{
			Position: rl.NewVector2(
				manager.randomGenerator.Float32()*float32(config.WindowWidth),
				manager.randomGenerator.Float32()*float32(config.WindowHeight),
			),
			Velocity: rl.NewVector2(
				(2*manager.randomGenerator.Float32()-1)*BOID_INIT_VELOCITY_MAX_MAGNITUDE,
				(2*manager.randomGenerator.Float32()-1)*BOID_INIT_VELOCITY_MAX_MAGNITUDE,
			),
		}

		slog.Debug("boid initialized", "boidIndex", i, "boidData", manager.Boids[i])
	}

	manager.config.NumWorkers = config.NumWorkers

	return manager
}

func (manager *BoidManager) TickBoids() {
	indexChannel := make(chan int)
	updatedBoids := make([]Boid, len(manager.Boids))

	var workerWaitGroup sync.WaitGroup
	for range manager.config.NumWorkers {
		workerWaitGroup.Add(1)
		go func() {
			defer workerWaitGroup.Done()
			tickBoidWorkerFunc(manager.Boids, updatedBoids, indexChannel)
		}()
	}

	for i := range len(manager.Boids) {
		indexChannel <- i
	}
	close(indexChannel)
	workerWaitGroup.Wait()

	// After workers are complete, updatedBoids contains the ticked boids.
	// We can replace the public Boids list immediately

	manager.Boids = updatedBoids
}

func tickBoidWorkerFunc(currentBoids []Boid, updatedBoids []Boid, indexChannel chan int) {
	for updateIndex := range indexChannel {
		targetBoid := currentBoids[updateIndex]
		targetBoid.Position = rl.Vector2Add(targetBoid.Position, targetBoid.Velocity)
		updatedBoids[updateIndex] = targetBoid
	}
}
