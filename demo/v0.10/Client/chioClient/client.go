package chioClient

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/diablozzc/chio/cnet"
)

var (
	Host = "127.0.0.1"
	Port = "9999"
)

type ChioClient struct {
	conn net.Conn
}

var instance *ChioClient

func (c *ChioClient) GetInstance() (*ChioClient, error) {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	if instance == nil {
		p, err := strconv.Atoi(Port)
		if err != nil {
			return nil, err
		}

		conn, err := connect(Host, p)
		if err != nil {
			log.Printf("连接服务器失败, Host: %s, Port: %d", Host, p)
			return nil, err
		}
		instance = &ChioClient{conn: conn}
	} else {
		log.Printf("已经连接服务器, Host: %s, Port: %s", Host, Port)
	}

	return instance, nil

}

// 连接chio server
func connect(host string, port int) (net.Conn, error) {
	return net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
}

// 关闭连接
func (c *ChioClient) Close() {
	c.conn.Close()
	instance = nil
}

// 发送消息
func (c *ChioClient) Send(msgId uint32, msg string) error {
	dp := cnet.NewDataPack()
	binaryMsg, err := dp.Pack(cnet.NewMsgPackage(msgId, []byte(msg)))
	if err != nil {
		fmt.Println("pack msg failed:", err)
		return err
	}

	if _, err := c.conn.Write(binaryMsg); err != nil {
		fmt.Println("send message to chio server failed:", err)
		return err
	}
	return nil
}

// 接收消息
func (c *ChioClient) Recv() (*cnet.Message, error) {
	dp := cnet.NewDataPack()
	binaryHead := make([]byte, dp.GetHeadLen())
	if _, err := io.ReadFull(c.conn, binaryHead); err != nil {
		if err == io.EOF {
			fmt.Println("no data")
			return nil, nil
		}
		return nil, err
	}

	msgHead, err := dp.Unpack(binaryHead)
	if err != nil {
		return nil, err
	}

	if msgHead.GetMsgLen() > 0 {
		msg := msgHead.(*cnet.Message)
		msg.Data = make([]byte, msg.GetMsgLen())
		if _, err := io.ReadFull(c.conn, msg.Data); err != nil {
			fmt.Println("read data failed:", err)
			return nil, err
		}
		return msg, nil
	}
	return nil, nil
}
