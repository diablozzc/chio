package ciface

import "net"

// 定义连接模块的接口

type IConnection interface {
	// 启动连接 让当前的连接准备开始工作
	Start()
	// 停止连接 结束当前连接的工作
	Stop()
	// 获取当前连接的绑定 socket conn
	GetTCPConnection() *net.TCPConn

	// 获取当前连接模块的连接ID
	GetConnID() uint32
	// 获取远程客户端的 TCP 状态 IP 端口等信息
	RemoteAddr() net.Addr

	// 发送数据给远程的客户端
	SendMsg(uint32, []byte) error

	// 设置连接属性
	SetProperty(key string, value interface{})
	// 获取连接属性
	GetProperty(key string) (interface{}, error)
	// 移除连接属性
	RemoveProperty(key string)
}

type HandleFunc func(*net.TCPConn, []byte, int) error
