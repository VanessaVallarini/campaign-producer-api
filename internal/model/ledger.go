package model

import (
	"time"

	"github.com/google/uuid"
)

type Ledger struct {
	Id         uuid.UUID `json:"id"`
	SpentId    uuid.UUID `json:"spent_id"`
	CampaignId uuid.UUID `json:"campaign_id"`
	MerchantId uuid.UUID `json:"merchant_id"`
	SlugName   string    `json:"slug_name"`
	RegionName string    `json:"region_name"`
	UserId     uuid.UUID `json:"user_id"`
	SessionId  uuid.UUID `json:"session_id"`
	EventType  EventType `json:"event_type"`
	Cost       float64   `json:"cost"`
	Ip         string    `json:"ip"`
	Lat        float64   `json:"lat"`
	Long       float64   `json:"long"`
	CreatedAt  time.Time `json:"created_at"`
	EventTime  time.Time `json:"event_time"`
}
