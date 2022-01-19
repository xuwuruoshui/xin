package xnet

import (
	"fmt"
	"github.com/xuwuruoshui/xin/xifs"
	"log"
)

type MessageHandle struct {
	// 存放每个msgId所对应的处理方法
	Apis map[uint32]xifs.XRouter
}

func NewMsgHandle() *MessageHandle {
	return &MessageHandle{
		Apis: make(map[uint32]xifs.XRouter),
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
