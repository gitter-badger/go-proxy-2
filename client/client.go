package client

import (
	"github.com/zhuozl/go-proxy/common"
	"log"
	"net"
	"fmt"
	"bytes"
	"strconv"
	"encoding/binary"
)

type LsLocal struct {
	*common.SecureSocket
	Config common.Config
}

// 新建一个本地端
// 本地端的职责是:
// 1. 监听来自本机浏览器的代理请求
// 2. 转发前加密数据
// 3. 转发socket数据到墙外代理服务端
// 4. 把服务端返回的数据转发给用户的浏览器
func New(password *common.Password, listenAddr, remoteAddr *net.TCPAddr) *LsLocal {
	return &LsLocal{
		SecureSocket: &common.SecureSocket{
			Cipher:     common.NewCipher(password),
			ListenAddr: listenAddr,
			RemoteAddr: remoteAddr,
		},
	}
}

// 本地端启动监听，接收来自本机浏览器的连接
func (local *LsLocal) Listen(didListen func(listenAddr net.Addr)) error {
	listener, err := net.ListenTCP("tcp", local.ListenAddr)
	if err != nil {
		return err
	}

	defer listener.Close()

	if didListen != nil {
		didListen(listener.Addr())
	}

	for {
		userConn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		// userConn被关闭时直接清除所有数据 不管没有发送的数据
		userConn.SetLinger(0)
		go local.handleConn(userConn)
	}
	return nil
}


//byte转16进制字符串
func ByteToHex(data []byte) string {
	buffer := new(bytes.Buffer)
	for _, b := range data {

		s := strconv.FormatInt(int64(b&0xff), 16)
		if len(s) == 1 {
			buffer.WriteString("0")
		}
		buffer.WriteString(s)
	}

	return buffer.String()
}

func (local *LsLocal) handleConn(userConn *net.TCPConn) {
	defer userConn.Close()

	proxyServer, err := local.DialRemote()
	if err != nil {
		log.Println(err)
		return
	}
	defer proxyServer.Close()
	// Conn被关闭时直接清除所有数据 不管没有发送的数据
	proxyServer.SetLinger(0)


	//进行二级穿透
	if(local.Config.Proxy!=""){
		proxyIp ,proxyPort ,_:= net.SplitHostPort(local.Config.Proxy)

		ip :=net.ParseIP(proxyIp)
		buff := new(bytes.Buffer)

		_port,_ :=strconv.Atoi(proxyPort)
		err = binary.Write(buff, binary.BigEndian, uint16(_port))
		if err != nil {
			fmt.Println(err)
		}

		var b [1024]byte

		//加密
		//proxyServer.Write([]byte{0x05,0x01,0x02})
		proxyServer.Write([]byte{0x05,0x00})
		proxyServer.Read(b[:])
		fmt.Println("X-1")
		fmt.Println(ByteToHex(b[:]))

		//如果是加密方式
		//proxyServer.Write([]byte{0x01,0x06,0x7a,0x68,0x75,0x6f,0x7a,0x6c,0x07,0x30,0x31,0x31,0x32,0x32,0x33,0x33});

		//proxyServer.Read(b[:])
		//fmt.Println("X-2")
		//fmt.Println(ByteToHex(b[:]))
		proxyServer.Write([]byte{0x05,0x01,0x00,0x01,ip[12],ip[13],ip[14],ip[15],buff.Bytes()[0],buff.Bytes()[1]})

		proxyServer.Read(b[:])
		fmt.Println("X-3")
		fmt.Println(ByteToHex(b[:]))
	}

	// 进行转发
	// 从 proxyServer 读取数据发送到 localUser
	go func() {
		err := local.DecodeCopy(userConn, proxyServer)
		if err != nil {
			fmt.Println(err.Error())
			// 在 copy 的过程中可能会存在网络超时等 error 被 return，只要有一个发生了错误就退出本次工作
			userConn.Close()
			proxyServer.Close()
		}
	}()
	// 从 localUser 发送数据发送到 proxyServer，这里因为处在翻墙阶段出现网络错误的概率更大
	local.EncodeCopy(proxyServer, userConn)

}
