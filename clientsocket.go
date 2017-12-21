package tcpsocket

import (
	"bytes"
	"encoding/binary"
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

func ClientTo(host string, port uint, action SocketAction) (client *ClientSocket, err error) {
	client = NewClient()
	client.SetHost(host)
	client.SetPort(port)
	client.SetAction(action)
	if err = client.Open(); err != nil {
		return nil, err
	}
	return client, nil
}

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

func (this *ClientSocket) close() error {
	if err := this.socket.Close(); err != nil {
		return err
	}
	this.stopReadThread()
	return nil
}

func (this *ClientSocket) Close() error {
	return this.close()
}

func (this *ClientSocket) startReadThread() {
	if this.readThreadActive {
		return
	}
	this.readThreadActive = true
	go this.threadHandleRead()
}

func (this *ClientSocket) stopReadThread() {
	this.readThreadActive = false
}

func (this *ClientSocket) readBuf(count uint64) (buf []byte, err error) {
	var (
		allSize  uint64
		currSize int
	)
	if count < 1 {
		return nil, nil
	}
	buf = make([]byte, count)
	allSize = 0
	for allSize < count {
		if currSize, err = this.socket.Read(buf[allSize:count]); err != nil {
			return nil, err
		}
		allSize += uint64(currSize)
	}
	return buf, nil
}

func (this *ClientSocket) readHead() (buf []byte, err error) {
	headSize := uint64(this.action.GetHeadSize())
	return this.readBuf(headSize)
}

func (this *ClientSocket) getBodySize(headBuf []byte) (size uint64, err error) {
	headSize := this.action.GetHeadSize()
	if headSize != len(headBuf) {
		return 0, fmt.Errorf("head buf size if error")
	}
	lob := this.action.LittleOrBig()
	iBit := this.action.GetBodySizeLength()
	offSet := this.action.GetBodySizeOffSet()
	sizeBuf := headBuf[offSet:]
	var bodySize uint64
	if lob == 'L' {
		switch iBit {
		case 2:
			bodySize = uint64(binary.LittleEndian.Uint16(sizeBuf))
		case 4:
			bodySize = uint64(binary.LittleEndian.Uint32(sizeBuf))
		case 8:
			bodySize = uint64(binary.LittleEndian.Uint64(sizeBuf))
		default:
			return 0, fmt.Errorf("body size bit error")
		}

	} else if lob == 'B' {
		switch iBit {
		case 2:
			bodySize = uint64(binary.BigEndian.Uint16(sizeBuf))
		case 4:
			bodySize = uint64(binary.BigEndian.Uint32(sizeBuf))
		case 8:
			bodySize = uint64(binary.BigEndian.Uint64(sizeBuf))
		default:
			return 0, fmt.Errorf("body size bit error")
		}
	} else {
		return 0, errors.New("body size parse method error")
	}
	return bodySize, nil
}

func (this *ClientSocket) readBody(headBuf []byte) (buf []byte, err error) {
	bodySize, err := this.getBodySize(headBuf)
	if err != nil {
		return nil, err
	}
	return this.readBuf(bodySize)
}

func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

func (this *ClientSocket) readData() (buf []byte, err error) {
	if buf, err = this.readHead(); err != nil {
		return nil, err
	}
	headBuf := buf
	if buf, err = this.readBody(headBuf); err != nil {
		return nil, err
	}

	if buf == nil {
		return headBuf, nil
	} else {
		return func(v ...[]byte) []byte {
			return bytes.Join(v, []byte(``))
		}(headBuf, buf), nil
	}
}

func (this *ClientSocket) handleOnConnect() {
	go this.action.OnConnect(this)
}

func (this *ClientSocket) handleOnDisconnect(err error) {
	this.Close()
	this.action.OnDisconnect(this, err)
}

func (this *ClientSocket) handleOnRead(data []byte) {
	go this.action.OnRead(this, data)
}

func (this *ClientSocket) threadHandleRead() {
	var (
		err  error = nil
		data []byte
	)
	defer this.handleOnDisconnect(err)
	this.handleOnConnect()
	for {
		if data, err = this.readData(); err != nil {
			break
		}
		this.handleOnRead(data)
	}
}

func (this *ClientSocket) RemoteIP() string {
	arrAddr := strings.Split(this.RemoteAddr(), `:`)
	return arrAddr[0]
}

func (this *ClientSocket) RemotePort() int {
	arrAddr := strings.Split(this.RemoteAddr(), `:`)
	if port, err := strconv.Atoi(arrAddr[1]); err == nil {
		return port
	}
	return 0
}

func (this *ClientSocket) RemoteAddr() string {
	return this.socket.RemoteAddr().String()
}

func (this *ClientSocket) Write(buf []byte, size int) error {
	return this.write(buf, size)
}

func (this *ClientSocket) WriteBuf(buf []byte) error {
	return this.write(buf, len(buf))
}

func (this *ClientSocket) write(buf []byte, size int) error {
	if buf == nil {
		return fmt.Errorf("write buf can not null")
	}

	if size < 1 {
		return fmt.Errorf("write buf size can not less than 1")
	}
	bufSize := len(buf)
	if bufSize < size {
		return fmt.Errorf("pending write size %d less than actual size %d", bufSize, size)
	}
	wBuf := buf[:size]
	wSize, err := this.socket.Write(wBuf)
	if err != nil {
		return err
	}

	if wSize != size {
		return fmt.Errorf("pending write %d bytes, actual write %d bytes", size, wSize)
	}
	return nil
}
