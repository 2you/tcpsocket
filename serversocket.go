package tcpsocket

import (
	"errors"
	"fmt"
	"net"
)

func NewServer(listenport uint) *ServerSocket {
	p := new(ServerSocket)
	p.listenPort = listenport
	return p
}

func (obj *ServerSocket) Open() error {
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
	debugln(`socket server listen port is`, obj.listenPort)
	if pTcpAddr, vErr = net.ResolveTCPAddr("tcp4", sPort); vErr != nil { //获取一个tcpAddr
		return vErr
	}
	if pListen, vErr = net.ListenTCP("tcp", pTcpAddr); vErr != nil {
		return vErr
	}
	go obj.accept(pListen)
	return nil
}

func (obj *ServerSocket) SetEventClientConnect(event ClientConnectEvent) {
	obj.eventClientConnect = event
}

func (obj *ServerSocket) SetEventClientIONew(event ClientIONewEvent) {
	obj.eventClientIONew = event
}

func (obj *ServerSocket) SetEventDataBlockNew(event DataBlockNewEvent) {
	obj.eventDataBlockNew = event
}

func (obj *ServerSocket) lock() {
	obj.mutex.Lock()
}

func (obj *ServerSocket) unlock() {
	obj.mutex.Unlock()
}

func (obj *ServerSocket) onClientConnect(client *ClientSocket) {
	if obj.eventClientConnect == nil {
		return
	}
	obj.eventClientConnect(client)
}

func (obj *ServerSocket) accept(listener *net.TCPListener) {
	for {
		tcpconn, err := listener.AcceptTCP()
		if err != nil {
			errorf("accept error is %s\n", err.Error())
			continue
		}

		func() {
			obj.lock()
			defer obj.unlock()
			client := NewClient(tcpconn, obj.eventClientIONew(tcpconn)) //新建客户端对象
			client.SetEventDataBlockNew(obj.eventDataBlockNew)
			obj.onClientConnect(client) //如果设置了eventClientConnect将会在此触发
			client.StartReadThread()
		}()
	}
}
