package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// newKafkaWriter 创建 Kafka Writer
func newKafkaWriter(brokers []string, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{}, // 使用最小字节负载均衡策略
	})
}

// SendMessage 发送消息到 Kafka
func SendMessage(ctx context.Context, w *kafka.Writer, key, value string) error {
	err := w.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	})
	if err != nil {
		zap.L().Error("w.WriteMessages failed",
			zap.String("key", key),
			zap.String("value", value),
			zap.Error(err),
		)
	}
	zap.L().Info("message sent successfully",
		zap.String("key", key),
		zap.String("value", value),
	)
	return nil
}
