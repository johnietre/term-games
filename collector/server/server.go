package server

import (
	"log"
	"net"
	"sync"
)

func Run() {
	log.SetFlags(0)
}

type Game struct {
	width  int
	height int

	players []Player
	mtx     sync.RWMutex
}

type Player struct {
	id   rune
	conn net.Conn
}

type Move struct {
	playerNum rune
}
