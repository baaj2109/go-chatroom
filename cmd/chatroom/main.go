package main

import (
	"fmt"
	"log"
	"net/http"

	_ "net/http/pprof"

	"chatroom/global"
	"chatroom/server"
)

var (
	addr   = ":2022"
	banner = `
    ____ .    .       _______
   |     |    |   /\     |
   |     |____|  /  \    | 
   |     |    | /----\   |
   |____ |    |/      \  |
go chatroom
`
)

func init() {
	global.Init()
}

func main() {
	fmt.Printf(banner, addr)

	server.RegisterHandle()

	log.Fatal(http.ListenAndServe(addr, nil))
}
