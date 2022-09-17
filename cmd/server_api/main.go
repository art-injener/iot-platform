package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/art-injener/iot-platform/internal/webserver"
)

func main() {
	time.LoadLocation(time.UTC.String())
	ctx, cancel := context.WithCancel(context.Background())
	application := webserver.App{}
	application.Initialize(ctx)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(quit)
		<-quit
		cancel()
	}()
	application.Run(ctx)
}
