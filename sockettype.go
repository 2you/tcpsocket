package tcpsocket

import (
	"net"
	"sync"
	"time"
)

type DataBlockHead []byte //数据包头部类型
type DataBlockBody []byte //数据包体类型

//通讯包接口
type IDataBlock interface {
	HeadSize() int                  //数据包头部长度
	BodySize() int                  //数据包体长度
	HeadContent() DataBlockHead     //获取数据包头部内容
	BodyContent() DataBlockBody     //获取数据包体内容
	SetHeadContent(v DataBlockHead) //设置数据包头部内容
	SetBodyContent(v DataBlockBody) //设置数据包体内容
}

//客户端IO事件接口
type IClientIO interface {
	OnRecv(client *ClientSocket, data IDataBlock)
	OnConnect(client *ClientSocket)
	OnClose(client *ClientSocket, err error)
}

type SocketAction interface {
	OnRead(client *ClientSocket, data []byte)
	OnConnect(client *ClientSocket)
	OnDisconnect(client *ClientSocket, err error)
}

type ReadMethod struct {
	HeadSize       int  //数据头长度
	BodySizeOffSet int  //数据体长度开始位置
	BodySizeLen    int  //表示数据体长度字节长度
	LB             byte //
}

//tcp服务端结构体
type ServerSocket struct {
	port   uint //服务端监听端口
	action SocketAction

	eventClientIONew   ClientIONewEvent   //新建客户端IO对象事件
	eventDataBlockNew  DataBlockNewEvent  //新建数据包对象事件
	eventClientConnect ClientConnectEvent //客户端连接事件
	mutex              sync.Mutex

	onClientConnect    ClientConnectEvent
	onClientDisconnect ClientDisconnectEvent
	onClientRead       ClientReadEvent
}

//tcp客户端结构体
type ClientSocket struct {
	active bool
	action SocketAction
	host   string
	port   uint
	svrclt bool
	socket *net.TCPConn

	io                IClientIO
	eventDataBlockNew DataBlockNewEvent //新建数据包对象事件
	readThreadActive  bool
	connectTime       time.Time //连接时间

	onConnect    ClientConnectEvent
	onDisconnect ClientDisconnectEvent
	onRead       ClientReadEvent
}

type ClientConnectEvent func(client *ClientSocket)
type ClientIONewEvent func(tcpconn *net.TCPConn) IClientIO
type DataBlockNewEvent func() IDataBlock

type ClientDisconnectEvent func(client *ClientSocket, err error)
type ClientReadEvent func(client *ClientSocket, data []byte)
