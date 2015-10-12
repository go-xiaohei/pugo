package core

import (
	"fmt"
	"gopkg.in/inconshreveable/log15.v2"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strings"
	"syscall"
)

type Service interface {
	Start()
	Stop()
}

// 启动服务，支持ctrl+c监听
func Start(s Service) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	signalChan := make(chan os.Signal)
	sName := reflect.TypeOf(s).String()

	// 监听关闭符号
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	log15.Info(fmt.Sprintf("%s.Start", strings.Title(sName)))

	// 启动服务，并监听关闭
	s.Start()
	<-signalChan
	s.Stop()
	WrapWait() // 等待全局的wg结束

	log15.Info(fmt.Sprintf("%s.Stop", strings.Title(sName)))
}
