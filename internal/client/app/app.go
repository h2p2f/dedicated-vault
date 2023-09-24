package app

import (
	"context"
	"crypto/sha256"

	//"encoding/pem"
	//"fmt"
	"github.com/h2p2f/dedicated-vault/internal/client/config"
	"github.com/h2p2f/dedicated-vault/internal/client/grpcclient"
	"github.com/h2p2f/dedicated-vault/internal/client/gui"
	"github.com/h2p2f/dedicated-vault/internal/client/storage"
	"github.com/h2p2f/dedicated-vault/internal/client/usecase"
	"go.uber.org/zap"
	//"os"
)

func Run(ctx context.Context) {
	conf := config.NewClientConfig()

	Key := sha256.Sum256([]byte(conf.Passphrase))

	conf.CryptoKey = Key[:]

	//fmt.Println(conf.CryptoKey)

	logger := zap.NewExample()
	db := storage.NewClientStorage(logger, conf)

	//cert, err := os.ReadFile("./crypto/public.crt")
	//if err != nil {
	//	panic(err)
	//}
	//block, _ := pem.Decode(cert)
	//conf.Cert.AppendCertsFromPEM(block.Bytes)
	//
	tr := grpcclient.NewClient(conf, logger)
	uc := usecase.NewClientUseCase(conf, db, tr)
	//
	//err = uc.CreateUser("test12", "test12", "test12")
	//if err != nil {
	//	fmt.Println(err)
	//}

	guiApp := gui.NewGraphicApp(uc, conf)
	guiApp.Run(ctx)

}
