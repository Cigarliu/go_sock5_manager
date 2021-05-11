package main

import (
	"fmt"
	_ "socks5_go/http"
	httpsocks "socks5_go/http"
	_ "socks5_go/sck5log"
	"socks5_go/sock"
)
func main() {
	//logrus.Error("****** 错误日志日志测试 ******")
	//logrus.Info("****** 信息日志日志测试 ******")
	//logrus.Debug("****** 调试日志日志测试 ******")

	fmt.Print("************* 欢迎使用  ************* \n" +
		"请使用SwitcheyOmega进行连接\n" +
		"连接地址：127.0.0.1  端口: 1080\n" +
		"************************************\n")
	httpsocks.InitDB()

	//httpsocks.GetUserInfo("cidgar")
	go httpsocks.WebStart()

	var a sock.MyConfig
	a.ServerAndListen()
	//a.Port = ":10081"
	//a.ServerAndListen()
}
