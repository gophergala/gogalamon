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

func mainLoop() {
	var entities []Entity
	ticker := time.Tick(time.Second / 30)
mainloop:
	for {
		{
			var i, place int
			for i < len(entities) {
				entity := entities[i]
				if entity.update() {
					entities[place] = entity
					place++
				} else {

					//OVERWORLD REMOVE ENTITY
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
		for {
			select {
			case <-ticker:
				continue mainloop
			case entity := <-NewEntity:
				entities = append(entities, entity)
			}

		}
	}
}

type Entity interface {
	update() (alive bool)
}

type V2 [2]float32
