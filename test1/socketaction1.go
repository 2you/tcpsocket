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
	size := len(data)
	rmAddr := client.RemoteAddr()
	bodySize, _ := client.GetBodySize(data[:this.GetHeadSize()])
	log.Printf("[%s] read bytes %d >>> body size %d\n", rmAddr, size, bodySize)
	err := client.Write(data, size)
	if err == nil {
		log.Printf("write %d bytes to [%s] succ", size, rmAddr)
	} else {
		log.Printf("write %d bytes to [%s] error [%s]", size, rmAddr, err.Error())
	}
}

func (this *SocketActionA) OnConnect(client *tcpsocket.ClientSocket) {
	log.Printf("[%s] connect\n", client.RemoteAddr())
}

func (this *SocketActionA) OnDisconnect(client *tcpsocket.ClientSocket, err error) {
	log.Printf("[%s] disconnect\n", client.RemoteAddr())
}
