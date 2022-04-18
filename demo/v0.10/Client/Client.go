package main

import (
	"chio/cnet"
	"fmt"
	"io"
	"net"
	"time"
)

// 模拟客户端
func main () {
	fmt.Println("client start...")
	time.Sleep(1 * time.Second)
	// 直接连接远程服务器
	conn, err:=net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println("client start err", err)
		return
	}
	// 连接调用write写数据
	for {
		// 发送封包的message消息
		dp := cnet.NewDataPack()
		binaryMsg,err := dp.Pack(cnet.NewMsgPackage(0, []byte("ping chio v1.0")))
		if err != nil {
			fmt.Println("pack msg err:", err)
			return
		}

		if _, err:=conn.Write(binaryMsg); err != nil {
			fmt.Println("write err:", err)
			return
		}
		// 服务器就应该回复一个message

		// 先读取流中的head部分 得到ID和dataLen
		
		binaryHead := make([]byte, dp.GetHeadLen())
		// ReadFull 会把msg填充满为止
		if _, err:=io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error", err)
			break
		}

		// 将二进制的head拆包到msg结构体中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("unpack msgHead error", err)
			break
		}
		
		if msgHead.GetMsgLen() > 0 {
			// 再根据DataLen进行第二次读取，将data读出来

			msg := msgHead.(*cnet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			
			if _, err:=io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error", err)
				return
			}

			fmt.Println("==> Recv Server Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}

		time.Sleep(1 * time.Second)
	}
}