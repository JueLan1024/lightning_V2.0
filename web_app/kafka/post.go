package kafka

import (
	"context"
	"strconv"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/tool"

	"go.uber.org/zap"
)

// 将帖子数据插入redis
func insertPostInRedis(ctx context.Context, msg map[string]interface{}) (err error) {
	idStr, ok := msg["post_id"].(string)
	if !ok {
		zap.L().Error("Invalid type for post_id")
		return ErrorInvalidDataType
	}
	title, ok := msg["title"].(string)
	if !ok {
		zap.L().Error("Invalid type for title")
		return ErrorInvalidDataType
	}
	content, ok := msg["content"].(string)
	if !ok {
		zap.L().Error("Invalid type for content")
		return ErrorInvalidDataType
	}
	authorIDStr, ok := msg["author_id"].(string)
	if !ok {
		zap.L().Error("Invalid type for author_id")
		return ErrorInvalidDataType
	}
	communityIDStr, ok := msg["community_id"].(string)
	if !ok {
		zap.L().Error("Invalid type for community_id")
		return ErrorInvalidDataType
	}
	createTimeStr, ok := msg["create_time"].(string)
	if !ok {
		zap.L().Error("Invalid type for create_time")
		return ErrorInvalidDataType
	}
	// 转换数据类型
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		zap.L().Error("strconv.ParseInt failed")
	}
	authorID, err := strconv.ParseInt(authorIDStr, 10, 64)
	if err != nil {
		zap.L().Error("strconv.ParseInt failed")
	}
	communityID, err := strconv.ParseInt(communityIDStr, 10, 64)
	if err != nil {
		zap.L().Error("strconv.ParseInt failed")
	}
	createTime, err := tool.ParseTime(createTimeStr)
	if err != nil {
		zap.L().Error("tool.ParseTime failed")
	}
	post := &models.Post{
		PostID:      postID,
		AuthorID:    authorID,
		CommunityID: communityID,
		Title:       title,
		Content:     content,
		CreatTime:   createTime,
	}

	// 将帖子存入redis
	if err = redis.CreatePost(ctx, post); err != nil {
		zap.L().Error("redis.CreatePost failed",
			zap.Int64("postID", postID),
			zap.String("authorID", authorIDStr),
		)
		return err
	}
	return nil
}
