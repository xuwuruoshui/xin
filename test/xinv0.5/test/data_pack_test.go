package test

import (
	"github.com/xuwuruoshui/xin/xnet"
	"log"
	"net"
	"testing"
)

// 测试拆包封包
func TestDataPack(t *testing.T) {

	// 服务端模拟
	listen, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		log.Println("server listen error:", err)
		return
	}

	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Println("server accpect error:", err)
			}
			go func(conn net.Conn) {

				dataPack := xnet.NewDataPack()

				for {
					// 1、从conn中读head
					data := make([]byte, dataPack.GetHeadLen())
					_, err := conn.Read(data)
					if err != nil {
						log.Println("read head error:", err)
						break
					}

					msg, err := dataPack.Unpack(data)
					if err != nil {
						log.Println("server unpack error:", err)
						break
					}
					if msg.GetLength() > 0 {
						data := make([]byte, msg.GetLength())
						_, err := conn.Read(data)
						if err != nil {
							log.Println("read data error:", err)
						}
						msg.SetData(data)
					}
					log.Println(msg.GetId(), msg.GetLength(), string(msg.GetData()))
				}
			}(conn)
		}
	}()

	// 客户端模拟
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		log.Println("client connect error:", err)
		return
	}
	// 创建一个封包对象
	// 模拟粘包过程

	msg := &xnet.Message{Id: 1, Length: 7, Data: []byte{'h', 'e', 'l', 'l', 'o', '!', '!'}}
	pack := xnet.NewDataPack()
	data, err := pack.Pack(msg)

	msg2 := &xnet.Message{Id: 2, Length: 8, Data: []byte{'x', 'i', 'n', 'n', 'i', 'u', 'b', 'i'}}
	pack2 := xnet.NewDataPack()
	data2, err := pack2.Pack(msg2)

	data = append(data, data2...)
	conn.Write(data)
	select {}
}
