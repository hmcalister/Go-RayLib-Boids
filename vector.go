package main

import "math"

type CartesianVector2 struct {
	X float64
	Y float64
}

func (v CartesianVector2) Magnitude() float64 {
	return math.Sqrt(v.DotProduct(v))
}

