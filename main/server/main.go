package main

import (
	"fmt"
	"github.com/zhuozl/go-proxy/common"
	"github.com/zhuozl/go-proxy/server"
	"log"
	"net"
)

var version = "master"

const (
	configPath = "server.conf"
	defalutPort=1111
)

func main() {
	log.SetFlags(log.Lshortfile)

	// 默认配置
	config := &common.Config{
		ListenAddr: fmt.Sprintf(":%d", defalutPort),
		// 密码随机生成
		Password: common.RandPassword().String(),
	}
	config.ReadConfig(configPath)
	config.SaveConfig(configPath)

	// 解析配置
	password, err := common.ParsePassword(config.Password)
	if err != nil {
		log.Fatalln(err)
	}
	listenAddr, err := net.ResolveTCPAddr("tcp", config.ListenAddr)
	if err != nil {
		log.Fatalln(err)
	}

	// 启动 server 端并监听
	lsServer := server.New(password, listenAddr)
	log.Fatalln(lsServer.Listen(func(listenAddr net.Addr) {
		log.Println("使用配置：", fmt.Sprintf(`
本地监听地址 listen：
%s
密码 password：
%s
	`, listenAddr, password))
		log.Printf("lightsocks-server:%s 启动成功 监听在 %s\n", version, listenAddr.String())
	}))
}
