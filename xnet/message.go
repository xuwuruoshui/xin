package xnet

type Message struct {
	// 消息Id
	Id uint32
	// 消息长度
	Length uint32
	// 消息内容
	Data []byte
}

func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		Id:     id,
		Length: uint32(len(data)),
		Data:   data,
	}
}

func (m *Message) GetId() uint32 {
	return m.Id
}
func (m *Message) GetLength() uint32 {
	return m.Length
}
func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetId(id uint32) {
	m.Id = id
}
func (m *Message) SetLength(len uint32) {
	m.Length = len
}
func (m *Message) SetData(data []byte) {
	m.Data = data
}
