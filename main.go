package main

import (
	"github.com/asim/go-micro/plugins/registry/etcd/v3"
	goMirco "github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	pb "github.com/lucky-cheerful-man/phoenix_apis/protobuf3.pb/user_info_manage"
	"github.com/lucky-cheerful-man/phoenix_server/pkg/gmysql"
	"github.com/lucky-cheerful-man/phoenix_server/pkg/gredis"
	"github.com/lucky-cheerful-man/phoenix_server/pkg/log"
	"github.com/lucky-cheerful-man/phoenix_server/pkg/setting"
	"github.com/lucky-cheerful-man/phoenix_server/service"
)

func main() {
	etcdReg := etcd.NewRegistry(
		registry.Addrs(setting.ReferGlobalConfig().ServerSetting.RegisterAddress),
	)
	// 初始化服务
	srv := goMirco.NewService(
		goMirco.Name(setting.ReferGlobalConfig().ServerSetting.RegisterServerName),
		goMirco.Version(setting.ReferGlobalConfig().ServerSetting.RegisterServerVersion),
		goMirco.Registry(etcdReg),
	)

	err := pb.RegisterUserServiceHandler(srv.Server(), &service.UserService{DB: gmysql.Setup(),
		Cache: gredis.Setup()})
	if err != nil {
		log.Errorf("RegisterUserServiceHandler failed, err:%s", err)
	}

	err = srv.Run()
	if err != nil {
		log.Errorf("run failed, err:%s", err)
	}
}
