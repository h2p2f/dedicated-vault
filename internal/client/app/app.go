// Package: app
// in this file we have main logic for client
package app

import (
	"context"

	"go.uber.org/zap"

	"github.com/h2p2f/dedicated-vault/internal/client/config"
	"github.com/h2p2f/dedicated-vault/internal/client/grpcclient"
	"github.com/h2p2f/dedicated-vault/internal/client/gui"
	"github.com/h2p2f/dedicated-vault/internal/client/storage"
	"github.com/h2p2f/dedicated-vault/internal/client/tlsloader"
	"github.com/h2p2f/dedicated-vault/internal/client/usecase"
)

// Run launches the main client logic
func Run(ctx context.Context) {
	var err error
	// read configuration
	conf := config.NewClientConfig()

	// create logger
	logger := zap.NewExample()
	//
	db := storage.NewClientStorage(logger, conf)
	//load tls
	conf.TLSConfig, err = tlsloader.LoadTLS(conf.ClientCA, conf.ClientCert, conf.ClientKey)
	if err != nil {
		logger.Fatal("tls", zap.Error(err))
	}
	// create grpc client
	tr := grpcclient.NewClient(conf, logger)
	uc := usecase.NewClientUseCase(conf, db, tr)
	// create gui
	guiApp := gui.NewGraphicApp(uc, conf)
	guiApp.Run(ctx)

	// i don't implement graceful shutdown because it's realized in gui
}
