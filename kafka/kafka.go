package kafka

import (
	"github.com/IBM/sarama"
)

type KafkaClient struct {
	Producer sarama.SyncProducer
	Consumer sarama.Consumer
}

func NewKafkaClient(brokers []string) (*KafkaClient, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, err
	}

	return &KafkaClient{
		Producer: producer,
		Consumer: consumer,
	}, nil
}

func (k *KafkaClient) SendMessage(topic string, message []byte) error {
	_, _, err := k.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	})
	return err
}

func (k *KafkaClient) ConsumeMessages(topic string, partition int32, offset int64) (chan *sarama.ConsumerMessage, error) {
	partitionConsumer, err := k.Consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		return nil, err
	}

	messages := make(chan *sarama.ConsumerMessage)
	go func() {
		for msg := range partitionConsumer.Messages() {
			messages <- msg
		}
		close(messages)
	}()

	return messages, nil
}

func (k *KafkaClient) Close() error {
	if err := k.Producer.Close(); err != nil {
		return err
	}
	return k.Consumer.Close()
}
