package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"golang.org/x/net/proxy"

)

type Mysocks5 interface {
	ServerAndListen()(interface{},interface{})
	AuthSocks5(client net.Conn)(interface{})
	ConnectSocks5()
	GetClientCallInfo(client net.Conn)(string,string)
	ForwardRequest(host string,port string,client net.Conn)(interface{})
	ProcessSocks5(client net.Conn)(interface{})
}

type  MyConfig struct {
	host string
	port string
}

func AuthSocks5(client net.Conn)(interface{},interface{})  {
	buf := make([]byte, 256)

	// 读取 VER 和 NMETHODS
	n, err := io.ReadFull(client, buf[:2])
	if n != 2 {
		return nil,errors.New("reading header: " + err.Error())
	}

	ver, nMethods := int(buf[0]), int(buf[1])
	if ver != 5 {
		return nil,errors.New("invalid version")
	}

	// 读取 METHODS 列表
	n, err = io.ReadFull(client, buf[:nMethods])
	if n != nMethods {
		return nil,errors.New("reading methods: " + err.Error())
	}

	//认证
	n, err = client.Write([]byte{0x05, 0x02})
	if n != 2 || err != nil {
		return nil,errors.New("write rsp: " + err.Error())
	}

	wbuff := make([]byte,1024)
	wn,err := client.Read(wbuff[:])

	if err != nil{
		fmt.Printf("授权阶段出现问题",err)
		client.Write([]byte{0x05,0x01})
		client.Close()
		return nil,nil
	}
	client.Write([]byte{0x05,0x00})
	ulen := int(wbuff[1])  // 用户长度
	plen := int(wbuff[2+ulen]) // 密码长度
	fmt.Printf("\n用户长度:",ulen)
	fmt.Printf("\n密码长度",plen)
	uname := string(wbuff[2:2+ulen])
	passwd := string(wbuff[3+plen:wn])
	fmt.Printf("\n用户名：",uname)
	fmt.Printf("\n密码:",passwd)
	return nil,nil
}



func (s MyConfig) ServerAndListen()(interface{},interface{}) {
	fmt.Print("\n服务启动中。。。")
	server,err := net.Listen("tcp",":1080")
	if err != nil{
		fmt.Printf("\n服务器监听出现错误:",err)
	}
	fmt.Println("\n启动成功")
	for true {
		client, err := server.Accept()
		if err != nil{
			fmt.Printf("\n建立连接发生错误：",err)
		}
		fmt.Print("\n连接ok")

		go  ProcessSocks5(client)
	}
	return nil, nil
}

func ForwardRequest(host string,port string,client net.Conn)(interface{}){

	// socks5  上游代理
	socks_server ,err :=proxy.SOCKS5("tcp","ssr.comeboy.cn:2933",nil,proxy.Direct)
	if err != nil {
		fmt.Print("\nGG 初始化代理失败！")
		return nil
	}

	server ,err := socks_server.Dial("tcp",net.JoinHostPort(host,port))
	if err !=nil {
		fmt.Printf("\n 使用代理访问出错！",err)

	}
	client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) //响应客户端连接成功

	go io.Copy(server,client)
	go io.Copy(client,server)
	return nil
}

func GetClientCallInfo(client net.Conn)(string,string){
	buff :=make([]byte,256)
	n,err :=client.Read(buff[:])
	if err != nil{
		fmt.Printf("\n解析请求不对劲：",err)

	}
	fmt.Println("\n----------解析请求信息中---------------")

	var host ,port string
	switch buff[3] {
	case 0x01:
		fmt.Println("\n阁下访问的-ip-")
		host = net.IPv4(buff[4],buff[5],buff[6],buff[7]).String()
	case 0x03:
		fmt.Println("\n阁下访问的-域名-")

		host = string(buff[5:n-2])
	default:
		fmt.Println("\n无法解析域名，不对劲")
		client.Close()
		return host,port
	}
	port = strconv.Itoa(int(buff[n-2])<<8 | int(buff[n-1]))

	fmt.Print("\n访问信息")
	fmt.Printf("\n域名： ",host)
	fmt.Printf("\n端口：",port)
	return host, port
}

func ProcessSocks5(client net.Conn)  {
	_,err := AuthSocks5(client)
	if err !=nil {
		fmt.Printf("\n发生错误:",err)
	}else {
		host,port := GetClientCallInfo(client)
		ForwardRequest(host,port,client)
	}


}