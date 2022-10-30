//go:build !gmath

package input

import (
	"math"
)

// We're not using any math/vector library to make it possible for the users
// to use any kind of math library they like without having to have more
// than one math library inside their project.
//
// This means that we'll have to implement some math here,
// but it's worth it.

type Vec struct {
	X float64
	Y float64
}

func vecDot(v, v2 Vec) float64 {
	return (v.X * v2.X) + (v.Y * v2.Y)
}

func vecLenSquared(v Vec) float64 {
	return vecDot(v, v)
}

func vecLen(v Vec) float64 {
	return math.Sqrt(vecLenSquared(v))
}

func vecAngle(v Vec) float64 {
	return math.Atan2(v.Y, v.X)
}

func angleNormalized(radians float64) float64 {
	radians -= math.Floor(radians/(2*math.Pi)) * 2 * math.Pi
	return radians
}
