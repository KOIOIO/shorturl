package main

import (
	"example.com/shorturl/short-url/zero_remake/middleware"
	"example.com/shorturl/short-url/zero_remake/shorturl_api/internal/config"
	"example.com/shorturl/short-url/zero_remake/shorturl_api/internal/handler"
	"example.com/shorturl/short-url/zero_remake/shorturl_api/internal/svc"
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/shorturl-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	server.Use(middleware.CorsMiddleware())
	server.Use(middleware.LoggerMiddleware(c.Log.Path))

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
