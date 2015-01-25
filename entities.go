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
	overworld.set(b, b.x, b.y, 8)
	return b.timeLeft > 0
}

func (b *Bullet) RenderInfo() RenderInfo {
	return RenderInfo{
		b.renderId, b.x, b.y, 0, "ball_plasma",
	}
}

type Planet struct {
	x, y     float32
	t        team
	rotation float32
	renderId int
	set      bool
}

func NewPlanet(x, y float32) {
	var p Planet
	p.x = x
	p.y = y

	p.renderId = <-NextRenderId
	NewEntity <- &p
}

func (p *Planet) update(overworld *Overworld) (alive bool) {
	if !p.set {
		overworld.set(p, p.x, p.y, 512)
		p.set = true
	}
	p.rotation += 1
	return true
}

func (p *Planet) RenderInfo() RenderInfo {
	img := "planet_python"

	return RenderInfo{
		p.renderId, p.x, p.y, p.rotation, img,
	}
}

func (p *Planet) planetInfo() PlanetInfo {
	return PlanetInfo{p.x, p.y, p.t.String()}
}

type PlanetInfo struct {
	X, Y float32
	Team string
}

type PlayerShip struct {
	transform
	user           *User
	radius         float32
	health         int
	maxHealth      int
	rotation       float32
	renderId       int
	reloadTime     int
	fullReloadTime int
}

func NewPlayerShip(user *User) {
	var p PlayerShip
	p.user = user
	p.accel = 0.8
	p.radius = 32
	p.maxHealth = 100
	p.speed = 16
	p.renderId = <-NextRenderId
	p.fullReloadTime = framesPerSecond / 5

	NewEntity <- &p
}

func (p *PlayerShip) update(overworld *Overworld) (alive bool) {
	if p.health <= 0 {
		p.health = p.maxHealth
		p.x = 0
		p.y = 0
		p.vx = 0
		p.vy = 0
	}

	var dx float32
	var dy float32
	if p.user.Key("a") {
		dx -= p.speed
	}
	if p.user.Key("d") {
		dx += p.speed
	}
	if p.user.Key("w") {
		dy -= p.speed
	}
	if p.user.Key("s") {
		dy += p.speed
	}
	if dy*dx != 0 {
		dx /= 1.41421356237
		dy /= 1.41421356237
	}
	{
		if dx != 0 || dy != 0 {
			p.rotation = float32(math.Atan2(float64(dx), float64(-1*dy))) /
				(math.Pi * 2) * 360
		}
		p.adjustV(dx, dy)
		p.applyV()
		p.user.viewX = p.x
		p.user.viewY = p.y
	}

	overworld.set(p, p.x, p.y, p.radius)

	//Collision testing code
	// log.Println("____________")
	// log.Println(p.x, p.y)
	//log.Println(overworld.query(nil, p.x, p.y, p.radius))
	// log.Println(overworld.query(nil, p.x, p.y+5, p.radius))

	if msg := p.user.GetChatMessage(); msg != nil {
		for _, other := range overworld.query(p, p.x, p.y, 200) {
			if other, ok := other.(*PlayerShip); ok {
				other.user.Send(msg)
			}
		}
	}

	p.reloadTime += 1
	if p.user.Key("f") && p.fullReloadTime < p.reloadTime {
		r := float64(p.rotation-90) / 180 * math.Pi
		vx := float32(math.Cos(r))*16 + p.vx
		vy := float32(math.Sin(r))*16 + p.vy
		p.reloadTime = 0
		go NewBullet(p.x, p.y, vx, vy, TeamGophers)
	}

	return p.user.Connected()
}

func (p *PlayerShip) RenderInfo() RenderInfo {
	return RenderInfo{
		p.renderId, p.x, p.y, p.rotation, "ship",
	}
}

func (p *PlayerShip) damage(damage int, teamSource team) {
	p.health -= damage
}
