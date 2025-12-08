package rest

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/elangreza/e-commerce/api/internal/constanta"
	"github.com/google/uuid"
)

type (
	AuthService interface {
		ProcessToken(ctx context.Context, reqToken string) (uuid.UUID, error)
	}

	AuthMiddleware struct {
		svc AuthService
	}
)

func (am *AuthMiddleware) MustAuthMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			rawAuthorization := r.Header["Authorization"]
			if len(rawAuthorization) == 0 {
				sendErrorResponse(w, http.StatusBadRequest, errors.New("token not valid"))
				return
			}

			authorization := rawAuthorization[0]

			rawToken := strings.Split(authorization, " ")
			if len(rawToken) != 2 {
				sendErrorResponse(w, http.StatusBadRequest, errors.New("token not valid"))
				return
			}

			bearer := rawToken[0]
			if strings.ToLower(bearer) != "bearer" {
				sendErrorResponse(w, http.StatusBadRequest, errors.New("token not valid. must be bearer + token"))
				return
			}

			token := rawToken[1]

			userID, err := am.svc.ProcessToken(r.Context(), token)
			if err != nil {
				sendErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorize user"))
				return
			}

			ctx := context.WithValue(r.Context(), constanta.LocalUserID, userID)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
