package cnet

import (
	"chio/ciface"
	"chio/utils"
	"fmt"
	"strconv"
)

/*
  消息处理模块的实现
*/

type MsgHandler struct {
	// 消息的id和router对应关系
	Apis map[uint32]ciface.IRouter

	// 负责Worker取任务的消息队列
	WorkerPoolSize uint32
	// 业务工作Worker池的worker数量
	TaskQueue []chan ciface.IRequest
}

// 初始化/创建消息处理模块
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ciface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ciface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// 根据msgId调度对应的router处理业务
func (mh *MsgHandler) DoMsgHandler(request ciface.IRequest) {
	// 从request中获取msgId
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgID(), " is not FOUND!")
		return
	}
	// 根据msgId调用相应的router对应的方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandler) AddRouter(msgId uint32, router ciface.IRouter) {
	// 判断当前msgId是否已经注册过
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api, msgId = " + string(rune(msgId)))
	}
	// 添加msgId和router的对应关系
	mh.Apis[msgId] = router
	fmt.Println("Add api msgId = " + strconv.Itoa(int(msgId)))
}

// 启动一个Worker工作池(开启工作池的动作只能发生一次，框架只有一个worker工作池)
func (mh *MsgHandler) StartWorkerPool() {
	// 根据workerPoolSize 分别开启Worker，每个Worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 当前的worker对应的channel消息队列 开辟空间
		mh.TaskQueue[i] = make(chan ciface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 启动当前的Worker，阻塞等待消息从channel传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (mh *MsgHandler) StartOneWorker(workerID int, taskQueue chan ciface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started.")
	// 不断的阻塞等待对应消息队列的消息

	// for {
	// 	select {
	// 		// 如果有消息过来，出列的就是一个客户端的Request，执行当前Request所绑定的业务
	// 	case request := <- taskQueue:
	// 		mh.DoMsgHandler(request)
	// 	}
	// }

	for request := range taskQueue {
		mh.DoMsgHandler(request)
	}

}

// 将消息交给TaskQueue,由worker进行处理
func (mh *MsgHandler) SendMsgToTaskQueue(request ciface.IRequest) {
	// 将消息平均分配给不同的worker
	// 根据客户端建立的链接ID来分配给不同的worker处理
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(), " request to workerID = ", workerID)
	// 将消息发送给任务队列
	mh.TaskQueue[workerID] <- request

}
