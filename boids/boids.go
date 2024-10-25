package boids

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	BOID_INIT_VELOCITY_MAX_MAGNITUDE float32 = 1.0
)

var (
	BOID_COLOR color.RGBA = rl.White
)

type boid struct {
	position rl.Vector2
	velocity rl.Vector2
}

func (b *boid) drawBoid(sideLength, angle float32) {
	side := rl.Vector2Scale(rl.Vector2Normalize(b.velocity), sideLength)
	v2 := rl.Vector2Add(rl.Vector2Rotate(side, -angle), b.position)
	v3 := rl.Vector2Add(rl.Vector2Rotate(side, angle), b.position)
	rl.DrawTriangle(rl.Vector2Add(b.position, side), v2, v3, BOID_COLOR)

	// rl.DrawCircle(int32(b.Position.X), int32(b.Position.Y), 3, rl.Blue)
}
