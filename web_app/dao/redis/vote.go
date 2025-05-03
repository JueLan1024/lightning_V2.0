package redis

import (
	"context"
	"strconv"
	"web_app/models"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// GetVoteType 获得当前帖子当前用户的旧投票类型
func GetVoteType(ctx context.Context, votePost *models.VotePost) (oVoteType int8, err error) {
	key := GetKeyVotePostHash(votePost.PostID)
	oVoteTypeStr, err := rdb.HGet(ctx, key, strconv.FormatInt(votePost.UserID, 10)).Result()
	if err == redis.Nil { // 用户没有给当前帖子投票
		oVoteType = 0
		return oVoteType, nil
	}
	if err != nil { //查询出错
		zap.L().Error("rdb.HGet Get voteType failed", zap.Error(err))
		return 0, err
	}
	// 返回投票数据
	parsedValue, err := strconv.ParseInt(oVoteTypeStr, 10, 8)
	if err != nil {
		zap.L().Error("strconv.ParseInt failed", zap.Error(err))
		return 0, err
	}
	oVoteType = int8(parsedValue)
	return oVoteType, nil
}

// VoteForPost 将投票数据存入redis
func VoteForPost(ctx context.Context, changeScore int, votePost *models.VotePost) (err error) {
	key := GetKeyVotePostHash(votePost.PostID)
	// 将用户名和投票类型存入redis
	txPipe := rdb.TxPipeline()
	txPipe.HSet(ctx, key, votePost.UserID, votePost.VoteType)
	// 更新帖子分数
	key = GetKeyPostScoreZSet()
	txPipe.ZIncrBy(ctx, key, float64(changeScore), strconv.FormatInt(votePost.PostID, 10))
	_, err = txPipe.Exec(ctx)
	return err
}
