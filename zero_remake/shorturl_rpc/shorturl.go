package main

import (
	"flag"
	"fmt"

	"example.com/shorturl/short-url/zero_remake/shorturl_rpc/internal/config"
	"example.com/shorturl/short-url/zero_remake/shorturl_rpc/internal/server"
	"example.com/shorturl/short-url/zero_remake/shorturl_rpc/internal/svc"
	"example.com/shorturl/short-url/zero_remake/shorturl_rpc/types/shortUrl"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "D:\\shorturl\\short-url\\zero_remake\\shorturl_rpc\\etc\\shorturl.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		shortUrl.RegisterShortUrlServer(grpcServer, server.NewShortUrlServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
