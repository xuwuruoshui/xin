package main

import "xin/xnet"

func main() {
	// 1、创建xin的server
	s := xnet.NewServer("[xin V0.2]")

	// 2、启动xin的server
	s.Run()
}