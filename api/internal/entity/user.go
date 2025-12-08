package entity

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uuid.UUID `db:"id"`
	Email    string    `db:"email"`
	Name     string    `db:"name"`
	password []byte    `db:"password"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewUser(email, password, name string) (*User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       id,
		Email:    email,
		Name:     name,
		password: pass,
	}, nil
}

func (u *User) IsPasswordValid(reqPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.password), []byte(reqPassword))
	return err == nil
}

func (u *User) GetPassword() []byte {
	return u.password
}

func (u *User) SetPassword(password []byte) {
	u.password = password
}
