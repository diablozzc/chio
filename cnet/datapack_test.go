package cnet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {
	// 模拟服务器
	// 创建socket TCP
	listenner, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}

	// 创建一个 go 承载 负责从客户端处理业务的 goroutine

	go func() {
		// 从客户端读数据，拆包处理
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("accept error:", err)
				return
			}
			go func(conn net.Conn) {
				// 处理客户端的请求
				// 拆包的过程
				// 定义一个拆包的对象
				dp := NewDataPack()
				for {
					// 第一次从conn读，把包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error:", err)
						break
					}

					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err:", err)
						return
					}

					if msgHead.GetMsgLen() > 0 {
						// msg是有数据的，需要进行第二次读取
						// 第二次从conn读，根据head中的datalen, 再读取data内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						// 根据datalen再次从io中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err:", err)
							return
						}

						// 完整的一个消息已经读取完毕，可以解析了
						fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.GetMsgLen(), ", data=", string(msg.Data))

					}
				}

			}(conn)
		}
	}()

	// 模拟客户端

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}

	// 准备封包发送的数据
	dp := NewDataPack()
	// 模拟粘包过程，封装两个msg一起发送
	// 封装第一个msg1包
	msg1 := &Message{Id: 1, DataLen: 4, Data: []byte{'c', 'h', 'i', 'o'}}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
		return
	}

	// 封装第二个msg2包
	msg2 := &Message{Id: 2, DataLen: 8, Data: []byte{'c', 'h', 'i', 'o', 't', 'e', 's', 't'}}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 err:", err)
		return
	}

	// 将两个包粘在一起
	sendData1 = append(sendData1, sendData2...)

	// 一次性发送给服务端
	conn.Write(sendData1)

	// 客户端阻塞
	select {}
}
