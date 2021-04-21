package sock

import (
	"errors"
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"os"
	httpsocks "socks5_go/http"
	"strconv"
	_ "socks5_go/http"

)

type MySocks5 interface {
	ServerAndListen() (interface{}, interface{})
	AuthSocks5(client net.Conn) interface{}
	ConnectSocks5()
	GetClientCallInfo(client net.Conn) (string, string)
	ForwardRequest(host string, port string, client net.Conn) interface{}
	ProcessSocks5(client net.Conn) interface{}
}

type MyConfig struct {
	Post string
	Port string
}

func AuthSocks5(client net.Conn) (interface{}, interface{}) {
	buf := make([]byte, 256)

	// 读取 VER 和 nMETHODS
	n, err := io.ReadFull(client, buf[:2])
	if n != 2 {
		return nil, errors.New("reading header: " + err.Error())
	}

	ver, nMethods := int(buf[0]), int(buf[1])
	if ver != 5 {
		return nil, errors.New("invalid version")
	}

	// 读取 METHODS 列表
	n, err = io.ReadFull(client, buf[:nMethods])
	if n != nMethods {
		return nil, errors.New("reading methods: " + err.Error())
	}

	// 不认证
	//n, err = client.Write([]byte{0x05, 0x00})

	//认证
	n, err = client.Write([]byte{0x05, 0x02})
	if n != 2 || err != nil {
		return nil, errors.New("write rsp: " + err.Error())
	}

	wBuff := make([]byte, 1024)
	wn, errReadBuff := client.Read(wBuff[:])
	if errReadBuff != nil {
		fmt.Println("授权阶段出现问题", errReadBuff)
		client.Write([]byte{0x05, 0x01})
		client.Close()
		return nil, nil
	}
	client.Write([]byte{0x05, 0x00})
	uLen := int(wBuff[1])      // 用户长度
	pLen := int(wBuff[2+uLen]) // 密码长度
	//fmt.Println("用户长度:", uLen)
	//fmt.Println("密码长度", pLen)
	uname := string(wBuff[2 : 2+uLen])
	passwd := string(wBuff[wn - pLen : wn])
	//fmt.Println("用户名：", uname)
	//fmt.Println("密码:", passwd)
	checkErr := httpsocks.CheckUser(uname,passwd)
	if checkErr != nil  {
		fmt.Println("User check fail\n")
		client.Close()
		return nil, checkErr
	}
	return nil, nil
}

func (s MyConfig) ServerAndListen() (interface{}, interface{}) {
	fmt.Println("服务启动中。。。")
	server, err := net.Listen("tcp", ":1080")
	if err != nil {
		fmt.Println("服务启动失败:", "您可能多次启动本程序，或服务端口被占用")
		os.Exit(0)
	}
	fmt.Println("启动成功")
	for true {
		client, err := server.Accept()
		fmt.Println("请求ip:", client.RemoteAddr())
		if err != nil {
			fmt.Println("建立连接发生错误：", err)
		}
		fmt.Println("连接ok")

		go ProcessSocks5(client)
	}
	return nil, nil
}

func ForwardRequest(host string, port string, client net.Conn) interface{} {
	// socks5  上游代理
	socksServer, err := proxy.SOCKS5("tcp", "ssr.comeboy.cn:2933", nil, proxy.Direct)
	if err != nil {
		fmt.Println("GG 初始化代理失败！")
		return nil
	}

	server, errDial := socksServer.Dial("tcp", net.JoinHostPort(host, port))
	if errDial != nil {
		fmt.Println("使用代理访问出错！")
		return nil
	}
	//响应客户端连接成功
	client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	forward := func(src,dest net.Conn) {
		defer src.Close()
		defer dest.Close()
		io.Copy(src,dest)
	}

	go forward(client,server)
	go forward(server,client)

	fmt.Println("Function Shutdown")
	return nil
}

func GetClientCallInfo(client net.Conn) (string, string, interface{}) {
	var host, port string
	buff := make([]byte, 256)
	n, err := client.Read(buff[:])
	if err != nil {
		fmt.Println("解析请求不对劲：", err)
		return host, port, err
	}
	//fmt.Println("---------------解析请求信息中---------------")

	switch buff[3] {
	case 0x01:
		host = net.IPv4(buff[4], buff[5], buff[6], buff[7]).String()
	case 0x03:
		host = string(buff[5 : n-2])
	default:
		client.Close()
		return host, port, err
	}
	port = strconv.Itoa(int(buff[n-2])<<8 | int(buff[n-1]))

	if len(host) < 5 || len(port) == 0 {
		return host, port, errors.New("域名或端口不正确")
	}

	//fmt.Println("访问信息")
	//fmt.Println("域名： ", host)
	//fmt.Println("端口：", port)
	return host, port, nil
}

func ProcessSocks5(client net.Conn) {
	_, err := AuthSocks5(client)
	if err != nil {
		//fmt.Println("发生错误:", err)
	} else {
		host, port, err := GetClientCallInfo(client)
		//host = "qd.hlwaqxz.cn"
		//port ="443"
		if err != nil {
			fmt.Println(err)
		} else {
			ForwardRequest(host, port, client)
		}
	}

}
