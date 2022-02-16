package xnet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/xuwuruoshui/xin/config"
	"github.com/xuwuruoshui/xin/xifs"
)

// 封包、拆包,解决粘包问题

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包的头的长度方法
func (d *DataPack) GetHeadLen() uint32 {
	//len + id = 4+4字节
	return 8
}

// 封包方法
func (d *DataPack) Pack(msg xifs.XMessage) ([]byte, error) {
	// 创建一个存放bytes的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// len写入buff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetLength()); err != nil {
		return nil, err
	}
	// id写入buff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetId()); err != nil {
		return nil, err
	}
	// data写入buff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包方法
func (d *DataPack) Unpack(binaryData []byte) (xifs.XMessage, error) {
	dataBuff := bytes.NewBuffer(binaryData)

	msg := &Message{}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Length); err != nil {
		return nil, err
	}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	if config.GloabalConf.MaxPackageSize < 0 && msg.Length > config.GloabalConf.MaxPackageSize {
		return nil, errors.New("too large Msg data recived!!!")
	}

	return msg, nil
}
