package main

import (
	"log"
	"net"
	"time"
)

func main() {

	// 1.链接远程服务器，得到conn
	conn, err := net.Dial("tcp", ":9999")
	if err != nil {
		log.Println("Client start error:", err)
		return
	}
	// 2.conn调用write
	for {
		_, err := conn.Write([]byte("Hello Xin V0.4"))
		if err != nil {
			log.Println("Write message error:", err)
			return
		}

		data := make([]byte, 512)
		n, err := conn.Read(data)
		if err != nil {
			log.Println("Read message error:", err)
			return
		}
		log.Println(string(data[:n]))
		time.Sleep(time.Second)
	}
}
