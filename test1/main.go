package main

import (
	"fmt"
	"net"
	"time"

	"github.com/2you/tcpsocket"
	"log"
)

var server *tcpsocket.ServerSocket
var logger log.Logger

func main() {
	wait := make(chan byte)
	server = tcpsocket.NewServer(12345)
	server.SetEventClientConnect(onClientConnect)
	server.SetEventClientIONew(NewD1ClientIO)
	server.SetEventDataBlockNew(NewD1DataBlock)
	err := server.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	<-wait
}

func onClientConnect(client *tcpsocket.ClientSocket) {

}

//////////////////////////////////////////////////////////////////////////////

type D1ClientIO struct {
	conn net.Conn
	tcpsocket.IClientIO
}

func NewD1ClientIO(tcpconn *net.TCPConn) tcpsocket.IClientIO {
	p := new(D1ClientIO)
	p.conn = tcpconn
	return p
}

func (obj *D1ClientIO) OnRecv(client *tcpsocket.ClientSocket, data tcpsocket.IDataBlock) {
	str := string(data.BodyContent())
	log.Printf("<%s>read[%s]\n", client.RemoteAddr(), str)
	if str == `12345` {
		log.Printf("<%s>准备阻塞\n", client.RemoteAddr())
		time.Sleep(8 * time.Second)
		log.Printf("<%s>退出阻塞\n", client.RemoteAddr())
	}
}

func (obj *D1ClientIO) OnConnect(client *tcpsocket.ClientSocket) {
	log.Printf("<%s>连接完成\n", client.RemoteAddr())
}

func (obj *D1ClientIO) OnClose(client *tcpsocket.ClientSocket, err error) {
	log.Printf("<%s>断开连接\n", client.RemoteAddr())
}
