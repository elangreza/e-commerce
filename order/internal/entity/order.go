package entity

import "github.com/google/uuid"

type Order struct {
	ID     uuid.UUID
	UserID string
	Cart   Cart
	Total  int64
	Status string
}
