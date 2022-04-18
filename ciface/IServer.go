package ciface

// 服务器接口
type IServer interface {
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 运行服务器
	Serve()

	// 路由器(路由调度)添加路由
	AddRouter(uint32, IRouter)

	// 获取当前server 的连接管理器
	GetConnMgr() IConnManager

	// 注册AfterConnStart钩子函数
	SetAfterConnStart(func(IConnection))
	// 
	SetBeforeConnStop(func(IConnection))

	// 调用AfterFunc钩子函数
	CallAfterConnStart(IConnection)
	// 调用BeforeFunc钩子函数
	CallBeforeConnStop(IConnection)


}
