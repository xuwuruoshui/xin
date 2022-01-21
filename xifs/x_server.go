package xifs

type XServer interface {

	// 启动服务器
	Start()

	// 停止服务器
	Stop()

	// 运行服务器
	Run()

	// 当前的Server添加一个MessageHandler,用于msgId和router的绑定
	AddRouter(msgId uint32, router XRouter)

	// 获取链接管理
	GetConnMgr() XConnectionManager

	// 注册OnConnStart方法
	SetOnConnStart(func(conn XConnection))

	// 注册OnConnStop方法
	SetOnConnStop(func(conn XConnection))

	// 注册OnConnStart方法
	CallOnConnStart(conn XConnection)

	// 注册OnConnStop方法
	CallOnConnStop(conn XConnection)
}
