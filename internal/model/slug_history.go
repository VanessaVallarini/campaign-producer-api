package model

import (
	"time"

	"github.com/google/uuid"
)

type SlugHistory struct {
	Id          uuid.UUID `json:"id"`
	SlugId      uuid.UUID `json:"slug_id"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}
