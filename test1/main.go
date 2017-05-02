package main

import (
	"fmt"
	"net"
	"time"

	"github.com/2you/tcpsocket"
)

var server *tcpsocket.TTcpServer

type tttmmm struct {
	key interface{}
}

func main() {
	m := make(map[string]string)
	m["1"] = "123"
	tcpsocket.Debugln(m["1"])
	m["1"] = "456"
	tcpsocket.Debugln(m["1"])

	tm := new(tttmmm)
	tm.key = 123
	tcpsocket.Debugln(tm.key)
	tm.key = "123"
	tcpsocket.Debugln(tm.key)

	wait := make(chan byte)
	server = tcpsocket.NewServer(11223)
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

func onClientConnect(client *tcpsocket.TTcpClient) {

}

//////////////////////////////////////////////////////////////////////////////

type D1ClientIO struct {
	conn net.Conn
	tcpsocket.IClientIO
}

func NewD1ClientIO(conn net.Conn) tcpsocket.IClientIO {
	p := new(D1ClientIO)
	p.conn = conn
	return p
}

func (obj *D1ClientIO) OnRecv(client *tcpsocket.TTcpClient, data tcpsocket.IDataBlock) {
	str := string(data.BodyContent())
	tcpsocket.Debugf("<%s>read[%s]\n", client.RemoteAddr(), str)
	if str == `12345` {
		tcpsocket.Debugf("<%s>准备阻塞\n", client.RemoteAddr())
		time.Sleep(8 * time.Second)
		tcpsocket.Debugf("<%s>退出阻塞\n", client.RemoteAddr())
	}
}

func (obj *D1ClientIO) OnConnect(client *tcpsocket.TTcpClient) {
	tcpsocket.Debugf("<%s>连接完成\n", client.RemoteAddr())
}

func (obj *D1ClientIO) OnClose(client *tcpsocket.TTcpClient, err error) {
	tcpsocket.Debugf("<%s>断开连接\n", client.RemoteAddr())
}
