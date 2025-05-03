package kafka

import (
	"context"
	"web_app/settings"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

var (
	communityReader *kafka.Reader
	postReader      *kafka.Reader
	votePostReader  *kafka.Reader
)
var (
	VotePostWriter *kafka.Writer
)

const KeySendVotePostMessage = "vote_post"

// Init 初始化 Kafka 消费者
func Init(ctx context.Context, cfg *settings.KafkaConfig) {
	// 创建 Kafka 消费者
	communityReader = newKafkaReader(cfg.Brokers, cfg.GroupIDCommunity, cfg.TopicCommunity)
	postReader = newKafkaReader(cfg.Brokers, cfg.GroupIDPost, cfg.TopicPost)
	votePostReader = newKafkaReader(cfg.Brokers, cfg.GroupIDVotePost, cfg.TopicVotePost)
	// 创建 Kafka 生产者
	VotePostWriter = newKafkaWriter(cfg.Brokers, cfg.TopicVotePost)
	go func() {
		defer communityReader.Close() //确保消费者在Goroutine 退出时被关闭
		if err := readMessageToRedis(ctx, communityReader); err != nil {
			if ctx.Err() == context.Canceled {
				zap.L().Info("Kafka consumer stopped gracefully")
			} else {
				zap.L().Error("Kafka consumer encountered an error", zap.Error(err))
			}
		}
	}()
	go func() {
		defer postReader.Close()
		if err := readMessageToRedis(ctx, postReader); err != nil {
			if ctx.Err() == context.Canceled {
				zap.L().Info("Kafka consumer stopped gracefully")
			} else {
				zap.L().Error("Kafka consumer encountered an error", zap.Error(err))
			}
		}
	}()
	go func() {
		defer votePostReader.Close()
		if err := readMessageToMysql(ctx, votePostReader); err != nil {
			if ctx.Err() == context.Canceled {
				zap.L().Info("Kafka consumer stopped gracefully")
			} else {
				zap.L().Error("Kafka consumer encountered an error", zap.Error(err))
			}
		}
	}()
}
