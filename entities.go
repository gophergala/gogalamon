package main

import "math"

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
	if dv > t.accel {
		ax = ax / dv * t.accel
		ay = ay / dv * t.accel
	}
	t.vx += ax
	t.vy += ay
}

func (t *transform) applyV() {
	t.x += t.vx
	t.y += t.vy
}

func NewBullet(x, y, vx, vy float32, t team) {
	var b Bullet
	b.x = x
	b.y = y
	b.vx = vx
	b.vy = vy
	b.timeLeft = framesPerSecond * 3

	b.renderId = <-NextRenderId
	NewEntity <- &b
}

type Bullet struct {
	transform
	t        team
	renderId int
	timeLeft int
}

func (b *Bullet) update(overworld *Overworld) (alive bool) {
	b.applyV()
	b.timeLeft -= 1
	overworld.set(b, b.x, b.y, 1)
	return b.timeLeft > 0
}

func (b *Bullet) RenderInfo() RenderInfo {
	return RenderInfo{
		b.renderId, b.x, b.y, 0, "ball_plasma",
	}
}
