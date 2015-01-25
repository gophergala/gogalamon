package main

import (
	"log"
	"runtime"
	"time"

	"net/http"
)
import "golang.org/x/net/websocket"

func main() {
	log.Println("Starting gogalamon server")

	runtime.GOMAXPROCS(runtime.NumCPU())

	http.Handle("/", http.FileServer(http.Dir("static/")))
	http.Handle("/sock/", websocket.Handler(wsHandler))

	go func() {
		var i int
		for {
			NextRenderId <- i
			i++
		}
	}()

	go mainLoop()
	err := http.ListenAndServe(":8080", nil)
	panic(err)
}

var NewEntity = make(chan Entity)
var NewUser = make(chan *User)

const framesPerSecond = 30

func mainLoop() {
	var entities []Entity
	users := make(map[*User]struct{})
	overworld := NewOverworld()
	ticker := time.Tick(time.Second / framesPerSecond)

	var planets []*Planet
	{

		go NewPlanet(0, 0)
		go NewPlanet(200, 0)
		go NewPlanet(0, 1000)
		go NewPlanet(9000, 9000)
		go NewPlanet(-9000, -9000)
		go NewPlanet(9000, -9000)

	}

	for {
		select {
		case <-ticker:
			{
				var i, place int
				for i < len(entities) {
					entity := entities[i]
					if entity.update(overworld) {
						entities[place] = entity
						place++
					} else {
						overworld.remove(entity)
					}
					i++
				}
				lastPlace := place
				for place < len(entities) {
					entities[place] = nil
					place++
				}
				entities = entities[:lastPlace]
			}
			{
				planetInfos := make([]PlanetInfo, len(planets))
				for i := range planetInfos {
					planetInfos[i] = planets[i].planetInfo()
				}
				nextUsers := make(map[*User]struct{})
				wait := make(chan *User)
				for user := range users {
					go user.render(overworld, planetInfos, wait)
				}
				for i := 0; i < len(users); i++ {
					user := <-wait
					if user != nil {
						nextUsers[user] = struct{}{}
					}
				}
				users = nextUsers
			}
		case entity := <-NewEntity:
			entities = append(entities, entity)
			if planet, ok := entity.(*Planet); ok {
				planets = append(planets, planet)
			}
		case user := <-NewUser:
			users[user] = struct{}{}
		}
	}
}

type Entity interface {
	update(overworld *Overworld) (alive bool)
	RenderInfo() RenderInfo
}

type team uint

type EntityTeam interface {
	team() team
}

const (
	TeamPirates = team(iota)
	TeamGophers
	TeamPythons
)

func (t team) String() string {
	switch t {
	case TeamPirates:
		return "Pirate"
	case TeamGophers:
		return "Gopher"
	case TeamPythons:
		return "Python"
	}
	return "Pirate"
}

type EntityDamage interface {
	damage(damage int, teamSource team)
}

type V2 [2]float32

var NextRenderId = make(chan int)

type RenderInfo struct {
	I int
	X float32
	Y float32
	R float32
	N string //name
}
