package xifs

type XMessage interface {
	GetId() uint32
	GetLength() uint32
	GetData() []byte

	SetId(id uint32)
	SetLength(len uint32)
	SetData(data []byte)
}
