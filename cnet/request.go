package cnet

import "chio/ciface"

type Request struct {
	// 已经和客户端建立好的链接
	conn ciface.IConnection

	// 客户端请求的消息
	msg ciface.IMessage
}

// 得到当前链接
func (r *Request) GetConnection() ciface.IConnection {
	return r.conn
}

// 得到请求的消息数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}
