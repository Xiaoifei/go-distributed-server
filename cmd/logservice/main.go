package main

import (
	"context"
	"fmt"
	stdlog "log"

	"xiaoifei.top/log"
	"xiaoifei.top/registry"
	"xiaoifei.top/service"
)

func main() {
	log.Run("./distributed.log")
	host, port := "localhost", "4000"
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		registry.Registration{
			ServiceName: "Log Service",
			ServiceURL:  fmt.Sprintf("http://%s:%s", host, port),
		},
		log.RegisterHandelers,
	)
	if err != nil {
		stdlog.Fatalln(err)
	}

	<-ctx.Done()
	fmt.Println("Shutdown Service")
}
