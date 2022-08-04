package main

import (
	"fmt"
	"github.com/lucky-cheerful-man/phoenix_server/pkg/gmysql"
	"github.com/lucky-cheerful-man/phoenix_server/pkg/gredis"
	"github.com/lucky-cheerful-man/phoenix_server/pkg/log"
	"github.com/lucky-cheerful-man/phoenix_server/pkg/rpc"
	"github.com/lucky-cheerful-man/phoenix_server/pkg/setting"
	"github.com/lucky-cheerful-man/phoenix_server/service"
	"google.golang.org/grpc" //nolint:goimports
	"net"
	"os"
	"os/signal"
	"syscall" //nolint:goimports
)

func dealSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	go func() {
		for s := range sigs {
			switch s {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT:
				log.Warnf("got signal:%v and try to exit: ", s)
				os.Exit(0)
			default:
				log.Warnf("other signal:%v: ", s)
			}
		}
	}()
}

func main() {
	// 初始化配置与日志
	setting.InitConfig()
	log.InitLog()

	// 监听处理信号
	dealSignal()

	// 初始化服务
	server := grpc.NewServer()
	rpc.RegisterUserServiceServer(server, &service.UserService{DB: gmysql.Setup(), Cache: gredis.Setup()})
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", setting.ServerSetting.Port))
	if err != nil {
		log.Fatalf("net.Listen err: %s", err)
	}
	log.Infof("start service")
	err = server.Serve(lis)
	if err != nil {
		log.Warnf("server init failed:%s", err)
	}
}
