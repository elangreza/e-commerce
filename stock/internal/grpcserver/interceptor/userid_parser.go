package interceptor

import (
	"context"
	"github/elangreza/e-commerce/stock/internal/constanta"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func UserIDParser() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		userID, ok := ctx.Value(constanta.UserIDKey).(string)

		if ok {
			uid, err := uuid.Parse(userID) // just to validate
			if err != nil {
				return nil, err
			}
			ctx = context.WithValue(ctx, constanta.UserIDKey, uid)
		}

		return handler(ctx, req)
	}
}
