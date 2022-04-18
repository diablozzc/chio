package main

import (
	"chio/ciface"
	"chio/cnet"
	"fmt"
)

// ping test 自定义路由

type PingRouter struct {
	cnet.BaseRouter
}

// Test Handle
func (r *PingRouter) Handle(request ciface.IRequest) {
	fmt.Println("call router Handle...")

	// 先读取客户端的数据，再回写ping...
	fmt.Println("recv from client : msgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(200, []byte("ping...\n"))
	if err != nil {
		fmt.Println("call back ping... err ", err)
	}

}

// ---------------------------------------------------

type HelloRouter struct {
	cnet.BaseRouter
}

func (r *HelloRouter) Handle(request ciface.IRequest) {
	fmt.Println("call router Handle...")

	// 先读取客户端的数据，再回写ping...
	fmt.Println("recv from client : msgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(201, []byte("hello...\n"))
	if err != nil {
		fmt.Println("call back hello... err ", err)
	}

}

// ---------------------------------------------------

// 创建连接之后执行钩子函数
func DoConnectionBegin(conn ciface.IConnection) {
	fmt.Println(" ==> DoConnectionBegin...")
	if err := conn.SendMsg(202, []byte("DoConnection BEGIN...\n")); err!=nil{
		fmt.Println("DoConnection BEGIN... err ", err)
	}

	fmt.Println("Set conn Name")
	conn.SetProperty("Name", "黑超特攻")
	conn.SetProperty("Home", "zhilogos.com")
	conn.SetProperty("Github", "github.com/zhilogos")
}

// 连接断开之前的需要执行的函数
func DoConnectionLost(conn ciface.IConnection) {
	fmt.Println(" ==> DoConnectionLost...")
	fmt.Println("conn ID = ", conn.GetConnID())

	// 获取连接属性

	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name = ", name)
	}
}




func main() {
	// 创建一个server句柄
	s := cnet.NewServer("chio v1.0")
	// 注册自定义路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	// 注册连接回调，当有新的连接进入时，会调用该方法
	s.SetAfterConnStart(DoConnectionBegin)
	s.SetBeforeConnStop(DoConnectionLost)

	s.Serve()
}