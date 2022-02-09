package hystrix

import (
	"context"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/eapache/go-resiliency/retrier"
	"github.com/micro/go-micro/client"
	"net"
	"net/http"
	"time"
)

type clientWrapper struct {
	client.Client
}

func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	return hystrix.Do(req.Service()+"."+req.Endpoint(), func() error {
		// add Repeater option
		// retrying 'n' times, and waiting 'amount' time
		re := retrier.New(retrier.ConstantBackoff(3, 100*time.Millisecond), nil)
		err := re.Run(func() error {
			return c.Client.Call(ctx, req, rsp, opts...)
		})
		return err
	}, nil)
}

// NewClientWrapper returns a hystrix client Wrapper.
func NewClientWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		return &clientWrapper{c}
	}
}

//Configure reset hystrix setting
/*
hystrix setting
Timeout：用于设置超时时间，超过该时间没有返回响应，意味着请求失败；
MaxConcurrent：用于设置同一类型请求的最大并发量，达到最大并发量后，接下来的请求会被拒绝；
VolumeThreshold：用于设置指定时间窗口内让断路器跳闸（开启）的最小请求数；
SleepWindow：断路器跳闸后，在此时间段内，新的请求都会被拒绝；
ErrorPercentThreshold：请求失败百分比，如果超过这个百分比，则断路器跳闸。
*/
func Configure(service string) {
	config := hystrix.CommandConfig{
		Timeout:               3000, //timeout 3s
		MaxConcurrentRequests: 100,  //
		ErrorPercentThreshold: 50,   //
	}
	configs := make(map[string]hystrix.CommandConfig)
	configs[service] = config
	hystrix.Configure(configs)
	// create hystrix dashboard
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "8181"), hystrixStreamHandler)
}
