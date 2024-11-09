package producer

import (
	"context"

	"github.com/IBM/sarama"

	"time"

	"github.com/VanessaVallarini/campaign-producer-api/internal/config"
	"github.com/VanessaVallarini/campaign-producer-api/internal/pkg/kafka/client"
	easyzap "github.com/lockp111/go-easyzap"
)

const (
	Avro string = "AVRO"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
	srClient client.SchemaRegistryClient
	topic    string
	subject  string
}

func NewProducer(
	ctx context.Context,
	brokerConfig config.KafkaConfig,
	rawSchema string) KafkaProducer {
	_, _, saramaConfig := client.NewKafkaClient(brokerConfig)
	srClient := client.NewSchemaRegistry(brokerConfig)

	saramaConfig.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerConfig.Brokers, saramaConfig)
	if err != nil {

		easyzap.Fatal(ctx, err, "failed create producer")
	}

	if err := srClient.ValidateSchema(rawSchema, brokerConfig.Subject, Avro); err != nil {

		easyzap.Fatal(ctx, err, "failed validate schema")
	}

	return KafkaProducer{
		producer: producer,
		srClient: srClient,
		topic:    brokerConfig.Topic,
		subject:  brokerConfig.Subject,
	}
}

func (kp KafkaProducer) Send(key string, msg interface{}) error {
	msgEncoder, err := kp.srClient.Encode(msg, kp.subject)
	if err != nil {
		easyzap.Errorf("kafka producer failed encode msg %v: %v", msg, err)

		return err
	}

	m := sarama.ProducerMessage{
		Topic:     kp.topic,
		Key:       sarama.ByteEncoder(key),
		Value:     sarama.ByteEncoder(msgEncoder),
		Timestamp: time.Now(),
	}

	_, _, err = kp.producer.SendMessage(&m)
	if err != nil {

		return err
	}

	return nil
}
