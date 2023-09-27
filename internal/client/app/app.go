package app

import (
	"context"
	"crypto/sha256"
	"github.com/h2p2f/dedicated-vault/internal/client/config"
	"github.com/h2p2f/dedicated-vault/internal/client/grpcclient"
	"github.com/h2p2f/dedicated-vault/internal/client/gui"
	"github.com/h2p2f/dedicated-vault/internal/client/storage"
	"github.com/h2p2f/dedicated-vault/internal/client/tlsloader"
	"github.com/h2p2f/dedicated-vault/internal/client/usecase"
	"go.uber.org/zap"
)

func Run(ctx context.Context) {
	var err error

	conf := config.NewClientConfig()

	Key := sha256.Sum256([]byte(conf.Passphrase))

	conf.CryptoKey = Key[:]

	logger := zap.NewExample()

	db := storage.NewClientStorage(logger, conf)

	conf.TLSConfig, err = tlsloader.LoadTLS()
	if err != nil {
		logger.Fatal("tls", zap.Error(err))
	}

	tr := grpcclient.NewClient(conf, logger)
	uc := usecase.NewClientUseCase(conf, db, tr)

	guiApp := gui.NewGraphicApp(uc, conf)
	guiApp.Run(ctx)

}
