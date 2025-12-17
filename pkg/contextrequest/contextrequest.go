package contextrequest

import (
	"context"

	"github.com/elangreza/e-commerce/pkg/globalcontanta"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

func AppendUserIDintoContextGrpcClient(ctx context.Context, userID uuid.UUID) context.Context {
	md := metadata.New(map[string]string{string(globalcontanta.UserIDKey): userID.String()})
	return metadata.NewOutgoingContext(ctx, md)
}
