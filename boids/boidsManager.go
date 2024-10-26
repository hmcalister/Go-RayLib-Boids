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
	BOID_INIT_VELOCITY_MAX_MAGNITUDE float32 = 4.0

	// How far two boids can be and still have an effect on one another
	MAXIMUM_BOID_VISION float32 = 128.0

	// The target proximity measure
	SEPARATION_OPTIMAL_PROXIMITY_MEASURE float32 = 0.05

	// How strongly boids are affected by the separation factor
	SEPARATION_COEFFICIENT float32 = 0.5

	// How strongly boids are affected by the alignment factor
	ALIGNMENT_COEFFICIENT float32 = 0.15

	// How strongly boids are affected by the cohesion factor
	COHESION_COEFFICIENT float32 = 0.05
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
	const windowEdgeSpringConstant float32 = 0.05
	for updateIndex := range indexChannel {
		targetBoid := currentBoids[updateIndex]

		// --------------------------------------------------------------------------------
		// Avoid flying off screen

		if targetBoid.position.X < WINDOW_EDGE_BUFFER_DISTANCE {
			targetBoid.velocity = rl.Vector2Add(targetBoid.velocity, rl.Vector2Scale(x_HAT, windowEdgeSpringConstant*(WINDOW_EDGE_BUFFER_DISTANCE-targetBoid.position.X)))
		} else if targetBoid.position.X > float32(config.WindowWidth)-WINDOW_EDGE_BUFFER_DISTANCE {
			targetBoid.velocity = rl.Vector2Add(targetBoid.velocity, rl.Vector2Scale(x_HAT, windowEdgeSpringConstant*(float32(config.WindowWidth)-WINDOW_EDGE_BUFFER_DISTANCE-targetBoid.position.X)))
		}

		if targetBoid.position.Y < WINDOW_EDGE_BUFFER_DISTANCE {
			targetBoid.velocity = rl.Vector2Add(targetBoid.velocity, rl.Vector2Scale(y_HAT, windowEdgeSpringConstant*(WINDOW_EDGE_BUFFER_DISTANCE-targetBoid.position.Y)))
		} else if targetBoid.position.Y > float32(config.WindowHeight)-WINDOW_EDGE_BUFFER_DISTANCE {
			targetBoid.velocity = rl.Vector2Add(targetBoid.velocity, rl.Vector2Scale(y_HAT, windowEdgeSpringConstant*(float32(config.WindowHeight)-WINDOW_EDGE_BUFFER_DISTANCE-targetBoid.position.Y)))
		}

		// Loop over all (other) boids and calculate the factors for update
		// Currently, check ALL boids, but perhaps changing this could improve performance...

		numProximalBoids := float32(0.0)
		centerOfMassOfProximalBoids := rl.Vector2Zero()
		for i := range len(currentBoids) {
			if i == updateIndex {
				continue
			}

			// --------------------------------------------------------------------------------'
			// Separation

			// --------------------------------------------------------------------------------'
			// Alignment

			// --------------------------------------------------------------------------------'
			// Cohesion

		}

		targetBoid.position = rl.Vector2Add(targetBoid.position, targetBoid.velocity)
		updatedBoids[updateIndex] = targetBoid
	}
}
