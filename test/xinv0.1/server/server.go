package main

import "github.com/xuwuruoshui/xin/xnet"

func main() {
	// 1、创建xin的server
	s := xnet.NewServer("[xin V0.1]")

	// 2、启动xin的server
	s.Run()
}
