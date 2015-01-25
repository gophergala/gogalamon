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
		go NewPlanet(0, 2000)
		go NewPlanet(0, 9000)
		go NewPlanet(0, -9000)
		go NewPlanet(9000, 0)
		go NewPlanet(-9000, 0)
	}

	for {
		select {
		case <-ticker:
			{
				var i, place int
				for i < len(entities) {
					entity := entities[i]
					if entity.update(overworld, planets) {
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

					if msg := user.GetChatMessage(); msg != nil {
						for other := range users {
							if other != user {
								other.Send(msg)
							}
						}
					}
				}
				for i := 0; i < len(users); i++ {
					user := <-wait
					if user != nil {
						nextUsers[user] = struct{}{}
					}
				}
				users = nextUsers
			}
			{
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
	update(overworld *Overworld, planets []*Planet) (alive bool)
	RenderInfo() RenderInfo
}

type team uint

type EntityTeam interface {
	team() team
}

const (
	TeamNone = team(iota)
	TeamPirates
	TeamGophers
	TeamPythons
	TeamMax
)

func (t team) String() string {
	switch t {
	case TeamPirates:
		return "pirate"
	case TeamGophers:
		return "gopher"
	case TeamPythons:
		return "python"
	}
	return "pirate"
}

type EntityDamage interface {
	Entity
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
