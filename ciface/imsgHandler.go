package ciface

/*
  消息管理抽象层
*/

type IMsgHandler interface {
	// 调度执行对应的Router消息处理方法
	DoMsgHandler(IRequest)
	// 为消息添加具体的处理逻辑
	AddRouter(uint32, IRouter)
	// 启动一个Worker工作池
	StartWorkerPool()
	// 将消息发送到任务队列
	SendMsgToTaskQueue(IRequest)
}
