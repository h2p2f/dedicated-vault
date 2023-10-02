package app

import (
	"context"
	"github.com/h2p2f/dedicated-vault/internal/server/config"
	"github.com/h2p2f/dedicated-vault/internal/server/grpcserver"
	"github.com/h2p2f/dedicated-vault/internal/server/grpcserver/middlewares"
	"github.com/h2p2f/dedicated-vault/internal/server/storage"
	"github.com/h2p2f/dedicated-vault/internal/server/tlsloader"
	pb "github.com/h2p2f/dedicated-vault/proto"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"os"
)

func Run(ctx context.Context, sigint chan os.Signal, connectionsClosed chan<- struct{}) {

	conf := config.NewServerConfig()

	atom, err := zap.ParseAtomicLevel(conf.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.Lock(os.Stdout),
		atom))

	defer logger.Sync() //nolint:errcheck

	db := storage.NewStorage(ctx, conf, logger)

	tlsConf, err := tlsloader.LoadTLS()
	if err != nil {
		logger.Fatal("tls", zap.Error(err))
	}
	tlsCredentials := credentials.NewTLS(tlsConf)
	opts := []grpc.ServerOption{
		grpc.Creds(tlsCredentials),
	}

	unprotectedMethods := map[string]bool{
		"/DedicatedVault/Register": true,
		"/DedicatedVault/Login":    true,
	}
	opts = append(
		opts,
		grpc.UnaryInterceptor(
			middlewares.JWTCheckingUnaryServerInterceptor(conf.JWTKey, unprotectedMethods),
		))
	listener, err := net.Listen("tcp", ":8090")

	if err != nil {
		panic(err)
	}

	server := grpc.NewServer(opts...)

	vaultServer := grpcserver.NewVaultServer(db, db, logger)

	pb.RegisterDedicatedVaultServer(server, vaultServer)

	go func() {
		if err := server.Serve(listener); err != nil {
			logger.Fatal("listen", zap.Error(err))
		}
	}()

	<-sigint
	logger.Info("Shutting down server...")
	server.GracefulStop()
	logger.Info("Server gracefully stopped")
	close(sigint)
	close(connectionsClosed)
}
