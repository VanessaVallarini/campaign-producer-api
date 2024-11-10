package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	RegionAvro = `{
		"type":"record",
		"name":"region",
		"namespace":"campaign.campaign_region_value",
		"fields":[
			{
				"name": "id",
				"type": {
				"type": "string",
				"logicalType": "UUID"
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
				"name":"lat",
				"type":"double"
			},
			{
				"name":"long",
				"type":"double"
			},
			{
				"name":"cost",
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

type Region struct {
	Id        uuid.UUID `json:"id" avro:"id"`
	Name      string    `json:"name" avro:"name"`
	Status    string    `json:"status" avro:"status"`
	Lat       float64   `json:"lat" avro:"lat"`
	Long      float64   `json:"long" avro:"long"`
	Cost      float64   `json:"cost" avro:"cost"`
	CreatedBy string    `json:"created_by" avro:"created_by"`
	UpdatedBy string    `json:"updated_by" avro:"updated_by"`
	CreatedAt time.Time `json:"created_at" avro:"created_at"`
	UpdatedAt time.Time `json:"updated_at" avro:"updated_at"`
}

type RegionCreateRequest struct {
	Name      string  `json:"name"`
	Lat       float64 `json:"lat"`
	Long      float64 `json:"long"`
	Cost      float64 `json:"cost"`
	CreatedBy string  `json:"createdBy"`
}
