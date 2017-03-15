package main

import (
	"github.com/wspl/arc"
)

func main() {
	go func() {
		s, _ := arc.ListenArc(":2000")
		c, _ := s.AcceptArc()
		println("Accepted:", c)
	}()

	go func() {
		arc.DialArc("localhost:2000")
	}()

	for { continue }
}