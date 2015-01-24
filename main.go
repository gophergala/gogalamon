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

	go mainLoop()
	err := http.ListenAndServe(":8080", nil)
	panic(err)
}

var NewEntity = make(chan Entity)

const framesPerSecond = 30

func mainLoop() {
	var entities []Entity
	overworld := NewOverworld()
	ticker := time.Tick(time.Second / framesPerSecond)

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
		case entity := <-NewEntity:
			entities = append(entities, entity)
		}

	}
}

type Entity interface {
	update(overworld *Overworld) (alive bool)
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

type EntityDamage interface {
	damage(damage int, teamSource team)
}

type V2 [2]float32
