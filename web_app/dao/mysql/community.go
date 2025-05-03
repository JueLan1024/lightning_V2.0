package mysql

import (
	"context"
	"database/sql"
	"errors"
	"web_app/models"
)

// CommunityExists 查看communityID是否存在
func CommunityExists(ctx context.Context, communityID int64) (err error) {
	sqlStr := `select count(community_id) from community where community_id = ?`
	var count int64
	if err = db.GetContext(ctx, &count, sqlStr, communityID); err != nil {
		return err
	}
	if count > 0 {
		return ErrorCommunityIDExist
	}
	return nil
}

// CreateCommunity 创建社区
func CreateCommunity(ctx context.Context, p *models.ParamCommunity) (err error) {
	sqlStr := `insert into community (community_id, community_name, introduction) values (?,?,?) `
	_, err = db.ExecContext(ctx, sqlStr, p.CommuntiyID, p.CommunityName, p.Introduction)
	return err
}

// GetCommunityList 从数据库中获得社区列表
func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := `select community_id, community_name from community`
	err = db.Select(&communityList, sqlStr)
	return communityList, err
}

// GetCommunityDetailList 从数据库中获得社区细节列表
func GetCommunityDetailList() (communityDetailList []*models.CommunityDetail, err error) {
	sqlStr := `select community_id, community_name, introduction, create_time from community`
	err = db.Select(&communityDetailList, sqlStr)
	return communityDetailList, err
}

// GetCommunityIDs 从数据库中获得社区ids
func GetCommunityIDs() (ids []int64, err error) {
	sqlStr := `select community_id from community`
	err = db.Select(&ids, sqlStr)
	return ids, err
}

// GetCommunityDetail 通过id从数据库获得社区信息
func GetCommunityDetail(communityID int64) (communityDetail *models.CommunityDetail, err error) {
	sqlStr := `select 
		community_id, community_name, introduction, create_time 
		from community 
		where community_id = ?
	`
	communityDetail = new(models.CommunityDetail)
	if err = db.Get(communityDetail, sqlStr, communityID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorCommunityNotExist
		}
		return nil, err
	}
	return communityDetail, nil
}
