package middlewares

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/h2p2f/dedicated-vault/internal/server/jwtprocessing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

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
		tk := jwtprocessing.Claims{}

		token, err := jwt.ParseWithClaims(authValues[0], &tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}
		if !token.Valid {
			return nil, status.Errorf(codes.Unauthenticated, "token's param is invalid")
		}
		fmt.Println("token valid", tk.Login)
		md.Set("user", tk.Login)
		ctx = metadata.NewIncomingContext(ctx, md)
		//ctx = context.WithValue(ctx, "user", tk.Login)

		return handler(ctx, req)
	}
}
