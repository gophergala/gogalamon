package main

import (
	"bytes"
	"encoding/json"
	"log"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

func wsHandler(s *websocket.Conn) {
	log.Println("User connected from", s.RemoteAddr())
	defer log.Println("User disconnected from", s.RemoteAddr())

	var u User
	u.s = s
	u.buf = bytes.NewBuffer(nil)
	u.jEnc = json.NewEncoder(u.buf)

	u.connected = true
	defer func() {
		u.connectedMutex.Lock()
		defer u.connectedMutex.Unlock()
		u.connected = false
	}()

	u.messages = make(chan *UserMessage)

	u.keys = make(map[string]bool)
	u.keyDown = make(map[string]bool)
	u.keyUp = make(map[string]bool)

	go u.recieveMessages()
	go NewPlayerShip(&u)

	for {
		select {
		case m := <-u.messages:
			err := u.send(m)
			if err != nil {
				log.Println("User error sending message,", err)
				return
			}
		}
	}
}

type User struct {
	s    *websocket.Conn
	buf  *bytes.Buffer
	jEnc *json.Encoder

	connectedMutex sync.Mutex
	connected      bool

	messages chan *UserMessage

	keysMutex sync.Mutex
	keys      map[string]bool
	keyDown   map[string]bool
	keyUp     map[string]bool
}

type UserMessage struct {
	Event string
	Data  interface{}
}

func (u *User) send(m *UserMessage) error {
	err := u.jEnc.Encode(m)
	if err != nil {
		return err
	}

	err = u.s.SetWriteDeadline(time.Now().Add(time.Second * 5))
	if err != nil {
		return err
	}
	_, err = u.s.Write(u.buf.Bytes())
	if err != nil {
		return err
	}
	u.buf.Reset()
	return nil
}

func (u *User) recieveMessages() {
	d := json.NewDecoder(u.s)

	for {
		m := make(map[string]interface{})
		err := d.Decode(&m)
		if err != nil {
			log.Println("Error reading message from client,", err)
			return
		}
		event, ok := m["Event"].(string)
		if !ok {
			log.Println("Unable to cast event from user to string")
			return
		}
		switch {
		default:
			if event[1:] == " down" {
				u.keysMutex.Lock()
				u.keys[event[:1]] = true
				u.keyDown[event[:1]] = true
				u.keysMutex.Unlock()
			} else if event[1:] == " up" {
				u.keysMutex.Lock()
				u.keys[event[:1]] = false
				u.keyUp[event[:1]] = true
				u.keysMutex.Unlock()
			} else {
				log.Println("Unkown user message:", m)
			}
		}

	}
}

type PlayerShip struct {
	user  *User
	x     float32
	y     float32
	vx    float32
	vy    float32
	accel float32
}

func NewPlayerShip(user *User) {
	var p PlayerShip
	p.user = user
	p.accel = 0.1

	NewEntity <- &p
}

func (p *PlayerShip) update() (alive bool) {
	p.user.keysMutex.Lock()
	defer p.user.keysMutex.Unlock()

	var dx float32
	var dy float32
	if p.user.keys["a"] {
		dx -= 1
	}
	if p.user.keys["d"] {
		dx += 1
	}
	if p.user.keys["w"] {
		dy += 1
	}
	if p.user.keys["s"] {
		dy -= 1
	}
	if dy*dx != 0 {
		dx /= 1.41421356237
		dy /= 1.41421356237
	}
	{
		a := p.accel
		ra := 1 - a
		p.vx = ra*p.vx + a*dx
		p.vy = ra*p.vy + a*dy
		p.x += p.vx
		p.y += p.vy
	}
	log.Println(p.x, p.y)

	for key := range p.user.keyDown {
		p.user.keyDown[key] = false
	}
	for key := range p.user.keyUp {
		p.user.keyUp[key] = false
	}

	p.user.connectedMutex.Lock()
	defer p.user.connectedMutex.Unlock()
	return p.user.connected
}
