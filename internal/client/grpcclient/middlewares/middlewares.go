package middlewares

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func JWTInjectorUnaryClientInterceptor(token string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		fmt.Println(method)
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