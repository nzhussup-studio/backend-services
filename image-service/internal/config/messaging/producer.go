package messaging

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type ProducerInterface interface {
	SendMessage(message []byte) error
	Close()
}

type Producer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokerList []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokerList...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) SendMessage(message []byte) error {
	err := p.writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: message,
		},
	)
	if err != nil {
		log.Printf("failed to write message to Kafka: %v", err)
		return err
	}
	return nil
}

func (p *Producer) Close() {
	if err := p.writer.Close(); err != nil {
		log.Fatalf("failed to close Kafka producer: %v", err)
	}
}
