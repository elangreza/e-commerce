package interceptor

import (
	"context"
	"github/elangreza/e-commerce/pkg/globalcontanta"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UserIDParser() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			// Get the userID from metadata
			values := md.Get(string(globalcontanta.UserIDKey))
			if len(values) > 0 {
				userID := values[0]
				uid, err := uuid.Parse(userID)
				if err != nil {
					return nil, err
				}
				ctx = context.WithValue(ctx, globalcontanta.UserIDKey, uid)
			}
		}

		return handler(ctx, req)
	}
}
