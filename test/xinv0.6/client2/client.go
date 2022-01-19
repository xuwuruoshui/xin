package main

import (
	"github.com/xuwuruoshui/xin/xnet"
	"log"
	"net"
	"time"
)

func main() {

	// 1.链接远程服务器，得到conn
	conn, err := net.Dial("tcp", ":7777")
	if err != nil {
		log.Println("Client start error:", err)
		return
	}
	// 2.conn调用write
	for {
		datapack := xnet.NewDataPack()

		// 1.发送消息
		msg := xnet.NewMessage(2, []byte("Xin V0.6 client Test message!!!"))
		packData, err := datapack.Pack(msg)
		if err != nil {
			log.Println("pack data error:", err)
			break
		}

		_, err = conn.Write(packData)
		if err != nil {
			log.Println("write data error:", err)
			break
		}

		//2.接受消息
		data := make([]byte, datapack.GetHeadLen())
		if _, err := conn.Read(data); err != nil {
			break
		}
		msgRecv, err := datapack.Unpack(data)
		if err != nil {
			break
		}

		if msgRecv.GetLength() > 0 {
			msgRecv.SetData(make([]byte, msgRecv.GetLength()))
			_, err := conn.Read(msgRecv.GetData())
			if err != nil {
				break
			}
			log.Println(msgRecv.GetId(), msgRecv.GetLength(), string(msgRecv.GetData()))
		}

		time.Sleep(time.Second)
	}
}
