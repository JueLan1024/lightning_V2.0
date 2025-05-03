package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"web_app/dao/mysql"
	"web_app/models"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// newKafkaReader 获取新消费者
func newKafkaReader(brokers []string, groupID string, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        groupID, // 指定消费者组id
		Topic:          topic,
		CommitInterval: 0, // 禁用自动提交偏移量
	})
}

// ReadCanalMessage 读取从Canal发送给Kafka的消息
func ReadCanalMessage(ctx context.Context, r *kafka.Reader) (message *models.CanalMessage, m kafka.Message, err error) {
	m, err = r.ReadMessage(ctx)
	if errors.Is(err, context.Canceled) {
		return nil, kafka.Message{}, nil
	}
	if err != nil {
		zap.L().Error("r.ReadMessage failed", zap.Error(err))
		return nil, kafka.Message{}, err
	}
	// 将数据写入message结构体
	message = new(models.CanalMessage)
	if err = json.Unmarshal(m.Value, message); err != nil {
		zap.L().Error("failed to unmarshal msg from kafka", zap.Error(err))
		return nil, kafka.Message{}, err
	}
	return message, m, nil

}

// readMessageToRedis 从kafka中读取新增的消息写入Redis
func readMessageToRedis(ctx context.Context, r *kafka.Reader) (err error) {
	for {
		select {
		case <-ctx.Done(): //检查上下文是否已取消
			zap.L().Info("readMessageToRedis stopped")
			return ctx.Err()
		default:
			// 获取消息
			msg, m, err := ReadCanalMessage(ctx, r)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return nil //正常退出
				}
				zap.L().Error("ReadCanalMessage failed", zap.Error(err))
				continue
			}
			// 将消息写入Redis
			if msg.Type == "INSERT" {
				if r == communityReader {
					if err = insertCommunityInRedis(ctx, msg.Data[0]); err != nil {
						zap.L().Error("insert community in redis failed", zap.Error(err))
						continue // 如果 Redis 写入失败，不提交偏移量
					}
				}
				if r == postReader {
					if err = insertPostInRedis(ctx, msg.Data[0]); err != nil {
						zap.L().Error("insertInRedis failed", zap.Error(err))
						continue // 如果 Redis 写入失败，不提交偏移量
					}
				}
				// 提交 Kafka 消息的偏移量
				if err = r.CommitMessages(ctx, m); err != nil {
					zap.L().Error("kafka commitMessages failed", zap.Error(err))
					continue
				}

			}
		}
	}
}

// readMessageToMysql 将消息存入mysql
func readMessageToMysql(ctx context.Context, r *kafka.Reader) (err error) {
	for {
		select {
		case <-ctx.Done():
			zap.L().Info("readMessageToMysql stoped")
			return ctx.Err()
		default:
			// 读取消息
			m, err := r.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return nil //正常退出
				}
				zap.L().Error("r.ReadMessage failed", zap.Error(err))
				continue
			}
			// 将消息提取到结构体中
			if string(m.Key) == KeySendVotePostMessage {
				// 解析数据
				msg := new(models.VotePost)
				if err = json.Unmarshal(m.Value, msg); err != nil {
					zap.L().Error("json.Unmarshal failed", zap.Error(err))
					continue
				}
				// 查询数据库是否已经有投票记录
				exist, err := mysql.VotePostExist(ctx, msg)
				if err != nil {
					zap.L().Error("mysql.VotePostExist failed",
						zap.Int64("post_id", msg.PostID),
						zap.Int64("user_id", msg.UserID),
						zap.Int8("vote_type", msg.VoteType),
						zap.Error(err),
					)
					continue
				}
				if !exist { //没有投票记录
					// 将消息存入数据库
					if err = mysql.CreateVotePost(ctx, msg); err != nil {
						zap.L().Error("mysql.CreateVotePost failed",
							zap.Int64("post_id", msg.PostID),
							zap.Int64("user_id", msg.UserID),
							zap.Int8("vote_type", msg.VoteType),
							zap.Error(err),
						)
						continue
					}
				} else { //有投票记录
					// 更新数据库
					if err = mysql.UpdateVotePost(ctx, msg); err != nil {
						zap.L().Error("mysql.UpdateVotePost failed",
							zap.Int64("post_id", msg.PostID),
							zap.Int64("user_id", msg.UserID),
							zap.Int8("vote_type", msg.VoteType),
							zap.Error(err),
						)
						continue
					}
				}
				// 提交 Kafka 消息的偏移量
				if err = r.CommitMessages(ctx, m); err != nil {
					zap.L().Error("kafka commitMessages failed", zap.Error(err))
					continue
				}
			}
		}
	}
}
