package main

import (
	"log"

	"github.com/2you/tcpsocket"
)

type SocketActionA struct {
	tcpsocket.SocketAction
}

func NewSocketActionA() tcpsocket.SocketAction {
	p := new(SocketActionA)
	return p
}

func (this *SocketActionA) GetHeadSize() int {
	return 32
}

func (this *SocketActionA) GetBodySizeOffSet() int {
	return 4
}

func (this *SocketActionA) GetBodySizeLength() int {
	return 4
}

func (this *SocketActionA) LittleOrBig() byte {
	return 'L'
}

func (this *SocketActionA) OnRead(client *tcpsocket.ClientSocket, data []byte) {
	log.Printf("[%s] read bytes %d\n", client.RemoteAddr(), len(data))
}

func (this *SocketActionA) OnConnect(client *tcpsocket.ClientSocket) {
	log.Printf("[%s] connect\n", client.RemoteAddr())
}

func (this *SocketActionA) OnDisconnect(client *tcpsocket.ClientSocket, err error) {
	log.Printf("[%s] disconnect\n", client.RemoteAddr())
}
