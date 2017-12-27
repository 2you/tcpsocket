package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"github.com/2you/gfunc"
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
	go func() {
		http.ListenAndServe(":12346", nil)
	}()

	//	var buf2 []byte
	buf2 := make([]byte, 10)
	idx := copy(buf2, buf1)
	copy(buf2[idx:], buf1)
	log.Println(buf2)

	buf3 := gfunc.BytesMerge(&buf1, &buf2, nil, &buf1)
	log.Println(buf3)

	buf4 := make([]byte, 3, 10)
	log.Println(len(buf4), cap(buf4))
	log.Println(buf4)

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
