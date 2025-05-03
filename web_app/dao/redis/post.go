package redis

import (
	"context"
	"strconv"
	"time"
	"web_app/models"
	"web_app/tool"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const (
	NewPostExpireTime           = 24 * time.Hour   // 新帖子过期时间设为一天
	OldPostExpireTime           = 24 * time.Hour   // 旧帖子过期时间设为一天
	CommunityPostListExpireTime = 60 * time.Second // 有序社区帖子列表过期时间设为一分钟
)

// CreatePost 创建帖子
func CreatePost(ctx context.Context, post *models.Post) (err error) {
	postMap := map[string]interface{}{
		"post_id":      post.PostID,
		"title":        post.Title,
		"content":      post.Content,
		"author_id":    post.AuthorID,
		"community_id": post.CommunityID,
		"create_time":  post.CreatTime,
	}
	// 将帖子信息存入 lightning:post:<post_id> Hash
	key := GetKeyPostHash(post.PostID)
	txPipe := rdb.TxPipeline()
	txPipe.HSet(ctx, key, postMap)
	txPipe.Expire(ctx, key, NewPostExpireTime) // 给新创建的帖子设置过期时间
	// 将帖子id和帖子创建时间存入 lightning:post:time ZSet
	key = GetKeyPostTimeZSet()
	createTimeUnix := post.CreatTime.Unix()
	txPipe.ZAdd(ctx, key, &redis.Z{
		Score:  float64(createTimeUnix),
		Member: post.PostID,
	})
	// 将帖子id和帖子分数存入 lightning:post:score ZSet
	key = GetKeyPostScoreZSet()
	txPipe.ZAdd(ctx, key, &redis.Z{
		Score:  float64(createTimeUnix),
		Member: post.PostID,
	})
	// 将帖子存入其社区 lightning:community:<community_id>:posts Set
	key = GetKeyCommunityPostsSet(post.CommunityID)
	txPipe.SAdd(ctx, key, post.PostID)
	_, err = txPipe.Exec(ctx)
	return err
}

// InsertPost 旧帖子设置到缓存中
func InsertPost(ctx context.Context, post *models.Post) (err error) {
	postMap := map[string]interface{}{
		"post_id":      post.PostID,
		"title":        post.Title,
		"content":      post.Content,
		"author_id":    post.AuthorID,
		"community_id": post.CommunityID,
		"create_time":  post.CreatTime,
	}
	// 将帖子信息存入 lightning:post:<post_id> Hash
	key := GetKeyPostHash(post.PostID)
	txPipe := rdb.TxPipeline()
	txPipe.HSet(ctx, key, postMap)
	txPipe.Expire(ctx, key, OldPostExpireTime) // 给写入缓存的旧帖子设置过期时间
	_, err = txPipe.Exec(ctx)
	return err
}

// GetPost 获取帖子信息
func GetPost(ctx context.Context, key string) (post *models.Post, err error) {
	data, err := rdb.HGetAll(ctx, key).Result()
	if len(data) == 0 { //没有数据返回错误
		return nil, ErrorDataNotFound
	}
	if err != nil {
		zap.L().Error("Get post failed", zap.Error(err))
		return nil, err
	}
	postIDStr := data["post_id"]
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		zap.L().Error("postIDStr strconv.ParseInt failed", zap.Error(err))
		return nil, err
	}
	communityIDStr := data["community_id"]
	communityID, err := strconv.ParseInt(communityIDStr, 10, 64)
	if err != nil {
		zap.L().Error("communityIDStr strconv.ParseInt failed", zap.Error(err))
		return nil, err
	}
	authorIDStr := data["author_id"]
	authorID, err := strconv.ParseInt(authorIDStr, 10, 64)
	if err != nil {
		zap.L().Error("authorIDStr strconv.ParseInt failed", zap.Error(err))
		return nil, err
	}
	createTimeStr := data["create_time"]
	createTime, err := tool.ParseTime(createTimeStr)
	if err != nil {
		zap.L().Error("createTimeStr time.Parse failed", zap.Error(err))
		return nil, err
	}
	post = &models.Post{
		PostID:      postID,
		AuthorID:    authorID,
		CommunityID: communityID,
		Title:       data["title"],
		Content:     data["content"],
		CreatTime:   createTime,
	}
	return post, nil
}

// GetVoteNum 通过帖子id获取投票数据
func GetVoteNum(ctx context.Context, postID int64) (voteNum int64, err error) {
	key := GetKeyVotePostHash(postID)
	results, err := rdb.HGetAll(ctx, key).Result()
	if err == redis.Nil {
		zap.L().Warn("vote data not found ",
			zap.Int64("post_id", postID),
		)
		return 0, nil
	}
	if err != nil {
		zap.L().Error("rdb.HGetAll failed",
			zap.Int64("post_id", postID),
			zap.Error(err),
		)
		return 0, err
	}
	voteNum = 0
	for _, valueStr := range results {
		value, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			zap.L().Error("strconv.ParseInt failed")
			continue
		}
		voteNum += value
	}
	return voteNum, nil
}

// GetPostList 获得帖子列表，返回帖子ID的字符串切片
func GetPostList(ctx context.Context, p *models.ParamGetPostsInOrder, cursor string, pageSize int64) (postIDStrs []string, err error) {
	// 将保存社区内的所有帖子的Set和保存所有帖子的ZSet合并
	ckey := GetKeyCommunityPostsSet(p.CommunityID)
	orderKey := GetKeyPostScoreZSet() //默认帖子按照分数排序
	key := GetKeyCommunityPostScoreZSet(p.CommunityID)
	if p.Order == "time" {
		orderKey = GetKeyPostTimeZSet()
		key = GetKeyCommunityPostTimeZSet(p.CommunityID)
	}
	// 如果合并后的 ZSet 不存在，创建
	if rdb.Exists(ctx, key).Val() < 1 {
		pipe := rdb.Pipeline()
		pipe.ZInterStore(ctx, key, &redis.ZStore{
			Keys:      []string{ckey, orderKey},
			Aggregate: "MAX",
		})
		pipe.Expire(ctx, key, CommunityPostListExpireTime)
		_, err = pipe.Exec(ctx)
		if err != nil {
			zap.L().Error("pipe.ZInterStore failed",
				zap.Int64("community_id", p.CommunityID),
				zap.Error(err),
			)
			return nil, err
		}
	}
	// 设置起始索引
	var start int64 = 0
	if cursor != "" {
		start, err = rdb.ZRevRank(ctx, key, cursor).Result() // 查询游标的索引
		if err != nil {
			zap.L().Warn("rdb.ZRevRank failed",
				zap.String("cursor", cursor),
				zap.Error(err),
			)
			start = 0
		}
	}
	// 根据索引查询帖子
	postIDStrs, err = rdb.ZRevRange(ctx, key, start, start+pageSize-1).Result()
	if err != nil {
		zap.L().Error("rdb.ZRevRange failed",
			zap.Int64("start", start),
			zap.Int64("stop", start+pageSize-1),
			zap.Error(err),
		)
		return nil, err
	}
	return postIDStrs, nil
}
