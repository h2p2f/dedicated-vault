package main

import (
	"context"
	"github.com/h2p2f/dedicated-vault/internal/client/app"
)

func main() {
	ctx := context.Background()
	app.Run(ctx)
}
