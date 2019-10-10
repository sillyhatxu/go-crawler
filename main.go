package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/sillyhatxu/consul-client"
	"github.com/sillyhatxu/environment-config"
	"github.com/sillyhatxu/go-crawler/config"
	"github.com/sillyhatxu/go-crawler/service"
	"github.com/sillyhatxu/logrus-client"
	"github.com/sillyhatxu/logrus-client/filehook"
	"github.com/sillyhatxu/logrus-client/logstashhook"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	hv1 "google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

func init() {
	cfgFile := flag.String("c", "config-local.conf", "config file")
	flag.Parse()
	err := envconfig.ParseEnvironmentConfig(&config.Conf.EnvConfig)
	if err != nil {
		panic(err)
	}
	envconfig.ParseConfig(*cfgFile, func(content []byte) {
		err := toml.Unmarshal(content, &config.Conf)
		if err != nil {
			panic(fmt.Sprintf("unmarshal toml object error. %v", err))
		}
	})
	fields := logrus.Fields{
		"project":  config.Conf.Project,
		"module":   config.Conf.Module,
		"@version": "1",
		"type":     "project-log",
	}
	logrusconf.NewLogrusConfig(
		logrusconf.Fields(fields),
		logrusconf.FileConfig(filehook.NewFileConfig(config.Conf.Log.FilePath, filehook.Open(config.Conf.Log.OpenLogfile))),
		logrusconf.LogstashConfig(logstashhook.NewLogstashConfig(config.Conf.EnvConfig.LogstashURL, logstashhook.Open(config.Conf.Log.OpenLogstash), logstashhook.Fields(fields))),
	).Initial()
}

func main() {
	consulServer := consul.NewConsulServer(
		config.Conf.EnvConfig.ConsulAddress,
		config.Conf.Module,
		config.Conf.Host,
		config.Conf.GRPCPort,
		consul.CheckType(consul.HealthCheckGRPC),
	)
	err := consulServer.Register()
	if err != nil {
		panic(err)
	}
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Conf.GRPCPort))
	if err != nil {
		panic(err)
	}

	go func() {
		server := grpc.NewServer()

		healthServer := health.NewServer()
		healthServer.SetServingStatus(consul.DefaultHealthCheckGRPCServerName, hv1.HealthCheckResponse_SERVING)
		hv1.RegisterHealthServer(server, healthServer)

		err := server.Serve(grpcListener)
		if err != nil {
			panic(err)
		}
	}()

	//schedulerClient, err := scheduler.NewScheduler(&service.CrawlerLongmanWord{Name: "crawler-longman-word"}, scheduler.Start("00:00:00"), scheduler.Interval(2*time.Second))
	//if err != nil {
	//	panic(err)
	//}
	//schedulerClient.Run()
	go service.AutoCrawlerLongmanWord()
	forever := make(chan bool)
	<-forever
}
