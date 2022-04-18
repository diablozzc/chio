package cnet

import (
	// "chio/ciface"
	// "chio/utils"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/diablozzc/chio/ciface"
	"github.com/diablozzc/chio/utils"
)

type Connection struct {
	// 当前Connection隶属于哪个Server
	TcpServer ciface.IServer
	// 当前连接的socket TCPConn
	Conn *net.TCPConn

	// 当前连接的ID
	ConnID uint32

	// 当前连接的状态
	isClose bool

	// 告知当前链接已经退出的/停止 channel (由Reader告知Writer退出)
	ExitChan chan bool

	// 无缓冲管道 用于读、写Goroutine之间的消息通信
	msgChan chan []byte

	// 消息的管理模块
	MsgHandler ciface.IMsgHandler

	// 连接属性集合
	property map[string]interface{}
	// 保护连接属性的锁
	propertyLock sync.RWMutex
}

func NewConnection(server ciface.IServer, conn *net.TCPConn, connID uint32, msgHandler ciface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		isClose:    false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		MsgHandler: msgHandler,
		property:   make(map[string]interface{}),
	}

	// 将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

// 连接的读业务
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println("connID = ", c.ConnID, "[Reader is exit] , remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到buf中
		// buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		// _, err:= c.Conn.Read(buf)
		// if err != nil {
		// 	fmt.Println("recv buf err ", err)
		// 	continue
		// }

		// 创建一个拆包解包对象
		dp := NewDataPack()
		// 读取客户端的Msg head 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head err ", err)
			break
		}

		// 拆包 得到msgID和 MsgDataLen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack err ", err)
			break
		}

		// 根据dataLen，再次读取，放在MsgData中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data err ", err)
				break
			}
		}

		msg.SetData(data)

		req := Request{
			conn: c,
			msg:  msg,
		}

		// 将消息交给路由处理
		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}

}

// 写消息Goroutine， 专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println("ConnID = ", c.ConnID, "[Writer is exit], remote addr is ", c.RemoteAddr().String())

	// 不断的阻塞的等待channel的消息， 如果有数据，就发送
	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data err ", err)
				return
			}
		case <-c.ExitChan:
			// 代表Reader已经退出，并且不再需要Write
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start().. ConnID = ", c.ConnID)
	// 启动从当前连接中读取数据流的业务
	go c.StartReader()
	// 启动从当前连接中写数据流的业务
	go c.StartWriter()

	c.TcpServer.CallAfterConnStart(c)

}

// 停止连接， 将连接从ConnMgr中删除
func (c *Connection) Stop() {
	fmt.Println("Conn Stop().. ConnID = ", c.ConnID)

	if c.isClose {
		return
	}

	c.isClose = true

	// 调用开发者注册 销毁连接之前 需要执行的业务Hook
	c.TcpServer.CallBeforeConnStop(c)

	// 关闭socket连接
	c.Conn.Close()
	// 告知Writer关闭
	c.ExitChan <- true

	// 将当前连接从ConnMgr中删除
	c.TcpServer.GetConnMgr().Remove(c)

	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 提供一个SendMsg的方法,将message数据发送给客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClose {
		return errors.New("Connection closed when send msg")
	}

	// 将data进行封包 MsgDataLen + MsgID + MsgData
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg id = " + string(rune(msgId)))
	}

	// 将数据发送给channel
	c.msgChan <- binaryMsg

	return nil
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
