package main

import (
	"fmt"
	"log"

	"github.com/2you/tcpsocket"
)

var server *tcpsocket.ServerSocket
var logger log.Logger

func main() {
	wait := make(chan byte)
	server = tcpsocket.NewServer()
	server.SetPort(12345)
	server.SetAction(NewSocketActionA())
	err := server.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	<-wait
}
