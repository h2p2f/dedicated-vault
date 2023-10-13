package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/h2p2f/dedicated-vault/internal/client/app"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(homeDir)
	ctx := context.Background()
	app.Run(ctx)
}
