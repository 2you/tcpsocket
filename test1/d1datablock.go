package main

import (
	"encoding/binary"

	"git.oschina.net/tigercat/tcfunc"
	"github.com/2you/tcpsocket"
)

type D1DataBlock struct {
	tcpsocket.IDataBlock
	headContent tcpsocket.DataBlockHead
	bodyContent tcpsocket.DataBlockBody
}

func NewD1DataBlock() tcpsocket.IDataBlock {
	var p = &D1DataBlock{}
	return p
}

func (obj *D1DataBlock) HeadSize() int {
	return 4
}

func (obj *D1DataBlock) BodySize() int {
	var nRet = binary.LittleEndian.Uint32(obj.headContent[0:4])
	return int(nRet)
}

//获取数据包头
func (obj *D1DataBlock) HeadContent() tcpsocket.DataBlockHead {
	return obj.headContent
}

//设置数据包头
func (obj *D1DataBlock) SetHeadContent(v tcpsocket.DataBlockHead) {
	obj.headContent = v
}

//获取数据包包体
func (obj *D1DataBlock) BodyContent() tcpsocket.DataBlockBody {
	return obj.bodyContent
}

//设置数据包包体
func (obj *D1DataBlock) SetBodyContent(v tcpsocket.DataBlockBody) {
	obj.bodyContent = v
	nLen := len(v)
	obj.SetBodyLength(nLen)
}

//获取命令id
func (obj *D1DataBlock) GetCmdId() int {
	var nRet = binary.LittleEndian.Uint32(obj.headContent[30:34])
	return int(nRet)
}

//设置命令id
func (obj *D1DataBlock) SetCmdId(v int) {
	binary.LittleEndian.PutUint32(obj.headContent[30:34], uint32(v))
}

//获取unixtime
func (obj *D1DataBlock) GetUnixTime() int64 {
	var nRet = int64(binary.LittleEndian.Uint64(obj.headContent[34:42]))
	return nRet
}

//设置unixtime
func (obj *D1DataBlock) SetUnixTime(value int64) {
	binary.LittleEndian.PutUint64(obj.headContent[34:42], uint64(value))
}

//设置sequnce
func (obj *D1DataBlock) GetSequnce() uint64 {
	var nRet = binary.LittleEndian.Uint64(obj.headContent[42:50])
	return nRet
}

//获取sequnce
func (obj *D1DataBlock) SetSequnce(value uint64) {
	binary.LittleEndian.PutUint64(obj.headContent[42:50], value)
}

//获取包体长度
func (obj *D1DataBlock) GetBodyLength() int {
	var nRet = binary.LittleEndian.Uint32(obj.headContent[100:104])
	return int(nRet)
}

//设置包体长度
func (obj *D1DataBlock) SetBodyLength(v int) {
	//	binary.LittleEndian.PutUint32(obj.headContent[100:104], uint32(v))
	binary.LittleEndian.PutUint32(obj.headContent[0:4], uint32(v))
}

//设置checksum
func (obj *D1DataBlock) SetChkeckSum() int {
	var nRet = binary.LittleEndian.Uint32(obj.headContent[120:124])
	return int(nRet)
}

//获取checksum
func (obj *D1DataBlock) SetCheckSum(v int) {
	binary.LittleEndian.PutUint32(obj.headContent[120:124], uint32(v))
}

func (obj *D1DataBlock) DecodeBase64Body2String() string {
	return string(obj.DecodeBase64Body2Bytes())
}

func (obj *D1DataBlock) DecodeBase64Body2Bytes() []byte {
	return tcfunc.Base64Decode(string(obj.bodyContent))
}
