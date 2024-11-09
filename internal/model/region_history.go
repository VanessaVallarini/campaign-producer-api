package model

import (
	"time"

	"github.com/google/uuid"
)

type RegionHistory struct {
	Id          uuid.UUID `json:"id"`
	RegionId    uuid.UUID `json:"region_id"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}
