package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	CampaignAvro = `{
		"type":"record",
		"name":"campaign",
		"namespace":"campaign.campaign_value",
		"fields":[
			{
				"name": "id",
				"type": {
				"type": "string",
				"logicalType": "UUID"
				}
			},
			{
				"name": "merchant_id",
				"type": {
				"type": "string",
				"logicalType": "UUID"
				}
			},
			{
				"name":"status",
				"type":"string"
			},
			{
				"name":"budget",
				"type":"double"
			},
			{
				"name":"created_by",
				"type":"string"
			},
			{
				"name":"updated_by",
				"type":"string"
			},
			{
				"name": "created_at",
				"type": {
				"type": "long",
				"logicalType": "timestamp-millis"
				}
			},
			{
				"name": "updated_at",
				"type": {
				"type": "long",
				"logicalType": "timestamp-millis"
				}
			}
		]
	 }`
)

type Campaign struct {
	Id         uuid.UUID `json:"id" avro:"id"`
	MerchantId uuid.UUID `json:"merchant_id" avro:"merchant_id"`
	Status     string    `json:"status" avro:"status"`
	Budget     float64   `json:"budget" avro:"budget"`
	CreatedBy  string    `json:"created_by" avro:"created_by"`
	UpdatedBy  string    `json:"updated_by" avro:"updated_by"`
	CreatedAt  time.Time `json:"created_at" avro:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" avro:"updated_at"`
}
