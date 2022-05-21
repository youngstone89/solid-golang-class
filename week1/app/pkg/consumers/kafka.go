package consumers

import (
	"context"
	"encoding/json"
	"event-data-pipeline/pkg/logger"
	"event-data-pipeline/pkg/sources"

	"event-data-pipeline/pkg/kafka"
)

// compile type assertion check
var _ Consumer = new(KafkaConsumerClient)
var _ ConsumerFactory = NewKafkaConsumerClient

// register kafka consumer client to factory
func init() {

	Register("kafka", NewKafkaConsumerClient)

}

type KafkaClientConfig struct {
	ClientName      string  `json:"client_name,omitempty"`
	Topic           string  `json:"topic,omitempty"`
	ConsumerOptions jsonObj `json:"consumer_options,omitempty"`
}

//implements Consumer interface
type KafkaConsumerClient struct {
	kafka.Consumer
	sources.Source
}

func NewKafkaConsumerClient(config jsonObj) Consumer {

	// Read config into KafkaClientConfig struct
	var kcCfg KafkaClientConfig
	_json, err := json.Marshal(config)
	if err != nil {
		logger.Panicf(err.Error())
	}
	json.Unmarshal(_json, &kcCfg)
	kfkCnsmr := kafka.NewKafkaConsumer(kcCfg.ConsumerOptions)

	// create a new Consumer concrete type - KafkaConsumerClient
	client := &KafkaConsumerClient{
		Consumer: kfkCnsmr,
		Source:   sources.NewKafkaSource(kfkCnsmr),
	}

	return client

}

// Init implements Consumer
func (kc *KafkaConsumerClient) Init() error {
	var err error

	err = kc.CreateConsumer()
	if err != nil {
		return err
	}
	err = kc.CreateAdminConsumer()
	if err != nil {
		return err
	}
	err = kc.GetPartitions()
	if err != nil {
		return err
	}
	return nil
}

// Consumer 인터페이스 구현
func (kc *KafkaConsumerClient) Consume(ctx context.Context) error {
	err := kc.Read(ctx)
	if err != nil {
		return err
	}
	return nil
}
