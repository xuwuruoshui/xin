# Xin
<a href="https://github.com/xuwuruoshui/xin/blob/main/LICENSE"><img src="https://img.shields.io/badge/LICENSE-GPL%203.0-blue"/></a>

Xin is a Golang-based TCP server framework

## Install
```go
go get -u github.com/xuwuruoshui/xin
```

## Usage
server
```go
func main(){
    // 1、create xin server
    s := xnet.NewServer()
    
    // 2、register Connection Hook
    s.SetOnConnStart(afterConnection)
    s.SetOnConnStop(beforeStop)
    
    // 3、Add Router(write your logic to router's handle)
    s.AddRouter(1, &PingRouter{})
    
    // 4、run server
    s.Run()
}

// imple BaseRouter
type HelloRouter struct {
    xnet.BaseRouter
}
// write your logic
func (h *HelloRouter) Handle(req xifs.XRequest) {
    log.Println("Call Router Handle")
    // Read Client Data
    log.Printf("Recv from clien:msgId=%d,data=%s\n", req.MsgId(), req.Data())
    // Write Client Data
    if err := req.Connection().SendMsg(req.MsgId(), []byte("XinV1.0 server Hello...")); err != nil {
        log.Println("send msg error:", err)
    }
}
```
client
```go
func main() {

	// 1.conn to server
	conn, err := net.Dial("tcp", ":7777")
	if err != nil {
		log.Println("Client start error:", err)
		return
	}
	// 2. write msg to server
	for {
		datapack := xnet.NewDataPack()

		// 1.发送消息
		msg := xnet.NewMessage(1, []byte("Xin V1.0 client Test message!!!"))
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

		//2. recv message
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
```
[example](https://github.com/xuwuruoshui/xin/tree/main/test/xinv1.0)
open server,client,client2
use `go run xxx.go`

