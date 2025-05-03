package logic

import (
	"context"
	"encoding/json"
	"web_app/dao/redis"
	"web_app/kafka"
	"web_app/models"

	"go.uber.org/zap"
)

const (
	scorePerVote = 432 //每票价值432分 86400/200  --> 200张赞成票可以给你的帖子续一天
)

// VoteForPost 帖子投票业务
func VoteForPost(ctx context.Context, userID int64, p *models.ParamVoteForPost) (err error) {
	votePost := &models.VotePost{
		PostID:   p.PostID,
		UserID:   userID,
		VoteType: p.VoteType,
	}
	// 将投票数据序列化为json格式
	data, err := json.Marshal(votePost)
	if err != nil {
		zap.L().Error("json.Marshal failed",
			zap.Int64("user_id", votePost.UserID),
			zap.Int64("post_id", votePost.PostID),
			zap.Error(err),
		)
		return err
	}
	// 获得当前帖子下的当前用户投票类型
	oVoteType, err := redis.GetVoteType(ctx, votePost)
	if err != nil {
		zap.L().Error("redis.GetVoteType failed",
			zap.Int64("user_id", votePost.UserID),
			zap.Int64("post_id", votePost.PostID),
			zap.Error(err),
		)
		return err
	}
	// 前后投票类型一致返回重复投票错误
	if oVoteType == votePost.VoteType {
		return ErrorVoteRepeated
	}

	// 计算新旧投票类型的差值
	diff := votePost.VoteType - oVoteType
	changeScore := int(diff) * scorePerVote
	// 将投票数据存入redis
	if err = redis.VoteForPost(ctx, changeScore, votePost); err != nil {
		zap.L().Error("redis.VoteForPost failed",
			zap.Int64("user_id", votePost.UserID),
			zap.Int64("post_id", votePost.PostID),
			zap.Error(err),
		)
		return err
	}

	// 将投票数据传给kafka
	if err = kafka.SendMessage(ctx, kafka.VotePostWriter, kafka.KeySendVotePostMessage, string(data)); err != nil {
		zap.L().Error("kafka.SendMessage failed",
			zap.Int64("user_id", votePost.UserID),
			zap.Int64("post_id", votePost.PostID),
			zap.Error(err),
		)
		return err
	}
	return nil
}
