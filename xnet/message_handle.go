package xnet

import (
	"fmt"
	"github.com/xuwuruoshui/xin/config"
	"github.com/xuwuruoshui/xin/xifs"
	"log"
)

type MessageHandle struct {
	// 存放每个msgId所对应的处理方法
	Apis map[uint32]xifs.XRouter

	// Worker池的大小
	WorkerPoolSize uint32

	// Worker取任务的消息队列
	TaskQueue []chan xifs.XRequest
}

func NewMsgHandle() *MessageHandle {
	return &MessageHandle{
		Apis: make(map[uint32]xifs.XRouter),
		// 从全局配置中获取
		WorkerPoolSize: config.GloabalConf.WorkerPoolSize,
		TaskQueue:      make([]chan xifs.XRequest, config.GloabalConf.WorkerPoolSize),
	}
}

// 调度/执行对应的Router消息处理方法
func (m *MessageHandle) DoMsgHandle(request xifs.XRequest) {
	//1、
	handler, ok := m.Apis[request.MsgId()]
	if !ok {
		log.Printf("Api msgId=%d is not found!!! Did you add the router?\n", request.MsgId())
		return
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (m *MessageHandle) AddRouter(msgId uint32, router xifs.XRouter) {
	// 1、判断当前msgId绑定的API处理方法是否已存在
	if _, ok := m.Apis[msgId]; ok {
		panic(fmt.Sprintf("repeat api,msgId=%d", msgId))
	}

	//2、添加msg与API的绑定关系
	m.Apis[msgId] = router
	log.Printf("Add api MsgId=%d successed!!!\n", msgId)
}

// 启动一个Worker工作池
func (m *MessageHandle) StartWorkerPool() {
	// 更具workerPoolSize分别开启Worker,每个Worker用一个go来承载
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 创建通道
		m.TaskQueue[i] = make(chan xifs.XRequest, config.GloabalConf.MaxWorkerTaskSize)
		// 启动worker
		go m.StartWorkflow(i, m.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (m *MessageHandle) StartWorkflow(workerId int, taskQueue chan xifs.XRequest) {
	log.Printf("Worker Id=%d is started...\n", workerId)

	for request := range taskQueue {
		m.DoMsgHandle(request)
	}
}

// 将消息交给TaskQueue,由Worker进行处理
func (m *MessageHandle) SendMsgToTaskQueue(request xifs.XRequest) {
	// 1、将消息平均分配给不同的worker
	workerId := request.Connection().ConnId() % m.WorkerPoolSize
	log.Printf("Add ConnId=%d,request MsgId=%d,WorkerId=%d\n", request.Connection().ConnId(), request.MsgId(), workerId)
	// 2、将消息发送给对应的worker的TaskQueue即可
	m.TaskQueue[workerId] <- request
}
