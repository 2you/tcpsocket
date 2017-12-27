package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"time"

	"github.com/2you/gfunc"
	"github.com/2you/tcpsocket"
)

var server *tcpsocket.ServerSocket
var logger log.Logger

func fn1(v []byte) {
	log.Println(&v)
	if v == nil {
		return
	}
	log.Println("-------------")
	v[0] = 1
	v = append(v, 3)
	log.Println(&v)
}

func fn2(v *[]byte) {
	log.Println(v)
	if v == nil || *v == nil {
		return
	}
	log.Println("-------------")
	(*v)[0] = 1
	*v = append(*v, 3)
	log.Println(v)
}

func test1() {
	buf1 := make([]byte, 5)
	for i := 0; i < 5; i++ {
		buf1[i] = byte(i + 1)
	}
	log.Println(buf1[4:5])
	//	var buf2 []byte
	buf2 := make([]byte, 10)
	idx := copy(buf2, buf1)
	copy(buf2[idx:], buf1)
	log.Println(buf2)

	buf3 := gfunc.BytesMergeA(&buf1, &buf2, nil, &buf1)
	log.Println(buf3)

	buf4 := make([]byte, 3, 10)
	log.Println(len(buf4), cap(buf4))
	log.Println(buf4)

	fn1(buf4)
	log.Println(buf4)

	fn2(&buf4)
	log.Println(buf4)
}

func poolNew1() interface{} {
	ret := make([]byte, 0)
	return ret
}

var pool *sync.Pool = new(sync.Pool)

func test2() {
	for i := 0; i < 5; i++ {
		o := pool.Get().([]byte)
		fmt.Println(o)
		time.Sleep(time.Second * 1)
	}
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	go func() {
		http.ListenAndServe(":12346", nil)
	}()

	test1()

	//	pool.New = poolNew1
	//	pool.Put([]byte("1"))
	//	pool.Put([]byte("2"))
	//	pool.Put([]byte("3"))
	//	go test2()
	//	go test2()

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
