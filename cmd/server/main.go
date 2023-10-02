package main

import (
	"context"
	"github.com/h2p2f/dedicated-vault/internal/server/app"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()
	sigint := make(chan os.Signal, 1)
	connectionsClosed := make(chan struct{})
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	app.Run(ctx, sigint, connectionsClosed)

	<-connectionsClosed
}
