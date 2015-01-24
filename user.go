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
}

type User struct {
	s    *websocket.Conn
	buf  *bytes.Buffer
	jEnc *json.Encoder

	connectedMutex sync.Mutex
	connected      bool
}

func (u *User) send(event string, data interface{}) error {
	type Message struct {
		event string
		data  interface{}
	}

	m := Message{event, data}

	err := u.jEnc.Encode(&m)
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

type PlayerShip struct {
	user *User
}

func (p *PlayerShip) update() (alive bool) {
	p.user.connectedMutex.Lock()
	defer p.user.connectedMutex.Unlock()
	return p.user.connected
}
