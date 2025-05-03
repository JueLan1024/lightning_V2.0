package redis

import (
	"context"
	"strconv"
	"time"
	"web_app/models"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// GetCommunitList 获得社区列表
func GetCommunitList(ctx context.Context, key string) (communitList []*models.Community, err error) {
	// 获得社区ids
	idStrs, err := getCommunityIDs(ctx, key)

	if err != nil {
		zap.L().Error("failed to get community IDs", zap.Error(err))
		return nil, err
	}
	if len(idStrs) == 0 { //没有数据返回错误
		return nil, ErrorDataNotFound
	}
	// 根据社区ids获得社区名称
	pipe := rdb.Pipeline()
	for _, idStr := range idStrs {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			zap.L().Error("strconv.ParseInt failed", zap.Error(err))
			continue
		}
		pipe.HMGet(ctx, GetKeyCommunityHash(id), "community_id", "community_name")
	}

	// 执行管道中的命令
	cmder, err := pipe.Exec(ctx)
	if err != nil {
		zap.L().Error("failed to execute pipeline", zap.Error(err))
		return
	}

	// 解析数据
	communitList = make([]*models.Community, 0, len(cmder))
	for _, cmd := range cmder {
		results, err := cmd.(*redis.SliceCmd).Result()
		if err != nil {
			zap.L().Error("cmd.(*redis.SliceCmd).Result() failed", zap.Error(err))
			continue
		}
		if len(results) == 0 || len(results)%2 != 0 || results[0] == nil || results[1] == nil { //result长度内容不正确
			zap.L().Warn("wrong results from cmder")
			continue
		}
		for i := 0; i < len(results); i += 2 {
			idStr, ok1 := results[i].(string)
			name, ok2 := results[i+1].(string)
			if !ok1 || !ok2 {
				zap.L().Warn("failed to get id&name from results")
				continue
			}
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				zap.L().Error("strconv.ParseInt failed", zap.Error(err))
				continue
			}
			community := &models.Community{
				CommunityID:   id,
				CommunityName: name,
			}
			communitList = append(communitList, community)
		}
	}
	return communitList, err
}

// getCommunityIDs 获得社区ids
func getCommunityIDs(ctx context.Context, key string) (idStrs []string, err error) {
	// 从Redis获得社区IDs
	return rdb.ZRange(ctx, key, 0, -1).Result()
}

// CreateCommunityDetail 创建社区信息 Key lightning:community:<community_id> ; lightning:community:list
func CreateCommunityDetail(ctx context.Context, community *models.CommunityDetail) (err error) {
	pipe := rdb.TxPipeline()
	key := GetKeyCommunityHash(community.CommunityID)
	// 存社区信息 Hash存储方式
	pipe.HMSet(ctx, key,
		"community_id", community.CommunityID,
		"community_name", community.CommunityName,
		"introduction", community.Introduction,
		"create_time", community.CreateTime,
	)
	// 存社区ID ZSet存储方式
	key = GetKeyCommunityIDsZSet()
	pipe.ZAdd(ctx, key, &redis.Z{
		Score:  float64(community.CreateTime.Unix()),
		Member: community.CommunityID,
	})
	_, err = pipe.Exec(ctx)
	return err
}

// GetCommunityDetail 通过id获取社区信息
func GetCommunityDetail(ctx context.Context, key string) (community *models.CommunityDetail, err error) {
	data, err := rdb.HGetAll(ctx, key).Result()
	if len(data) == 0 { // key不存在返回数据未找到错误
		return nil, ErrorDataNotFound
	}
	if err != nil {
		zap.L().Error("GetCommunityDetail failed", zap.Error(err))
		return nil, err
	}
	idStr := data["community_id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		zap.L().Error("strconv.ParseInt(idStr,10,64) failed", zap.Error(err))
		return nil, err
	}
	createTimeStr := data["create_time"]
	createTime, err := time.Parse(time.RFC3339, createTimeStr)
	if err != nil {
		zap.L().Error("createTimeStr time.Parse failed failed", zap.Error(err))
		return nil, err
	}
	community = &models.CommunityDetail{
		CommunityID:   id,
		CommunityName: data["community_name"],
		Introduction:  data["introduction"],
		CreateTime:    createTime,
	}
	return community, nil
}

// // GetCommunityName 获得社区名称
// func GetCommunityName(ctx context.Context, key string) (communityList []*models.Community, err error) {
// 	// 从Redis获取指定字段
// 	result, err := rdb.HMGet(ctx, key, "community_id", "community_name").Result()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 检查是否有数据
// 	if len(result) == 0 || result[0] == nil || result[1] == nil {
// 		return nil, ErrorDataNotFound
// 	}

// 	// 检查结果长度是否匹配
// 	if len(result)%2 != 0 {
// 		return nil, ErrorInvalidDataFormat
// 	}

// 	// 遍历结果并解析为 []*models.Community
// 	for i := 0; i < len(result); i += 2 {
// 		idStr, ok1 := result[i].(string)
// 		name, ok2 := result[i+1].(string)
// 		if !ok1 || !ok2 {
// 			return nil, ErrorParseDataFailed
// 		}
// 		// 将字符串ID转换为int64
// 		id, err := strconv.ParseInt(idStr, 10, 64)
// 		if err != nil {
// 			return nil, ErrorParseDataFailed
// 		}

// 		community := &models.Community{
// 			CommunityID:   id,
// 			CommunityName: name,
// 		}
// 		communityList = append(communityList, community)
// 	}
// 	return communityList, nil
// }

// // GetCommunityDetail 获得社区信息(测试用)
// func GetCommunityDetail(ctx context.Context, communityID int64) (err error) {
// 	key := GetKeyCommunityHash(communityID)
// 	datas, err := rdb.HGetAll(ctx, key).Result()
// 	for _, data := range datas {
// 		fmt.Println(data)
// 	}
// 	return
// }

// // CreateCommunityIDs 创建社区ID列表，ZSet存储方式，Member为社区ID，Score为社区帖子数量（Score暂时为0）
// func CreateCommunityIDs(ctx context.Context, ids []int64) (err error) {
// 	key := GetKeyCommunityIDsZSet()
// 	pipeline := rdb.Pipeline()
// 	for _, id := range ids {
// 		// 将int64转为string
// 		idStr := strconv.FormatInt(id, 10)
// 		pipeline.ZAdd(ctx, key, &redis.Z{
// 			Score:  0,
// 			Member: idStr,
// 		})
// 	}
// 	_, err = pipeline.Exec(ctx)
// 	return err
// }
