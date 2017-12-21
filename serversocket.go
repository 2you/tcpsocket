package tcpsocket

import (
	"fmt"
	"net"
)

func NewServer() *ServerSocket {
	p := new(ServerSocket)
	return p
}

func (this *ServerSocket) SetPort(port uint) {
	this.port = port
}

func (this *ServerSocket) SetAction(act SocketAction) {
	this.action = act
}

func (this *ServerSocket) Open() error {
	var (
		err         error
		tcpAddr     *net.TCPAddr
		tcpListener *net.TCPListener
	)
	sPort := fmt.Sprintf(":%d", this.port)
	debugln(`socket server listen port is`, this.port)
	if tcpAddr, err = net.ResolveTCPAddr("tcp4", sPort); err != nil { //获取一个tcpAddr
		return err
	}
	if tcpListener, err = net.ListenTCP("tcp", tcpAddr); err != nil {
		return err
	}
	go this.accept(tcpListener)
	return nil
}

func (obj *ServerSocket) lock() {
	obj.mutex.Lock()
}

func (obj *ServerSocket) unlock() {
	obj.mutex.Unlock()
}

func (this *ServerSocket) accept(listener *net.TCPListener) {
	for {
		tcpconn, err := listener.AcceptTCP()
		if err != nil {
			errorf("accept error is %s\n", err.Error())
			continue
		}

		func() {
			this.lock()
			defer this.unlock()
			client := NewClient() //新建客户端对象
			client.accept(tcpconn, this.action)
		}()
	}
}
