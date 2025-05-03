package logic

import (
	"context"
	"strconv"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/pkg/bloom"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

var g = new(singleflight.Group)

// CreateCommunity 新建社区
func CreateCommunity(ctx context.Context, p *models.ParamCommunity) (err error) {
	// 查看社区是否存在
	if err = mysql.CommunityExists(ctx, p.CommuntiyID); err != nil {
		zap.L().Error("failed to create community", zap.Int64("community id", p.CommuntiyID), zap.Error(err))
		if err == mysql.ErrorCommunityIDExist {
			return ErrorCommunityExist
		}
		return err
	}
	// 将社区存入Mysql
	if err = mysql.CreateCommunity(ctx, p); err != nil {
		zap.L().Error("falied to insert community in mysql", zap.Int64("community_id", p.CommuntiyID), zap.Error(err))
		return err
	}
	// 将社区ID存入布隆过滤器
	bloom.CommunityBloomFilter.AddString(strconv.FormatInt(p.CommuntiyID, 10))
	return nil
}

// GetCommunityList 获得社区列表
func GetCommunityList(ctx context.Context) (communityList []*models.Community, err error) {
	// 使用singleflight防止缓存击穿
	data, err := getCommunityListSingleFlight(ctx, redis.GetKeyCommunityIDsZSet())
	if err != nil {
		zap.L().Error("getCommunityListSingleFlight failed", zap.Error(err))
		return nil, err
	}
	communityList, ok := data.([]*models.Community)
	if !ok {
		zap.L().Error("parse interface{} failed")
		return nil, err
	}
	return communityList, nil
}

// getCommunityListSingleFlight 使用singleflight获取社区列表
func getCommunityListSingleFlight(ctx context.Context, key string) (interface{}, error) {
	v, err, _ := g.Do(key, func() (interface{}, error) {
		// 查缓存
		list, err := redis.GetCommunitList(ctx, key)
		if err == nil {
			return list, nil
		}
		if err == redis.ErrorDataNotFound {
			zap.L().Warn("community list not found in redis")
			// redis中没有数据查mysql
			detailList, err := mysql.GetCommunityDetailList()
			if err == nil { // 查到数据设置缓存
				communityList := make([]*models.Community, 0, len(detailList))
				for _, detail := range detailList {
					err = redis.CreateCommunityDetail(ctx, detail)
					if err != nil {
						zap.L().Error("redis.CreateCommunityDetail failed",
							zap.Int64("community_id", detail.CommunityID),
							zap.Error(err),
						)
						continue
					}
					// 提取数据
					community := &models.Community{
						CommunityID:   detail.CommunityID,
						CommunityName: detail.CommunityName,
					}
					communityList = append(communityList, community)
				}
				return communityList, nil
			}
			zap.L().Error("failed to get community list in mysql", zap.Error(err))
			return nil, err
		}
		return nil, err // 缓存出错直接返回，防止灾难传递至DB
	})

	if err != nil {
		return nil, err
	}
	return v, nil
}

// GetCommunityDetail 通过id获得社区信息
func GetCommunityDetail(ctx context.Context, id int64) (communityDetail *models.CommunityDetail, err error) {
	// 用布隆过滤器判断社区是否存在
	if !bloom.IsCommunityIDExist(id) {
		return nil, ErrorCommunityNotExist
	}

	// 使用singleflight防止缓存击穿
	communityDetail, err = getCommunityDetailSingleFlight(ctx, id)
	if err != nil {
		zap.L().Error("getCommunityDetailSingleFlight failed",
			zap.Int64("community_id", id),
			zap.Error(err),
		)
		return nil, err
	}

	return communityDetail, nil
}

// getCommunityDetailSingleFlight 使用singleFlight获得社区信息
func getCommunityDetailSingleFlight(ctx context.Context, communityID int64) (communityDetail *models.CommunityDetail, err error) {
	key := redis.GetKeyCommunityHash(communityID)
	v, err, _ := g.Do(key, func() (interface{}, error) {
		// 查缓存
		communityDetail, err := redis.GetCommunityDetail(ctx, key)
		if err == nil {
			return communityDetail, nil
		}
		if err == redis.ErrorDataNotFound { //缓存没数据查数据库
			zap.L().Warn("community not found in redis", zap.Int64("community_id", communityID))
			communityDetail, err := mysql.GetCommunityDetail(communityID)
			if err == nil { //查到数据设置缓存
				err = redis.CreateCommunityDetail(ctx, communityDetail)
				if err != nil {
					zap.L().Error("redis.CreateCommunityDetail failed",
						zap.Int64("community_id", communityID),
						zap.Error(err),
					)
					return nil, err
				}
				// 返回数据
				return communityDetail, nil
			}
			if err == mysql.ErrorCommunityNotExist {
				zap.L().Error("community not found in mysql",
					zap.Int64("community_id", communityID),
					zap.Error(err),
				)
				return nil, ErrorCommunityNotExist
			}
			zap.L().Error("mysql.GetCommunityDetail(communityID) failed",
				zap.Int64("community_id", communityID),
				zap.Error(err),
			)
			return nil, err
		}
		// 缓存出错直接返回，防止灾难传递至DB
		zap.L().Error("redis.GetCommunityDetail() failed",
			zap.Int64("community_id", communityID),
			zap.Error(err),
		)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	// 格式转换
	communityDetail, ok := v.(*models.CommunityDetail)
	if !ok {
		zap.L().Error("parse interface{} failed")
		return nil, err
	}
	return communityDetail, nil
}

// // SetCommunityIDs 从mysql中获取所有社区的id并存入redis
// func SetCommunityIDs(ctx context.Context) (err error) {
// 	ids, err := mysql.GetCommunityIDs()
// 	if err != nil {
// 		zap.L().Error("failed to get community ids in mysql", zap.Error(err))
// 		return err
// 	}
// 	err = redis.CreateCommunityIDs(ctx, ids)
// 	if err != nil {
// 		zap.L().Error("failed to create community ids in redis", zap.Error(err))
// 		return err
// 	}
// 	return nil
// }
