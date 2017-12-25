package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"

	"github.com/2you/tcpsocket"
)

var server *tcpsocket.ServerSocket
var logger log.Logger

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	buf1 := make([]byte, 5)
	for i := 0; i < 5; i++ {
		buf1[i] = byte(i + 1)
	}
	log.Println(buf1[4:5])

	wait := make(chan byte)
	server = tcpsocket.NewServer()
	server.SetPort(12345)
	server.SetAction(NewSocketActionA())
	err := server.Open()
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Println(runtime.GOOS)

	<-wait
}
