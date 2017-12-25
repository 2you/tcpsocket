package tcpsocket

import (
	"net"
	"sync"
	"time"
)

type SocketAction interface {
	GetHeadSize() int       //数据头长度
	GetBodySizeOffSet() int //数据体长度开始位置
	GetBodySizeLength() int //表示数据体长度字节长度
	LittleOrBig() byte      //小端表示还是大端表示

	OnRead(client *ClientSocket, data []byte)
	OnConnect(client *ClientSocket)
	OnDisconnect(client *ClientSocket, err error)
}

//tcp服务端结构体
type ServerSocket struct {
	port   uint //服务端监听端口
	action SocketAction
	mutex  sync.Mutex
}

//tcp客户端结构体
type ClientSocket struct {
	wbLock           sync.Mutex
	active           bool
	action           SocketAction
	host             string
	port             uint
	svrclt           bool
	socket           *net.TCPConn
	readThreadActive bool
	connectTime      time.Time //连接时间
}
