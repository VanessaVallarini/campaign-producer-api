package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	MerchantAvro = `{
		"type":"record",
		"name":"merchant",
		"namespace":"campaign.campaign_merchant_value",
		"fields":[
			{
				"name": "id",
				"type": {
				"type": "string",
				"logicalType": "UUID"
				}
			},
			{
				"name": "owner_id",
				"type": {
				"type": "string",
				"logicalType": "UUID"
				}
			},
			{
				"name": "region_id",
				"type": {
				"type": "string",
				"logicalType": "UUID"
				}
			},
			{
				"name": "slugs",
				"type": {
				"type": "array",
				"items": {
					"type": "string",
					"logicalType": "UUID"
				}
				}
			},
			{
				"name":"name",
				"type":"string"
			},
			{
				"name":"status",
				"type":"string"
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

type Merchant struct {
	Id        uuid.UUID   `json:"id" avro:"id"`
	OwnerId   uuid.UUID   `json:"owner_id" avro:"owner_id"`
	RegionId  uuid.UUID   `json:"region_id" avro:"region_id"`
	Slugs     []uuid.UUID `json:"slugs" avro:"slugs"`
	Name      string      `json:"name" avro:"name"`
	Status    string      `json:"status" avro:"status"`
	CreatedBy string      `json:"created_by" avro:"created_by"`
	UpdatedBy string      `json:"updated_by" avro:"updated_by"`
	CreatedAt time.Time   `json:"created_at" avro:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" avro:"updated_at"`
}
