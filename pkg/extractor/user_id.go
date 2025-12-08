package extractor

import (
	"context"

	"github.com/elangreza/e-commerce/pkg/globalcontanta"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ExtractUserIDFromMetadata(ctx context.Context) (uuid.UUID, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return uuid.Nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}
	rawUserID := md.Get(string(globalcontanta.UserIDKey))

	if len(rawUserID) == 0 {
		return uuid.Nil, status.Errorf(codes.Unauthenticated, "not valid userID")
	}

	userID, err := uuid.Parse(rawUserID[0])
	if err != nil {
		return uuid.Nil, status.Errorf(codes.Unauthenticated, "failed to parse userID")
	}

	return userID, nil
}
