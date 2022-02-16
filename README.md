# Xin
<hr>
[![](https://img.shields.io/badge/LICENSE-GPL%203.0-blue "title")](https://github.com/xuwuruoshui/xin/blob/main/LICENSE)
Xin is a Golang-based TCP server framework

## USAGE
**Install**
```go
go get -u github.com/xuwuruoshui/xin
```


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


