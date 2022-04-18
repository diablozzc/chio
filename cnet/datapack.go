package cnet

import (
	"bytes"
	"chio/ciface"
	"chio/utils"
	"encoding/binary"
	"errors"
)

type DataPack struct {

}

// 拆包封包实例初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包头长度
func (dp *DataPack) GetHeadLen() uint32 {
	// Datalen uint32(4) + ID uint32(4)
	return 8
}

// 封包
func (dp *DataPack) Pack(msg ciface.IMessage) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 将dataLen 写入到dataBuff中
	if err:=binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	// 将MsgId 写入到dataBuff中
		if err:=binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	// 将data 写入到dataBuff中
		if err:=binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包 将包的Head信息读取出来 再根据Head信息创建一个Msg
func (dp *DataPack) Unpack(data []byte) (ciface.IMessage, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(data)

	// 只解压head的信息，得到dataLen和msgID
	msg:=&Message{}

	// 读dataLen
	if err:=binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 读msgID
	if err:=binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断dataLen的长度是否超出我们允许的最大包长度
	if (utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize) {
		return nil, errors.New("too large msg data recieved")
	}


	return msg, nil
}