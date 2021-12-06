package client

import (
	"github.com/astaxie/beego"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/broker/nats"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
	"github.com/wengoldx/wing/tool/hystrix"
)

// NewClient create new client for call any method
func NewClient(sername string) client.Client {
	natsAddr := beego.AppConfig.String("natsAddr")
	etcdAddr := beego.AppConfig.String("etcdAddr")
	reg := etcd.NewRegistry(
		registry.Addrs(etcdAddr),
	)
	bro := nats.NewBroker(
		broker.Addrs(natsAddr),
	)
	hystrix.Configure(sername)
	service := micro.NewService(
		micro.Registry(reg),
		micro.Broker(bro),
		// micro.Client(c),//change selector
		micro.WrapClient(hystrix.NewClientWrapper()),
	)
	service.Init()
	return service.Client()
}
