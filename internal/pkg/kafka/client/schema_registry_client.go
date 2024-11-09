package client

import (
	"encoding/binary"
	"errors"
	"strings"

	"github.com/VanessaVallarini/campaign-producer-api/internal/config"
	"github.com/hamba/avro"
	easyzap "github.com/lockp111/go-easyzap"
	"github.com/riferrei/srclient"
)

type SchemaRegistryClient struct {
	client *srclient.SchemaRegistryClient
}

func NewSchemaRegistry(brokerConfig config.KafkaConfig) SchemaRegistryClient {
	srClient := srclient.CreateSchemaRegistryClient(brokerConfig.SchemaRegistryConfig.Host)
	if brokerConfig.UseAuthentication {
		srClient.SetCredentials(brokerConfig.SchemaRegistryConfig.User, brokerConfig.SchemaRegistryConfig.Password)
	}

	return SchemaRegistryClient{
		client: srClient,
	}
}

// ValidateSchema checks for the existence and compatibility of a schema.
// if the subject does not exist it will be created, if it is incompatible it will return an error.
func (sr SchemaRegistryClient) ValidateSchema(rawSchema, subject string, schemaType string) error {
	schema, err := sr.client.GetLatestSchema(subject)

	if err != nil && !strings.Contains(err.Error(), "404") {
		easyzap.Error(err, "kafka schema registry")

		return err
	}

	if schema == nil {
		_, err := sr.client.CreateSchema(subject, rawSchema, srclient.SchemaType(schemaType))
		if err != nil {

			return err
		}

		return nil
	}

	isCompatible, err := sr.client.IsSchemaCompatible(subject, rawSchema, "latest", srclient.SchemaType(schemaType))
	if err != nil || !isCompatible {

		return err
	}

	if !isCompatible {

		return errors.New("schema registry invalid schema is not compatible")
	}

	return nil
}

func (sr SchemaRegistryClient) GetSchema(subject string) (*srclient.Schema, error) {
	schema, err := sr.client.GetLatestSchema(subject)
	if err != nil {

		return nil, err
	}

	if schema == nil {

		return nil, errors.New("schema registry unexpected behavior retrieving schema, got 'nil' from registry")
	}

	return schema, nil
}

func (sr SchemaRegistryClient) Encode(value interface{}, subject string) ([]byte, error) {
	schema, err := sr.GetSchema(subject)
	if err != nil {

		return nil, err
	}

	schemaEncoder, err := avro.Parse(schema.Schema())
	if err != nil {

		return []byte{}, err
	}

	avroNative, err := avro.Marshal(schemaEncoder, value)
	if err != nil {

		return []byte{}, err
	}

	var recordValue []byte
	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schema.ID()))
	recordValue = append(recordValue, byte(0))
	recordValue = append(recordValue, schemaIDBytes...)
	recordValue = append(recordValue, avroNative...)

	return recordValue, nil
}
