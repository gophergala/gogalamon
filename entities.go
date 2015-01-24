package main

import (
	"math"
)

type transform struct {
	x, y   float32
	vx, vy float32
	accel  float32
	speed  float32
}

func (t *transform) adjustV(vx, vy float32) {
	ax := vx - t.vx
	ay := vy - t.vy
	dv := float32(math.Sqrt(float64(ax*ax + ay*ay)))
	if dv2 > t.accel {
		ax = ax / dv * t.accel
		ay = ay / dv * t.accel
	}
	t.vx += ax
	t.vy += ay
}
