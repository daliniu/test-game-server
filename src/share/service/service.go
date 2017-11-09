package service

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"share/config"
	"share/wrapper"
	"strconv"
	"strings"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	microError "github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/transport"
	"github.com/micro/go-os/metrics"
	"github.com/micro/go-os/trace"
	"github.com/micro/go-plugins/metrics/statsd"
	"github.com/micro/go-plugins/trace/zipkin"
	"github.com/micro/misc/lib/addr"
	"golang.org/x/net/context"
)

// NewService 创建服务
func NewService(serviceName, version string, isRoot bool, flags ...cli.Flag) micro.Service {
	cfg := config.Instance()

	// 两分钟读写超时
	server.DefaultServer.Init(server.Transport(transport.NewTransport(transport.Timeout(2 * time.Minute))))
	// 服务发现的服务列表数据缓存时间
	// cache.DefaultTTL = 10 * time.Minute
	// 默认客户端连接池
	client.DefaultPoolSize = 20
	// 设置客户端重试
	client.DefaultRetries = 2
	client.DefaultRetry = func(ctx context.Context, req client.Request, retryCount int, err error) (bool, error) {
		if err == nil {
			return false, nil
		}
		e := microError.Parse(err.Error())
		if e.Code < 600 {
			return true, nil
		}
		return false, nil
	}

	gracefulRestartWrapper := wrapper.GracefulRestartWrapper{}

	clientWrappers := []client.Wrapper{wrapper.NewHystrixClientWrapper()}
	serverWrappers := []server.HandlerWrapper{gracefulRestartWrapper.HandlerWrapper()}

	// 初始化统计
	statsdBrokerAddress := cfg.Get("statsd-broker").String("")
	if len(statsdBrokerAddress) > 0 {
		m := statsd.NewMetrics(metrics.Collectors())
		stats := wrapper.NewServerStats(m, config.CurrentIP, serviceName, server.DefaultId)
		stats.Run()
		serverWrappers = append(serverWrappers, stats.HandlerStatsWrapper())
	}

	// 初始化trace
	traceKafkaAddress := cfg.Get("trace-kafka").StringSlice([]string{})
	if len(traceKafkaAddress) > 0 {
		t := zipkin.NewTrace(trace.Collectors(traceKafkaAddress...))
		srv := &registry.Service{Name: serviceName, Nodes: []*registry.Node{{Address: config.CurrentIP, Port: 8888}}}
		clientWrappers = append(clientWrappers, trace.ClientWrapper(t, srv, cfg, isRoot))
		serverWrappers = append(serverWrappers, trace.HandlerWrapper(t, srv, cfg))
	}

	service := micro.NewService(
		micro.Name(serviceName),
		micro.Version(version),
		micro.RegisterInterval(5*time.Second),
		micro.RegisterTTL(10*time.Second),
		micro.WrapClient(clientWrappers...),
		micro.WrapHandler(serverWrappers...),
		micro.Flags(flags...),

		micro.AfterStop(func() error {
			gracefulRestartWrapper.Wait()
			time.Sleep(3 * time.Second)
			return nil
		}),
	)

	// 写入pid
	srvName, ok := config.ServiceNameAlias[serviceName]
	if !ok {
		srvName = strings.Replace(serviceName, config.Namespace, "", 1)
	}
	WritePid(srvName)
	return service
}

// WritePid 写入当前进程的pid
func WritePid(srvName string) {
	pidfilepath := config.Instance().Get("WWHOME").String("")
	if len(pidfilepath) == 0 {
		pidfilepath = path.Join("..", "pidfile", srvName+".pid")
	} else {
		pidfilepath = path.Join(pidfilepath, "pidfile", srvName+".pid")
	}
	if err := ioutil.WriteFile(pidfilepath, []byte(strconv.Itoa(os.Getpid())), 0766); err != nil {
		panic("写入pid失败:" + err.Error())
	}
}

// GetServerAddr GetServerAddr
func GetServerAddr(s server.Server) (string, error) {
	if len(s.Options().Address) > 0 {
		idx := strings.LastIndex(s.Options().Address, ":")
		ip, err := addr.Extract("")
		if err != nil {
			return "", err
		}
		if ip == "" {
			return "", errors.New("无法获取当前监听地址")
		}
		addr := ip + s.Options().Address[idx:]
		return addr, nil
		//log.Printf("房间服务器初始化地址：m.addr:%+v ip:%+v address:%+v idx:%+v", addr, ip, s.Options().Address, idx)
	}
	return "", errors.New("无法获取当前监听地址")
}
