package main

import (
	"flag"
	"httpProxy/httpProxy"

	"net/http"
	_ "net/http/pprof" // 自动注册 /debug/pprof/*

	"github.com/spf13/viper"
)

type Config struct {
	Port        string
	ProxyConfig []httpProxy.ProxyConfig
}

func getConfig(configPath string) (Config, error) {
	var conf Config
	viper.SetConfigName("config")
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	viper.Unmarshal(&conf)
	return conf, err
}

func main() {
	var configPath = flag.String("path", "./", "config's path")
	flag.Parse()

	conf, err := getConfig(*configPath)
	if err != nil {
		panic(err)
	}
	httpProxy.Serve(conf.ProxyConfig, conf.Port)
}

func prof() {
	go func() {
		http.ListenAndServe(":6060", nil) // 独占端口，避免与业务路由冲突
	}()
}
