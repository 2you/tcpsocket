package main

import (
	"encoding/binary"
	"log"
	"sync"
	"time"

	"github.com/2you/tcpsocket"
)

type SocketActionA struct {
	tcpsocket.SocketAction
	clientMap map[*tcpsocket.ClientSocket]*TestClient
	mutex     sync.Mutex
}

func NewSocketActionA() tcpsocket.SocketAction {
	p := new(SocketActionA)
	p.clientMap = make(map[*tcpsocket.ClientSocket]*TestClient)
	return p
}

func (this *SocketActionA) GetHeadSize() int {
	return 32 //数据包头部长度32字节
}

func (this *SocketActionA) GetBodySizeOffSet() int {
	return 4 //数据包头部第4位开始记录包体长度
}

func (this *SocketActionA) GetBodySizeLength() int {
	return 4 //用4个字节表示包体长度
}

func (this *SocketActionA) LittleOrBig() byte {
	return 'L'
}

func (this *SocketActionA) OnRead(client *tcpsocket.ClientSocket, data []byte) {
	//	size := len(data)
	//	rmAddr := client.RemoteAddr()
	//	bodySize, _ := client.GetBodySize(data[:this.GetHeadSize()])
	//	log.Printf("[%s] read bytes %d >>> body size %d\n", rmAddr, size, bodySize)
	//	err := client.Write(data, size)
	//	if err == nil {
	//		log.Printf("write %d bytes to [%s] succ\n", size, rmAddr)
	//	} else {
	//		log.Printf("write %d bytes to [%s] error [%s]\n", size, rmAddr, err.Error())
	//	}
	//	data = nil
	this.mutex.Lock()
	testcl := this.clientMap[client]
	this.mutex.Unlock()

	if testcl != nil {
		testcl.ReadTest(data)
	} else {
		log.Printf("[%s] deleted\n", client.RemoteAddr())
	}
}

func (this *SocketActionA) OnConnect(client *tcpsocket.ClientSocket) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	testcl := NewTestClient(client)
	this.clientMap[client] = testcl
	go testcl.WriteTest()
	log.Printf("[%s] connect\n", client.RemoteAddr())
}

func (this *SocketActionA) OnDisconnect(client *tcpsocket.ClientSocket, err error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	delete(this.clientMap, client)
	log.Printf("[%s] disconnect\n", client.RemoteAddr())
}

////////////////////////////////////////////////////////////////

type TestClient struct {
	isWait  bool
	currIdx uint64
	cmdId   uint32
	waitSec int
	client  *tcpsocket.ClientSocket
}

func NewTestClient(client *tcpsocket.ClientSocket) *TestClient {
	p := new(TestClient)
	p.client = client
	p.isWait = false
	p.waitSec = 0
	p.currIdx = 0
	p.cmdId = 8192
	return p
}

func (this *TestClient) WriteTest() {
	for {
		if !this.client.Active() {
			break
		}

		if this.isWait {
			time.Sleep(time.Second * 3)
			this.waitSec++
			if this.waitSec > 60 {
				log.Printf("[%s] wait id %d greater than %d sec\n", this.client.RemoteAddr(), this.currIdx, this.waitSec)
				if this.waitSec > 180 {
					this.client.Close()
					log.Printf("[%s] wait id %d greater than %d sec [manual close client]\n",
						this.client.RemoteAddr(), this.currIdx, this.waitSec)
					return
				}
			}
			continue
		}

		this.waitSec = 0
		this.isWait = true
		//8192 + 32 = 8224
		buf := make([]byte, 8224)
		binary.LittleEndian.PutUint32(buf[4:], 8192)
		binary.LittleEndian.PutUint32(buf[10:], this.cmdId)
		binary.LittleEndian.PutUint64(buf[14:], this.currIdx)

		if err := this.client.WriteBuf(buf); err != nil {
			this.client.Close()
			log.Println("[%s] write test error [%s]\n", this.client.RemoteAddr(), err.Error())
			return
		}
	}
}

func (this *TestClient) ReadTest(data []byte) {
	cmdId := binary.LittleEndian.Uint32(data[10:])
	if cmdId != this.cmdId {
		return
	}
	currIdx := binary.LittleEndian.Uint64(data[14:])
	if currIdx != this.currIdx {
		this.client.Close()
		log.Println("[%s] read idx not same as write idx\n", this.client.RemoteAddr())
		return
	}
	this.currIdx++
	this.isWait = false
}
