package boids

import (
	"log/slog"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/hmcalister/Go-RayLib-Boids/config"
	"golang.org/x/exp/rand"
)

const (
	WINDOW_EDGE_BUFFER_DISTANCE      float32 = 25.0
	BOID_INIT_VELOCITY_MAX_MAGNITUDE float32 = 3.0
)

var (
	x_HAT rl.Vector2 = rl.NewVector2(1.0, 0.0)
	y_HAT rl.Vector2 = rl.NewVector2(0.0, 1.0)
)

type BoidManager struct {
	Boids           []boid
	randomGenerator *rand.Rand
	config          config.Config
}

// Create a new BoidManager, which in turn makes a number of new Boids.
// randomSeed is used to initialize the Boids. If passed randomSeed is exactly 0 then a new seed is generated based on the timestamp.
// Boids are initialized randomly with both velocity and position.
func NewBoidManager(config config.Config) BoidManager {
	slog.Debug("start BoidManager initialization")

	manager := BoidManager{}
	if config.RandomSeed == 0 {
		config.RandomSeed = uint64(time.Now().UnixMicro())
		slog.Info("randomSeed set based on timestamp", "randomSeed", config.RandomSeed)
	}
	manager.randomGenerator = rand.New(rand.NewSource(config.RandomSeed))

	manager.Boids = make([]boid, config.NumBoids)
	for i := range config.NumBoids {
		manager.Boids[i] = boid{
			position: rl.NewVector2(
				manager.randomGenerator.Float32()*(float32(config.WindowWidth)-2*WINDOW_EDGE_BUFFER_DISTANCE)+WINDOW_EDGE_BUFFER_DISTANCE,
				manager.randomGenerator.Float32()*(float32(config.WindowHeight)-2*WINDOW_EDGE_BUFFER_DISTANCE)+WINDOW_EDGE_BUFFER_DISTANCE,
			),
			velocity: rl.NewVector2(
				(2*manager.randomGenerator.Float32()-1)*BOID_INIT_VELOCITY_MAX_MAGNITUDE,
				(2*manager.randomGenerator.Float32()-1)*BOID_INIT_VELOCITY_MAX_MAGNITUDE,
			),
		}

		// slog.Debug("boid initialized", "boidIndex", i, "boidData", manager.Boids[i])
	}

	manager.config = config

	return manager
}

func (manager *BoidManager) DrawBoids() {
	for _, b := range manager.Boids {
		b.draw()
	}
}

func (manager *BoidManager) TickBoids() {
	indexChannel := make(chan int)
	updatedBoids := make([]boid, len(manager.Boids))

	var workerWaitGroup sync.WaitGroup
	for range manager.config.NumWorkers {
		workerWaitGroup.Add(1)
		go func() {
			defer workerWaitGroup.Done()
			tickBoidWorkerFunc(manager.Boids, updatedBoids, manager.config, indexChannel)
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

func tickBoidWorkerFunc(currentBoids []boid, updatedBoids []boid, config config.Config, indexChannel chan int) {
	for updateIndex := range indexChannel {
		targetBoid := currentBoids[updateIndex]
		targetBoid.position = rl.Vector2Add(targetBoid.position, targetBoid.velocity)
		updatedBoids[updateIndex] = targetBoid
	}
}
