package main

import (
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"

	"net/http"
)
import "golang.org/x/net/websocket"

func main() {
	log.Println("Starting gogalamon server")

	err := readImageSizes()
	if err != nil {
		panic(err)
	}
	log.Println(imageSizes)

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
	err = http.ListenAndServe(":8080", nil)
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
				nextUsers := make(map[*User]struct{})
				wait := make(chan *User)
				for user := range users {
					go user.render(overworld, wait)
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

var imageSizes = make(map[string]float32)

func readImageSizes() error {
	infos, err := ioutil.ReadDir("static/img/")
	if err != nil {
		return err
	}
	for _, info := range infos {
		if !info.IsDir() {
			name := info.Name()
			file, err := os.Open("static/img/" + name)
			defer file.Close()
			if err != nil {
				return err
			}
			img, err := png.Decode(file)
			if err != nil {
				return err
			}
			imageSizes[name[:len(name)-4]] = float32(img.Bounds().Dx()) / 2
		}
	}
	return nil
}
