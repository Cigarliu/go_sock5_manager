package sck5log

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

const (
	fileName = "sck5log.log"      // 日志文件前缀名称
	logLevel = logrus.DebugLevel  // 日志级别设置
)

func init() {
	//设置输出样式，自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}
	logrus.SetFormatter(&logrus.TextFormatter{})
	//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
	logrus.SetOutput(os.Stdout)
	/* 设置日志文件名*/
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// 将文件设置为终端输出流 + 文件重定向
	writers := []io.Writer{file, os.Stdout}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	if err == nil {
		logrus.SetOutput(fileAndStdoutWriter)
	} else {
		logrus.Info("failed to log to file.")
	}
	logrus.SetLevel(logLevel)
}
