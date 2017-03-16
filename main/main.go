package main

import (
	"github.com/wspl/arc"
)

func main() {
	go func() {
		s, _ := arc.ListenArc(":2000")
		c, _ := s.AcceptArc()
		println("Server accepted: ", c.DumpRemoteAddr())
	}()

	go func() {
		c, _ := arc.DialArc("localhost:2000")
		println("Client accepted: ", c.DumpRemoteAddr())
	}()

	for { continue }
}