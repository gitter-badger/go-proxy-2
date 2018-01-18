package main

import (
	"io"
	"log"
	"net"
	"strconv"
	"bytes"
	"fmt"
	"encoding/binary"
)

func main() {



	buff := new(bytes.Buffer)

	port,_ :=strconv.Atoi("1111")

	err := binary.Write(buff, binary.BigEndian, uint16(port))
	if err != nil {
		fmt.Println(err)
	}


	intByteArray := buff.Bytes()
	fmt.Printf("intByteArray : % x\n", intByteArray)

	buf := make([]byte, 4)
	n := binary.PutVarint(buf, 1111)
	b := buf[:n]
	fmt.Println(b)

	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, 888)
	fmt.Println(bs)

	return ;
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	l, err := net.Listen("tcp", ":1080")
	if err != nil {
		log.Panic(err)
	}

	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}

		go handleClientRequest(client)
	}
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

func handleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()

	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}

	if b[0] == 0x05 { //只处理Socks5协议
		//客户端回应：Socks服务端不需要验证方式
		client.Write([]byte{0x05, 0x00})
		n, err = client.Read(b[:])
		var host, port string
		switch b[3] {
		case 0x01: //IP V4
			host = net.IPv4(b[4], b[5], b[6], b[7]).String()
		case 0x03: //域名
			host = string(b[5 : n-2]) //b[4]表示域名的长度
		case 0x04: //IP V6
			host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
		}
		port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))

		server, err := net.Dial("tcp", net.JoinHostPort(host, port))
		if err != nil {
			log.Println(err)
			return
		}
		defer server.Close()
		client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) //响应客户端连接成功

		log.Printf("server:%s,client:%s",server.LocalAddr().String(),client.LocalAddr().String())
		//进行转发
		go io.Copy(server, client)
		io.Copy(client, server)
	}

}
