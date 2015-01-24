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
	u.disconnect = make(chan struct{}, 1) //Buffer should be size of things
	//which can cause disconnect

	go u.recieveMessages()

	for {
		select {
		case m := <-u.messages:
			err := u.send(m)
			if err != nil {
				log.Println("User error sending message,", err)
				return
			}
		case <-u.disconnect:
			return
		}
	}
}

type User struct {
	s    *websocket.Conn
	buf  *bytes.Buffer
	jEnc *json.Encoder

	connectedMutex sync.Mutex
	connected      bool

	messages   chan *UserMessage
	disconnect chan struct{}
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
	defer func() {
		u.disconnect <- struct{}{}
	}()
	d := json.NewDecoder(u.s)

	for {
		m := make(map[string]interface{})
		err := d.Decode(&m)
		if err != nil {
			log.Println("Error reading message from client,", err)
			return
		}
		log.Println("User message:", m)
	}
}

type PlayerShip struct {
	user *User
}

func (p *PlayerShip) update() (alive bool) {

	p.user.connectedMutex.Lock()
	defer p.user.connectedMutex.Unlock()
	return p.user.connected
}
