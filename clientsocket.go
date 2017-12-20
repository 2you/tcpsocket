package tcpsocket

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func NewClient() *ClientSocket {
	p := new(ClientSocket)
	p.active = false
	p.svrclt = false
	p.readThreadActive = false
	return p
}

//func ClientTo(io IClientIO, addr string, port int) *ClientSocket {
//	tcpaddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", addr, port))
//	if err != nil {
//		panic(err)
//	}
//	tcpconn, err := net.DialTCP("tcp", nil, tcpaddr)
//	if err != nil {
//		panic(err)
//	}
//	return NewClient(tcpconn, io)
//}

func (this *ClientSocket) SetAction(act SocketAction) {
	this.action = act
}

func (this *ClientSocket) SetHost(host string) {
	this.host = host
}

func (this *ClientSocket) SetPort(port uint) {
	this.port = port
}

func (this *ClientSocket) accept(tcpconn *net.TCPConn, act SocketAction) error {
	this.svrclt = true
	return this.open(tcpconn, act)
}

func (this *ClientSocket) open(tcpconn *net.TCPConn, act SocketAction) error {
	this.action = act
	this.socket = tcpconn
	this.active = true
	this.connectTime = time.Now()
	this.startReadThread()
	return nil
}

func (this *ClientSocket) Open() error {
	if this.svrclt {
		return fmt.Errorf("server accept client can not open")
	}

	tcpaddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", this.host, this.port))
	if err != nil {
		return err
	}
	tcpconn, err := net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		return err
	}
	return this.open(tcpconn, this.action)
}

func (this *ClientSocket) Close() error {

	return this.socket.Close()
}

func (obj *ClientSocket) SetEventDataBlockNew(event DataBlockNewEvent) {
	obj.eventDataBlockNew = event
}

func (this *ClientSocket) startReadThread() {
	go this.handleRead()
}

func (obj *ClientSocket) RemoteIP() string {
	vAddr := strings.Split(obj.RemoteAddr(), `:`)
	return vAddr[0]
}

func (obj *ClientSocket) RemotePort() int {
	vAddr := strings.Split(obj.RemoteAddr(), `:`)
	if port, err := strconv.Atoi(vAddr[1]); err == nil {
		return port
	}
	return 0
}

func (obj *ClientSocket) RemoteAddr() string {
	return obj.socket.RemoteAddr().String()
}

func (obj *ClientSocket) Close() error {
	return obj.socket.Close()
}

func (obj *ClientSocket) WriteA(v []byte, size int) error {
	return obj.writeA(v, size)
}

func (obj *ClientSocket) WriteB(v []byte) error {
	return obj.writeB(v)
}

func (obj *ClientSocket) writeA(v []byte, size int) error {
	if v == nil {
		return fmt.Errorf("发送数据不能为空")
	}

	if size < 1 {
		return fmt.Errorf("发送数据长度必须大于0")
	}
	bufsize := len(v)
	if bufsize < size {
		return fmt.Errorf("待发数据长度[%d]小于应发数据长度[%d]", bufsize, size)
	}
	wbuf := v[:size]
	n, err := obj.socket.Write(wbuf)
	if err != nil {
		return err
	}

	if n != size {
		return fmt.Errorf("应发送%d个字节 实际发送%d个字节", size, n)
	}
	return nil
}

func (obj *ClientSocket) writeB(v []byte) error {
	return obj.writeA(v, len(v))
}

func (obj *ClientSocket) handleOnConnect() {
	obj.io.OnConnect(obj)
}

func (obj *ClientSocket) onClose(err error) {
	obj.io.OnClose(obj, err)
}

func (obj *ClientSocket) onRecv(data IDataBlock) {
	defer func() {

	}()
	obj.io.OnRecv(obj, data)
}

func (obj *ClientSocket) handleRead() {
	var (
		err   error = nil
		vRead IDataBlock
	)
	defer obj.onClose(err)
	go obj.handleOnConnect()
	for {
		if vRead, err = obj.recvData(); err != nil {
			errorln(err.Error())
			break
		}
		debugf("<%s>数据包体长度为%d\n", obj.RemoteAddr(), vRead.BodySize())
		go obj.onRecv(vRead)
	}
}

func (obj *ClientSocket) recvData() (data IDataBlock, err error) {
	data = obj.eventDataBlockNew()
	if err = obj.recvHead(data); err != nil {
		return nil, err
	}

	if err = obj.recvBody(data); err != nil {
		return nil, err
	}
	return data, nil
}

func (obj *ClientSocket) recvBody(data IDataBlock) error {
	var (
		vBuf []byte
		vErr error
	)
	nBodySize := data.BodySize()
	if nBodySize < 1 {
		return nil
	}
	if vBuf, vErr = obj.recvBuf(nBodySize); vErr != nil {
		return vErr
	}
	data.SetBodyContent(vBuf)
	return nil
}

func (obj *ClientSocket) recvHead(data IDataBlock) error {
	var (
		vBuf []byte
		vErr error
	)
	nHeadSize := data.HeadSize()
	if vBuf, vErr = obj.recvBuf(nHeadSize); vErr != nil {
		return vErr
	}
	data.SetHeadContent(vBuf)
	return nil
}

func (obj *ClientSocket) recvBuf(count int) ([]byte, error) {
	var (
		vRet               []byte
		nAllRead, nCurRead int
		vErr               error
	)
	if count < 1 {
		return nil, errors.New(`read count less than 1`)
	}
	vRet = make([]byte, count)
	nAllRead = 0
	for nAllRead < count {
		if nCurRead, vErr = obj.socket.Read(vRet[nAllRead:count]); vErr != nil {
			return nil, vErr
		}
		nAllRead += nCurRead
	}
	return vRet, nil
}
