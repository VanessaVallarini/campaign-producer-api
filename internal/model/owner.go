package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	OwnerAvro = `{
		"type":"record",
		"name":"owner",
		"namespace":"campaign.campaign_owner_value",
		"fields":[
			{
				"name": "id",
				"type": {
				"type": "string",
				"logicalType": "UUID"
				}
			},
			{
				"name":"email",
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

type Owner struct {
	Id        uuid.UUID `json:"id" avro:"id"`
	Email     string    `json:"email" avro:"email"`
	Status    string    `json:"status" avro:"status"`
	CreatedBy string    `json:"created_by" avro:"created_by"`
	UpdatedBy string    `json:"updated_by" avro:"updated_by"`
	CreatedAt time.Time `json:"created_at" avro:"created_at"`
	UpdatedAt time.Time `json:"updated_at" avro:"updated_at"`
}

type OwnerCreateRequest struct {
	Email     string `json:"email"`
	CreatedBy string `json:"createdBy"`
}
