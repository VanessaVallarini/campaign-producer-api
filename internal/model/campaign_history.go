package model

import (
	"time"

	"github.com/google/uuid"
)

type CampaignHistory struct {
	Id          uuid.UUID `json:"id"`
	CampaignId  uuid.UUID `json:"campaign_id"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}
