package cnet

import (
	"chio/ciface"
	"chio/utils"
	"fmt"
	"net"
)

// iServer的接口实现, 定义一个Server的服务器模块
type Server struct {
	// 服务器名称
	Name string
	// 服务器绑定的ip版本
	IPVersion string
	// 服务器监听的IP
	IP string
	// 服务器监听的端口
	Port int
	// 当前server的消息管理模块，用来绑定MsgID和对应的处理业务的Handle实例
	MsgHandler ciface.IMsgHandler
	// 当前server的链接管理器
	ConnMgr ciface.IConnManager

	// AfterConnStart 是一个钩子函数，当Conn执行完初始化方法之后，会调用这个方法
	AfterConnStart func(ciface.IConnection)

	// BeforeConnStop 是一个钩子函数，当Conn执行完关闭方法之后，会调用这个方法
	BeforeConnStop func(ciface.IConnection)
}

func (s *Server) Start() {
	fmt.Printf(
		"[Chio] Server Name: %s listenner at IP: %s, Port: %d is starting\n",
		utils.GlobalObject.Name,
		utils.GlobalObject.Host,
		utils.GlobalObject.TcpPort,
	)
	fmt.Printf(
		"[Chio] Version %s, MaxConn: %d, MaxPackageSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize,
	)

	go func() {
		// 0. 启动worker工作池机制
		s.MsgHandler.StartWorkerPool()

		// 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		// 监听服务器地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}

		fmt.Println("start server succ, ", s.Name, "succ, now listenning...")
		var cid uint32 = 0

		// 阻塞等待客户端链接，处理客户端业务
		for {
			// 如果有客户端链接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 设置最大链接数，如果超过最大链接，那么则关闭此新的连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				// TODO 给客户端回复一个超出链接数的信息
				fmt.Println("[Too many connections] , MaxConn = ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			// 将处理新连接的业务方法和 conn 进行绑定，得到我们的业务方法
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			// 启动当前链接的处理业务
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	// TODO 将一些服务器的资源、状态、工作结束
	fmt.Println("[Chio] Server ", s.Name, "is Stopping ...")
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	// TODO 做一些启动服务器之后的额外业务

	// 阻塞状态
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ciface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Succ!")
}

func (s *Server) GetConnMgr() ciface.IConnManager {
	return s.ConnMgr
}

func NewServer(name string) ciface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandler(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

// 注册AfterConnStart钩子函数
func (s *Server) SetAfterConnStart(hookFunc func(ciface.IConnection)) {
	s.AfterConnStart = hookFunc
}

// 注册BeforeConnStop钩子函数
func (s *Server) SetBeforeConnStop(hookFunc func(ciface.IConnection)) {
	s.BeforeConnStop = hookFunc
}

// 调用AfterConnStart钩子函数
func (s *Server) CallAfterConnStart(conn ciface.IConnection) {
	if nil != s.AfterConnStart {
		fmt.Println("---> CallAfterConnStart...")
		s.AfterConnStart(conn)
	}
}

// 调用BeforeConnStop钩子函数
func (s *Server) CallBeforeConnStop(conn ciface.IConnection) {
	if nil != s.BeforeConnStop {
		fmt.Println("---> CallBeforeConnStop...")
		s.BeforeConnStop(conn)
	}
}
