package main

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/diablozzc/chio/cnet"
	"github.com/diablozzc/chio/demo/v0.10/Client/chioClient"
)

var (
	done      = make(chan bool, 1)
	interrupt = make(chan os.Signal, 1)
	reconnect = make(chan bool, 1)
)

// 模拟客户端
func main() {
	reconnect <- true

	// 监听循环
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt:", interrupt)
			done <- true
		case <-reconnect:
			doConnect()
		}
	}
}

func doConnect() error {
	bs := backoff.NewExponentialBackOff()
	bs.MaxElapsedTime = 0
	bs.MaxInterval = time.Second * 20
	bs.Multiplier = 1.5
	bs.InitialInterval = time.Second * 2
	err := backoff.Retry(bootstrap, bs)
	if err != nil {
		return err
	}
	return nil
}

func bootstrap() error {
	// 创建连接
	log.Println("连接xdr server...")
	cc, err := new(chioClient.ChioClient).GetInstance()
	if err != nil {
		return errors.New("创建连接失败")
	}

	// 连接成功后开始消息监听
	go startMessageHandle(cc)
	// 定时发送ping
	go func() {
		for {
			err := cc.Send(3, "ping server")
			if err != nil {
				return
			}
			time.Sleep(time.Second * 3)
		}
	}()

	return nil
}

// 开始消息监听
func startMessageHandle(cc *chioClient.ChioClient) {
	for {
		message, err := cc.Recv()
		if err != nil {
			log.Println("recv message failed:", err)
			continue
		}
		if message == nil {
			log.Println("recv message is nil")
			cc.Close()
			reconnect <- true
			return
		}
		route(message)
	}
}

// 消息路由
func route(msg *cnet.Message) {
	msgType := msg.GetMsgID()

	switch msgType {
	case 4:
		log.Println("收到服务器消息:", string(msg.GetData()))
	default:
		log.Println("接收到未知消息:", string(msg.GetData()))
	}
}
