package client

import (
	"time"

	"github.com/IBM/sarama"
	"github.com/VanessaVallarini/campaign-producer-api/internal/config"
	easyzap "github.com/lockp111/go-easyzap"
)

const _saramaTimeoutFlushMs = 500

func NewKafkaClient(brokerConfig config.KafkaConfig) (sarama.Client, sarama.ClusterAdmin, *sarama.Config) {
	saramaConfig := generateSaramaConfig(brokerConfig)

	kafkaClient, err := sarama.NewClient(brokerConfig.Brokers, saramaConfig)
	if err != nil {
		easyzap.Fatal(err, "kafka client failed to new kafka client")
	}

	srClusterAdmin, err := sarama.NewClusterAdminFromClient(kafkaClient)
	if err != nil {
		easyzap.Fatal(err, "kafka client failed to new cluster admin")
	}

	if !brokerConfig.UseAuthentication {
		createTopic(brokerConfig.Topic, srClusterAdmin)
	}

	return kafkaClient, srClusterAdmin, saramaConfig
}

func generateSaramaConfig(brokerConfig config.KafkaConfig) *sarama.Config {
	srConfig := sarama.NewConfig()
	saramaTimeout := brokerConfig.Timeout * time.Millisecond

	srConfig.ClientID = brokerConfig.ClientId
	srConfig.Version = sarama.V3_0_0_0
	srConfig.Net.DialTimeout = saramaTimeout
	srConfig.Net.ReadTimeout = saramaTimeout
	srConfig.Net.WriteTimeout = saramaTimeout
	srConfig.Metadata.Timeout = saramaTimeout
	srConfig.Producer.RequiredAcks = sarama.WaitForLocal
	srConfig.Producer.Flush.Frequency = _saramaTimeoutFlushMs * time.Millisecond
	srConfig.Metadata.Retry.Max = brokerConfig.RetryMax
	srConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{brokerConfig.BalanceStrategy}

	if brokerConfig.UseAuthentication {
		srConfig.Net.SASL.Mechanism = sarama.SASLMechanism(brokerConfig.SaslMechanism)
		srConfig.Net.SASL.User = brokerConfig.User
		srConfig.Net.SASL.Password = brokerConfig.Password
		srConfig.Net.TLS.Enable = brokerConfig.EnableTLS
		srConfig.Net.SASL.Enable = true
		setAuthentication(srConfig)
	}

	return srConfig
}

func setAuthentication(srConfig *sarama.Config) {
	switch srConfig.Net.SASL.Mechanism {
	case sarama.SASLTypeSCRAMSHA512:
		scram512Fn := func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
		srConfig.Net.SASL.SCRAMClientGeneratorFunc = scram512Fn
	case sarama.SASLTypeSCRAMSHA256:
		scram256Fn := func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }
		srConfig.Net.SASL.SCRAMClientGeneratorFunc = scram256Fn
	}
}

func createTopic(topic string, srClusterAdmin sarama.ClusterAdmin) error {
	err := srClusterAdmin.CreateTopic(topic,
		&sarama.TopicDetail{
			NumPartitions:     4,
			ReplicationFactor: 1,
		},
		false)
	if err != nil {
		easyzap.Warnf("kafka client failed create topic %s: %v", topic, err)
	}
	return nil
}
