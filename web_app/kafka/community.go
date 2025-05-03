package kafka

import (
	"context"
	"strconv"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/tool"

	"go.uber.org/zap"
)

// insertCommunityInRedis 将社区数据插入redis
func insertCommunityInRedis(ctx context.Context, d map[string]interface{}) (err error) {
	idStr, ok := d["community_id"].(string)
	if !ok {
		zap.L().Error("invalid type for community_id")
		return ErrorInvalidDataType
	}
	name, ok := d["community_name"].(string)
	if !ok {
		zap.L().Error("invalid type for community_name")
		return ErrorInvalidDataType
	}
	introduction, ok := d["introduction"].(string)
	if !ok {
		zap.L().Error("invalid type for introduction")
		return ErrorInvalidDataType
	}
	createTimeStr, ok := d["create_time"].(string)
	if !ok {
		zap.L().Error("invalid type for create_time")
		return ErrorInvalidDataType
	}
	createTime, err := tool.ParseTime(createTimeStr)
	if err != nil {
		zap.L().Error("tool.ParseTime failed", zap.Error(err))
		return err
	}
	id, _ := strconv.ParseInt(idStr, 10, 64)
	community := &models.CommunityDetail{
		CommunityID:   id,
		CommunityName: name,
		Introduction:  introduction,
		CreateTime:    createTime,
	}
	//社区信息存入redis,Key  lightning:community:<community_id> 社区id存入redis,Key  lightning:community:list
	if err = redis.CreateCommunityDetail(ctx, community); err != nil {
		zap.L().Error("redis.CreateCommunityDetail failed", zap.Error(err))
		return err
	}
	return nil
}
