package main

import (
	"fmt"
	"github.com/zhuozl/go-proxy/common"
	"github.com/zhuozl/go-proxy/client"
	"log"
	"net"
)

const (
	DefaultListenAddr = ":7448"
	configPath = "client.conf"
)

var version = "master"

func main() {


	log.SetFlags(log.Lshortfile)

	// 默认配置
	config := &common.Config{
		ListenAddr: DefaultListenAddr,
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
	remoteAddr, err := net.ResolveTCPAddr("tcp", config.RemoteAddr)
	if err != nil {
		log.Fatalln(err)
	}

	// 启动 local 端并监听
	lsLocal := client.New(password, listenAddr, remoteAddr)
	lsLocal.Config = *config
	log.Fatalln(lsLocal.Listen(func(listenAddr net.Addr) {
		log.Println("使用配置：", fmt.Sprintf(`
本地监听地址 listen：
%s
远程服务地址 remote：
%s
密码 password：
%s
	`, listenAddr, remoteAddr, password))
		log.Printf("lightsocks-local:%s 启动成功 监听在 %s\n", version, listenAddr.String())
	}))
}
