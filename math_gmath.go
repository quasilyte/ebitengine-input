//go:build gmath

package input

import (
	"github.com/quasilyte/gmath"
)

type Vec = gmath.Vec

func vecDistance(v, v2 Vec) float64 {
	return v.DistanceTo(v2)
}

func vecDot(v, v2 Vec) float64 {
	return v.Dot(v2)
}

func vecLenSquared(v Vec) float64 {
	return v.LenSquared()
}

func vecLen(v Vec) float64 {
	return v.Len()
}

func vecAngle(v Vec) float64 {
	return float64(v.Angle())
}

func angleNormalized(radians float64) float64 {
	return float64(gmath.Rad(radians).Normalized())
}
