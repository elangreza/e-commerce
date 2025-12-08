package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/elangreza/e-commerce/api/internal/constanta"
	"github.com/elangreza/e-commerce/api/internal/entity"
	errs "github.com/elangreza/e-commerce/api/internal/error"
	"github.com/elangreza/e-commerce/api/internal/params"
	"github.com/google/uuid"
)

type (
	userRepo interface {
		CreateUser(ctx context.Context, user entity.User) error
		GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
		GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	}

	tokenRepo interface {
		CreateToken(ctx context.Context, token entity.Token) error
		GetTokenByUserID(ctx context.Context, userID uuid.UUID) (*entity.Token, error)
		GetTokenByTokenID(ctx context.Context, tokenID uuid.UUID) (*entity.Token, error)
	}

	AuthService struct {
		UserRepo  userRepo
		TokenRepo tokenRepo
	}
)

func NewAuthService(userRepo userRepo, tokenRepo tokenRepo) *AuthService {
	return &AuthService{
		UserRepo:  userRepo,
		TokenRepo: tokenRepo,
	}
}

func (as *AuthService) RegisterUser(ctx context.Context, req params.RegisterUserRequest) error {
	user, err := as.UserRepo.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if user != nil {
		return errs.AlreadyExist{Name: fmt.Sprintf("email %s", req.Email)}
	}

	user, err = entity.NewUser(req.Email, req.Password, req.Name)
	if err != nil {
		return err
	}

	err = as.UserRepo.CreateUser(ctx, *user)
	if err != nil {
		return err
	}

	return nil
}

func (as *AuthService) LoginUser(ctx context.Context, req params.LoginUserRequest) (string, error) {
	user, err := as.UserRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.NotFound{Message: fmt.Sprintf("email %s", req.Email)}
		}
		return "", err
	}

	ok := user.IsPasswordValid(req.Password)
	if !ok {
		return "", errs.InvalidCredential{}
	}

	token, err := as.TokenRepo.GetTokenByUserID(ctx, user.ID)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	if token != nil {
		_, err = token.IsTokenValid([]byte(constanta.AuthenticationSigningKey))
		if err == nil {
			return token.Token, nil
		}
	}

	token, err = entity.NewToken([]byte(constanta.AuthenticationSigningKey), user.ID, "LOGIN")
	if err != nil {
		return "", err
	}

	err = as.TokenRepo.CreateToken(ctx, *token)
	if err != nil {
		return "", err
	}

	return token.Token, nil
}

func (as *AuthService) ProcessToken(ctx context.Context, reqToken string) (uuid.UUID, error) {
	token := &entity.Token{Token: reqToken}

	tokenID, err := token.IsTokenValid([]byte(constanta.AuthenticationSigningKey))
	if err != nil {
		return uuid.UUID{}, err
	}

	token, err = as.TokenRepo.GetTokenByTokenID(ctx, tokenID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.UUID{}, errs.NotFound{Message: "token"}
		}
		return uuid.UUID{}, err
	}

	return token.UserID, nil
}
