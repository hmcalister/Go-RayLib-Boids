package main

import "math"

type CartesianVector2 struct {
	X float64
	Y float64
}

func (v CartesianVector2) Magnitude() float64 {
	return math.Sqrt(v.DotProduct(v))
}

func (v CartesianVector2) Normalize() CartesianVector2 {
	mag := v.Magnitude()
	return CartesianVector2{
		v.X / mag,
		v.Y / mag,
	}
}

