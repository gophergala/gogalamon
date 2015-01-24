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

	keysMutex   sync.Mutex
	keys        map[string]bool
	keyDown     map[string]bool
	keyUp       map[string]bool
	chatMessage string
	Username    string
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
		switch event {
		case "chatMessage":
			if msg, ok := m["Message"].(string); ok {
				u.keysMutex.Lock()
				u.chatMessage = msg
				u.keysMutex.Unlock()
			}
		case "username":
			if name, ok := m["User"].(string); ok {
				u.keysMutex.Lock()
				u.Username = name
				u.keysMutex.Unlock()
			}
		default:
			if event[len(event)-5:] == " down" {
				u.keysMutex.Lock()
				u.keys[event[:1]] = true
				u.keyDown[event[:1]] = true
				u.keysMutex.Unlock()
			} else if event[len(event)-3:] == " up" {
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
	user      *User
	x         float32
	y         float32
	vx        float32
	vy        float32
	accel     float32
	radius    float32
	health    int
	maxHealth int
	speed     float32
}

func NewPlayerShip(user *User) {
	var p PlayerShip
	p.user = user
	p.accel = 0.1
	p.radius = 0.1
	p.maxHealth = 100
	p.speed = 1

	NewEntity <- &p
}

func (p *PlayerShip) update(overworld *Overworld) (alive bool) {
	p.user.keysMutex.Lock()
	defer p.user.keysMutex.Unlock()

	if p.health <= 0 {
		p.health = p.maxHealth
		p.x = 0
		p.y = 0
		p.vx = 0
		p.vy = 0
	}

	var dx float32
	var dy float32
	if p.user.keys["a"] {
		dx -= p.speed
	}
	if p.user.keys["d"] {
		dx += p.speed
	}
	if p.user.keys["w"] {
		dy += p.speed
	}
	if p.user.keys["s"] {
		dy -= p.speed
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

	overworld.set(p, p.x, p.y, p.radius)

	//Collision testing code
	// log.Println("____________")
	// log.Println(p.x, p.y)
	//log.Println(overworld.query(nil, p.x, p.y, p.radius))
	// log.Println(overworld.query(nil, p.x, p.y+5, p.radius))

	if p.user.chatMessage != "" {
		log.Println(p.user.Username, ":", p.user.chatMessage)
		type Chatmsg struct {
			User    string
			Message string
		}
		msg := UserMessage{
			"chatMessage",
			Chatmsg{
				p.user.Username,
				p.user.chatMessage,
			},
		}
		for _, other := range overworld.query(p, p.x, p.y, 200) {
			if other, ok := other.(*PlayerShip); ok {
				other.user.messages <- &msg
			}
		}
		p.user.chatMessage = ""
	}

	//Keys cleanup
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

func (p *PlayerShip) damage(damage int, teamSource team) {
	p.health -= damage
}
