package app

import (
	"context"
	"github.com/h2p2f/dedicated-vault/internal/server/grpcserver/middlewares"

	//"crypto/tls"
	"github.com/h2p2f/dedicated-vault/internal/server/config"
	"github.com/h2p2f/dedicated-vault/internal/server/grpcserver"
	"github.com/h2p2f/dedicated-vault/internal/server/storage"
	pb "github.com/h2p2f/dedicated-vault/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
)

func Run(ctx context.Context) {

	logger := zap.NewExample()

	conf := &config.ServerConfig{
		StorageAddress: "mongodb://localhost:27017",
		JWTKey:         "secret",
	}

	db := storage.NewStorage(ctx, conf, logger)

	//cert, err := tls.LoadX509KeyPair("./crypto/public.crt", "./crypto/private.key")
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//opts := []grpc.ServerOption{
	//	grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	//}
	var opts []grpc.ServerOption
	unprotectedMethods := map[string]bool{
		"/DedicatedVault/Register":       true,
		"/DedicatedVault/Login":          true,
		"/DedicatedVault/ChangePassword": true,
	}
	opts = append(opts, grpc.UnaryInterceptor(middlewares.JWTCheckingUnaryServerInterceptor(conf.JWTKey, unprotectedMethods)))
	listener, err := net.Listen("tcp", ":8090")

	if err != nil {
		panic(err)
	}

	server := grpc.NewServer(opts...)

	vaultServer := grpcserver.NewVaultServer(db, db, logger)

	pb.RegisterDedicatedVaultServer(server, vaultServer)

	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
	//go func() {
	//	logger.Info("starting server")
	//	if err := server.Serve(listener); err != nil {
	//		logger.Fatal("listen", zap.Error(err))
	//	}
	//}()

	//err = db.Close(ctx)
	//if err != nil {
	//	panic(err)
	//}

}
