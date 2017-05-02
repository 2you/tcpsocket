package tcpsocket

import (
	"errors"
	"fmt"
	"net"
)

func NewServer(listenport uint) *TTcpServer {
	p := new(TTcpServer)
	p.listenPort = listenport
	return p
}

func (obj *TTcpServer) Open() error {
	var (
		vErr     error
		pTcpAddr *net.TCPAddr
		pListen  *net.TCPListener
		sPort    string
	)
	if obj.eventClientIONew == nil {
		return errors.New("未设置新建客户端IO对象事件")
	}
	if obj.eventDataBlockNew == nil {
		return errors.New("未设置新建数据包对象事件")
	}
	sPort = fmt.Sprintf(":%d", obj.listenPort)
	Debugln(`socket server listen port is`, obj.listenPort)
	if pTcpAddr, vErr = net.ResolveTCPAddr("tcp4", sPort); vErr != nil { //获取一个tcpAddr
		return vErr
	}
	if pListen, vErr = net.ListenTCP("tcp", pTcpAddr); vErr != nil {
		return vErr
	}
	go obj.accept(pListen)
	return nil
}

func (obj *TTcpServer) SetEventClientConnect(event ClientConnectEvent) {
	obj.eventClientConnect = event
}

func (obj *TTcpServer) SetEventClientIONew(event ClientIONewEvent) {
	obj.eventClientIONew = event
}

func (obj *TTcpServer) SetEventDataBlockNew(event DataBlockNewEvent) {
	obj.eventDataBlockNew = event
}

func (obj *TTcpServer) lock() {
	obj.mutex.Lock()
}

func (obj *TTcpServer) unlock() {
	obj.mutex.Unlock()
}

func (obj *TTcpServer) onClientConnect(client *TTcpClient) {
	if obj.eventClientConnect == nil {
		return
	}
	obj.eventClientConnect(client)
}

func (obj *TTcpServer) accept(listener *net.TCPListener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			Errorf("accept error is %s\n", err.Error())
			continue
		}

		func() {
			obj.lock()
			defer obj.unlock()
			//			conn.Write([]byte("hello\n"))
			//			conn.Close()
			//			return
			client := NewTcpClient(conn, obj.eventClientIONew(conn)) //新建客户端对象
			client.SetEventDataBlockNew(obj.eventDataBlockNew)
			obj.onClientConnect(client) //如果设置了eventClientConnect将会在此触发
			client.StartReadThread()
		}()
	}
}
