package main

import (
	"context"
	"fmt"
	stdlog "log"

	"xiaoifei.top/registry"
	"xiaoifei.top/separation"
	"xiaoifei.top/service"
)

func main() {
	host, port := "localhost", "6000"
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		registry.Registration{
			ServiceName: registry.SeperationService,
			ServiceURL:  fmt.Sprintf("http://%s:%s", host, port),
		},
		separation.RegisterHandlers,
	)
	if err != nil {
		stdlog.Fatalln(err)
	}

	<-ctx.Done()
	fmt.Println("Shutdown Service")
}
