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
	OnRecv(client *TTcpClient, data IDataBlock)
	OnConnect(client *TTcpClient)
	OnClose(client *TTcpClient, err error)
}

//tcp服务端结构体
type TTcpServer struct {
	listenPort         uint               //服务端监听端口
	eventClientIONew   ClientIONewEvent   //新建客户端IO对象事件
	eventDataBlockNew  DataBlockNewEvent  //新建数据包对象事件
	eventClientConnect ClientConnectEvent //客户端连接事件
	mutex              sync.Mutex
}

//tcp客户端结构体
type TTcpClient struct {
	io                IClientIO
	eventDataBlockNew DataBlockNewEvent //新建数据包对象事件
	socket            net.Conn
	readThreadActive  bool
	connectTime       time.Time //连接时间
}

type ClientConnectEvent func(client *TTcpClient)
type ClientIONewEvent func(conn net.Conn) IClientIO
type DataBlockNewEvent func() IDataBlock
