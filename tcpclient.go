package tcpsocket

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func NewTcpClient(conn net.Conn, io IClientIO) *TTcpClient {
	p := new(TTcpClient)
	p.connectTime = time.Now()
	p.readThreadActive = false
	p.socket = conn
	p.io = io
	return p
}

func ClientTo(io IClientIO, addr string, port int) *TTcpClient {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		panic(err)
	}
	return NewTcpClient(conn, io)
}

func (obj *TTcpClient) SetEventDataBlockNew(event DataBlockNewEvent) {
	obj.eventDataBlockNew = event
}

func (obj *TTcpClient) StartReadThread() {
	if !obj.readThreadActive {
		obj.readThreadActive = true
		go obj.handle()
	}
}

func (obj *TTcpClient) RemoteIP() string {
	vAddr := strings.Split(obj.RemoteAddr(), `:`)
	return vAddr[0]
}

func (obj *TTcpClient) RemotePort() int {
	vAddr := strings.Split(obj.RemoteAddr(), `:`)
	if port, err := strconv.Atoi(vAddr[1]); err == nil {
		return port
	}
	return 0
}

func (obj *TTcpClient) RemoteAddr() string {
	return obj.socket.RemoteAddr().String()
}

func (obj *TTcpClient) Close() error {
	return obj.socket.Close()
}

func (obj *TTcpClient) onConnect() {
	obj.io.OnConnect(obj)
}

func (obj *TTcpClient) onClose(err error) {
	obj.io.OnClose(obj, err)
}

func (obj *TTcpClient) onRecv(data IDataBlock) {
	defer func() {

	}()
	obj.io.OnRecv(obj, data)
}

func (obj *TTcpClient) handle() {
	var (
		err   error = nil
		vRead IDataBlock
	)
	defer obj.onClose(err)
	go obj.onConnect()
	for {
		if vRead, err = obj.recvData(); err != nil {
			Errorln(err.Error())
			break
		}
		Debugf("<%s>数据包体长度为%d\n", obj.RemoteAddr(), vRead.BodySize())
		go obj.onRecv(vRead)
	}
}

func (obj *TTcpClient) recvData() (data IDataBlock, err error) {
	data = obj.eventDataBlockNew()
	if err = obj.recvHead(data); err != nil {
		return nil, err
	}

	if err = obj.recvBody(data); err != nil {
		return nil, err
	}
	return data, nil
}

func (obj *TTcpClient) recvBody(data IDataBlock) error {
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

func (obj *TTcpClient) recvHead(data IDataBlock) error {
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

func (obj *TTcpClient) recvBuf(count int) ([]byte, error) {
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
