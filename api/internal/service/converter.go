package service

import (
	errs "github.com/elangreza/e-commerce/api/internal/error"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func convertErrGrpc(err error) error {
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.InvalidArgument:
			return errs.ValidationError{
				Message: st.Message(),
			}
		case codes.NotFound:
			return errs.NotFound{
				Message: st.Message(),
			}
		case codes.Unauthenticated:
			return errs.InvalidCredential{}
		}
	}

	return err
}
