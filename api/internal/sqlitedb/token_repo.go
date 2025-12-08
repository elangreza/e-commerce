package sqlitedb

import (
	"context"
	"database/sql"
	"time"

	"github.com/elangreza/e-commerce/api/internal/entity"
	"github.com/google/uuid"
)

type (
	TokenRepo struct {
		db *sql.DB
	}
)

func NewTokenRepo(db *sql.DB) *TokenRepo {
	return &TokenRepo{
		db: db,
	}
}

const (
	createTokenQuery = `INSERT INTO tokens
(id, user_id, "token", token_type, issued_at, expired_at, duration)
VALUES(?, ?, ?, ?, ?, ?, ?);`
)

// CreateToken implements tokenRepo.
func (u *TokenRepo) CreateToken(ctx context.Context, token entity.Token) error {
	_, err := u.db.ExecContext(ctx, createTokenQuery,
		token.ID,
		token.UserID,
		token.Token,
		token.TokenType,
		token.IssuedAt,
		token.ExpiredAt,
		token.Duration,
	)
	if err != nil {
		return err
	}

	return nil
}

const (
	getTokenByTokenIDQuery = `SELECT 
		id, 
		user_id, 
		"token", 
		token_type, 
		issued_at, 
		expired_at, 
		duration
	FROM tokens
	WHERE id = ?
	;`
)

// GetTokenByTokenID implements tokenRepo.
func (u *TokenRepo) GetTokenByTokenID(ctx context.Context, tokenID uuid.UUID) (*entity.Token, error) {
	token := &entity.Token{}
	err := u.db.QueryRowContext(ctx, getTokenByTokenIDQuery, tokenID).Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.TokenType,
		&token.IssuedAt,
		&token.ExpiredAt,
		&token.Duration,
	)
	if err != nil {
		return nil, err
	}

	return token, nil
}

const (
	getTokenByUserIDQuery = `SELECT 
		id, 
		user_id, 
		"token", 
		token_type, 
		issued_at, 
		expired_at, 
		duration
	FROM tokens
	WHERE user_id = ? AND expired_at > ?
	;`
)

// GetTokenByUserID implements tokenRepo.
func (u *TokenRepo) GetTokenByUserID(ctx context.Context, userID uuid.UUID) (*entity.Token, error) {
	token := &entity.Token{}
	err := u.db.QueryRowContext(ctx, getTokenByUserIDQuery, userID, time.Now()).Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.TokenType,
		&token.IssuedAt,
		&token.ExpiredAt,
		&token.Duration,
	)
	if err != nil {
		return nil, err
	}

	return token, nil
}
