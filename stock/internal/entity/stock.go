package entity

import "github.com/google/uuid"

type Stock struct {
	ProductID uuid.UUID `json:"id"`
	Quantity  int32     `json:"quantity"`
}
