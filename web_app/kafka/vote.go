package kafka

import (
	"context"
	"web_app/dao/mysql"
	"web_app/models"

	"go.uber.org/zap"
)

// insertVoteInMysql 将投票数据插入mysql
func insertVoteInMysql(ctx context.Context, msg *models.VotePost) (err error) {
	// 查询数据库是否已经有投票记录
	exist, err := mysql.VotePostExist(ctx, msg)
	if err != nil {
		zap.L().Error("VotePost existed",
			zap.Int64("post_id", msg.PostID),
			zap.Int64("user_id", msg.UserID),
			zap.Int8("vote_type", msg.VoteType),
			zap.Error(err),
		)
		return err
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
			return err
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
			return err
		}
	}
	return nil
}
