# xin tcp框架设计

- 基础server
    - 方法
        - 启动服务器(Start)
            - 基本的服务器开发:1.创建addr 2.创建listen 3.处理客户端的基本业务(回显)
        - 停止服务器(Stop)
        - 运行服务器(Run)
            - 调用Start(),然后阻塞处理,之后可以做一些扩展功能
        - 初始化server(NewServer)
    - 属性
        - 名称(Name)
        - 监听ip(IP)
        - 监听端口(Port)
        - IP版本(IPVersion)

## 简单的链接封装和业务绑定

- 方法
    - 启动链接(Start)
    - 停止链接(Stop)
    - 获取当前链接的conn(TCPConnetion)
    - 获取链接ID(ConnId)
    - 得到客户端链接地址和端口(RemoteAddr)
    - 发送数据的方法(Send)
- 属性
    - socket TCP套接字(Conn)
    - 链接ID(ConnID)
    - 当前链接的状态是否已经关闭(isClosed)
    - 与当前链接所绑定的处理业务方法(handleAPI)
    - 等待链接被动退出的channel(ExitChan)

## 基础router模块
- Request请求封装
  - 将链接和数据绑定在一起
    - 属性
      - 链接(conn) 
      - 请求数据(data)
    - 方法
      - 得到当前链接(Connection)
      - 得到当前数据(Data)
- Router模块
  - 抽象Router
    - 处理业务之前的方法(PreHandle)
    - 处理业务的主方法(Handle)
    - 处理业务之后的方法(PostHandle)
  - 具体的BaseRouter(具体实现)
      - 处理业务之前的方法
      - 处理业务的主方法
      - 处理业务之后的方法
- xin集成router模块
  - XServer新增一个路由添加功能(++AddRouter)
  - Server类新增一个Router成员(--HandleAPI)
  - Connection类绑定一个Router
  - 在Connection调用已经注册的Router处理业务
- 使用xinV0.3开发
  1. 实现Router(直接拷贝BaseRouter的写法就行了)
  2. 创建Server
  3. 添加实现的路由
  4. 启动
## 全局配置
- 通过yaml进行配置
- 创建一个global_config结构体,读取yaml配置
- 替换xin中的配置
- 打包到github

## 消息封装
- 定义一个消息结构体(Message)
  - 属性
    - 消息ID
    - 消息长度
    - 消息内容
- 定义一个解决TCP粘包问题的封包拆包的模块
  - 针对Message进行TLV(type(id)、lenght、value)格式的封装(Pack)
    - 写消息长度
    - 写消息ID
    - 写消息内容
  - 针对Message进行TLV格式的拆包(UnPack)
    - 先读取固定长度的head->消息内容的长度和消息类型
    - 再根据消息内容的长度,再进行一次读写,从conn中读取消息内容
- 集成消息封装到Xin中
  - 将Message添加到Request中
  - 修改链接的读取机制,把单纯的读取byte改成拆包
  - 给Connection提供一个发包机制:将发送的消息封包后发送