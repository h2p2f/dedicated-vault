// Package app
// configuring the server, logging, database, tls, grpc server
package app

import (
	"context"
	"log"
	"net"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/h2p2f/dedicated-vault/internal/server/config"
	"github.com/h2p2f/dedicated-vault/internal/server/grpcserver"
	"github.com/h2p2f/dedicated-vault/internal/server/grpcserver/middlewares"
	"github.com/h2p2f/dedicated-vault/internal/server/storage"
	"github.com/h2p2f/dedicated-vault/internal/server/tlsloader"
	pb "github.com/h2p2f/dedicated-vault/proto"
)

// Run starts the application
func Run(ctx context.Context, sigint chan os.Signal, connectionsClosed chan<- struct{}) {
	// read configuration
	conf := config.NewServerConfig()
	// create logger
	atom, err := zap.ParseAtomicLevel(conf.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.Lock(os.Stdout),
		atom))
	defer logger.Sync() //nolint:errcheck
	// create storage
	db := storage.NewStorage(ctx, conf, logger)
	// create grpc server
	// load tls
	tlsConf, err := tlsloader.LoadTLS(conf)
	if err != nil {
		logger.Fatal("tls", zap.Error(err))
	}
	tlsCredentials := credentials.NewTLS(tlsConf)
	opts := []grpc.ServerOption{
		grpc.Creds(tlsCredentials),
	}
	// add jwt middleware with unprotected methods
	unprotectedMethods := map[string]bool{
		"/DedicatedVault/Register": true,
		"/DedicatedVault/Login":    true,
	}
	opts = append(
		opts,
		grpc.UnaryInterceptor(
			middlewares.JWTCheckingUnaryServerInterceptor(conf.JWTKey, unprotectedMethods),
		))
	// create listener
	listener, err := net.Listen("tcp", ":8090")

	if err != nil {
		panic(err)
	}
	// create grpc server
	server := grpc.NewServer(opts...)

	vaultServer := grpcserver.NewVaultServer(db, db, logger)
	// register grpc server
	pb.RegisterDedicatedVaultServer(server, vaultServer)
	// run grpc server
	logger.Info("Starting server...",
		zap.String("address", conf.GRPCAddress),
		zap.String("tls cert", conf.ServerCert),
		zap.String("tls key", conf.ServerKey),
		zap.String("storage address", conf.StorageAddress),
		zap.String("log level", conf.LogLevel))
	go func() {
		if err := server.Serve(listener); err != nil {
			logger.Fatal("listen", zap.Error(err))
		}
	}()
	// wait for a signal to stop the server
	<-sigint
	logger.Info("Shutting down server...")
	server.GracefulStop()
	logger.Info("Server gracefully stopped")
	close(sigint)
	close(connectionsClosed)
}
