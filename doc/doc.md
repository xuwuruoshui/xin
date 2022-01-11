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