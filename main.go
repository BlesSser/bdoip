package main

import "bdoip/server"

func main() {
	s, _ := server.CreateServer()
	s.Start()
}
