package utils

import (
	"chio/ciface"
	"encoding/json"
	"io/ioutil"
)

/*
	存储一切有关 chio框架的全局参数，供其他模块使用
	一些参数是可以通过zinx.json由用户进行配置
*/

type GlobalObj struct {
	// Server

	TcpServer ciface.IServer // 存储一个Server句柄
	Host      string         // 当前服务器主机监听的IP
	TcpPort   int            // 当前服务器主机监听的端口号
	Name      string         // 当前服务器的名称

	// Chio

	Version          string // chio版本号
	MaxConn          int    // 当前服务器主机允许的最大链接数
	MaxPackageSize   uint32 // 当前Chio框架数据包的最大值
	WorkerPoolSize   uint32 // 当前Chio框架业务工作Worker池的worker数量
	MaxWorkerTaskLen uint32 // 当前Chio框架业务工作Worker对应的任务队列最大值
}

/*
  定义一个全局的对外GlobalObj
*/

var GlobalObject *GlobalObj

/*
	从 chio.json加载用户配置
*/
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/chio.json")
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}

}

// 定义一个初始化方法
func init() {
	// 默认配置
	GlobalObject = &GlobalObj{
		Name:             "ChioServerApp",
		Version:          "V0.9",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	// 应该尝试从 conf/chio.json中加载用户配置
	GlobalObject.Reload()
}
