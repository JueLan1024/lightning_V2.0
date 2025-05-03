package mysql

import (
	"context"
	"web_app/models"
)

// VotePostExist 查询投票数据是否存在
func VotePostExist(ctx context.Context, msg *models.VotePost) (exist bool, err error) {
	sqlStr := `select count(id) from vote_post where post_id = ? and user_id = ?`
	var count int
	if err = db.Get(&count, sqlStr, msg.PostID, msg.UserID); err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

// CreateVotePost 创建用户给帖子投票表
func CreateVotePost(ctx context.Context, msg *models.VotePost) (err error) {
	sqlStr := `insert into vote_post (post_id, user_id, vote_type) values (?,?,?)`
	_, err = db.Exec(sqlStr, msg.PostID, msg.UserID, msg.VoteType)
	return err
}

// UpdateVotePost 更新投票表
func UpdateVotePost(ctx context.Context, msg *models.VotePost) (err error) {
	sqlStr := `update vote_post
				set vote_type = ?
				where post_id = ? and user_id = ?
	`
	_, err = db.Exec(sqlStr, msg.VoteType, msg.PostID, msg.UserID)
	return err
}
