package interceptor

import (
	"context"
	"github/elangreza/e-commerce/stock/internal/constanta"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func UserIDParser() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		userID := ctx.Value(constanta.UserIDKey).(uuid.UUID)
		ctx = context.WithValue(ctx, constanta.UserIDKey, userID)
		return handler(ctx, req)
	}
}
