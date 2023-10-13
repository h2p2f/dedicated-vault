// Package: middlewares
// in this file we have grpc client middlewares
package middlewares

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// JWTInjectorUnaryClientInterceptor is a middleware for injecting jwt token into metadata
func JWTInjectorUnaryClientInterceptor(token string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		var md metadata.MD
		// check if metadata exists
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			//if not, create new metadata
			md = metadata.New(nil)
		}
		md.Set("authorization", token)
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
