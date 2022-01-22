package config

import (
	"github.com/xuwuruoshui/xin/xifs"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type GloabalConfig struct {
	// 当前server对象
	Server  xifs.XServer
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	// 最大连接数
	MaxConn uint `yaml:"maxConn"`
	// 发送一次数据包的最大值
	MaxPackageSize uint32 `yaml:"maxPackageSize"`
	// 工作池channel数量
	WorkerPoolSize uint32 `yaml:"workerPoolSize"`
	// 一个工作池最大容量
	MaxWorkerTaskSize uint32
}

var GloabalConf *GloabalConfig

func init() {
	// 默认值
	GloabalConf = &GloabalConfig{
		Name:              "Xin Server",
		Version:           "V1.0",
		Port:              9998,
		Host:              "0.0.0.0",
		MaxConn:           1000,
		MaxPackageSize:    4096,
		WorkerPoolSize:    10,
		MaxWorkerTaskSize: 1024,
	}
	// 从yaml中去加载用户自定义信息
	GloabalConf.Reload()
}

func (g *GloabalConfig) Reload() {
	data, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		log.Println("Read config error:", err)
	}
	err = yaml.Unmarshal(data, g)
	if err != nil {
		log.Println("Data exchange error:", err)
	}

}
