// Package: middlewares
package middlewares

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/h2p2f/dedicated-vault/internal/server/jwtprocessing"
)

// JWTCheckingUnaryServerInterceptor is an interceptor for checking jwt token
func JWTCheckingUnaryServerInterceptor(key string, fullAccessMethods map[string]bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		fmt.Println("interceptor")
		if fullAccessMethods[info.FullMethod] {
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "metadata not found")
		}
		authValues := md.Get("authorization")

		user, err := jwtprocessing.ParseToken(authValues[0], key)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		md.Set("user", user)
		ctx = metadata.NewIncomingContext(ctx, md)

		return handler(ctx, req)
	}
}
