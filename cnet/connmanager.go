package cnet

import (
	"chio/ciface"
	"fmt"
	"sync"
)

/*
  连接管理模块
*/

type ConnManager struct {
	// 管理的连接集合
	connections map[uint32] ciface.IConnection
	// 保护连接集合的读写锁
	connLock sync.RWMutex

}

// 创建ConnManager的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32] ciface.IConnection),
	}
}

	// 添加连接
	func (cm *ConnManager) Add(conn ciface.IConnection) {
		// 保护共享资源map, 加写锁
		cm.connLock.Lock()
		defer cm.connLock.Unlock()

		// 将conn加入到connections中
		cm.connections[conn.GetConnID()] = conn
		fmt.Println("connID = ", conn.GetConnID(), "add to ConnManager successfully: conn num = ", cm.Len())
	}
	// 删除连接
	func (cm *ConnManager) Remove(conn ciface.IConnection) {
		// 保护共享资源map, 加写锁
		cm.connLock.Lock()
		defer cm.connLock.Unlock()

		// 删除连接信息
		delete(cm.connections, conn.GetConnID())
		fmt.Println("connID = ", conn.GetConnID(), "remove from ConnManager successfully: conn num = ", cm.Len())

	}
  // 根据connID获取连接
	func (cm *ConnManager) Get(connID uint32) (ciface.IConnection, error) {
		// 保护共享资源map, 加读锁
		cm.connLock.RLock()
		defer cm.connLock.RUnlock()

		if conn, ok := cm.connections[connID]; ok {
			return conn, nil
		} else {
			return nil, fmt.Errorf("connID = %d not found", connID)
		}

	}
	// 得到当前连接总数
	func (cm *ConnManager) Len() int {
		return len(cm.connections)
	}
	// 清除并终止所有连接
	func (cm *ConnManager) ClearConn() {
		cm.connLock.Lock()
		defer cm.connLock.Unlock()
		// 删除conn并停止conn的工作

		for connID, conn := range cm.connections {
			// 停止
			conn.Stop()
			// 删除
			delete(cm.connections, connID)
		}

		fmt.Println("Clear all connections successfully: conn num = ", cm.Len())
	}

