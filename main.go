package main

import (
	"fmt"
	_ "socks5_go/http"
	_ "socks5_go/sck5log"
	"socks5_go/sock"
)
func main() {
	//logrus.Error("****** 错误日志日志测试 ******")
	//logrus.Info("****** 信息日志日志测试 ******")
	//logrus.Debug("****** 调试日志日志测试 ******")

	//fmt.Print("************* 欢迎使用  ************* \n" +
	//	"请使用SwitcheyOmega进行连接\n" +
	//	"连接地址：127.0.0.1  端口: 1080\n" +
	//	"************************************\n")

	//httpsocks.WebStart()

	//httpsocks.Get("http://ssr.comeboy.cn:8989/login")
	//httpsocks.CheckUser("cigar","123456")
	var a sock.MyConfig
	a.ServerAndListen()
	var abc string
	fmt.Print("按任意键结束")
	fmt.Scanln(&abc)
}
