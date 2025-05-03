package mysql

import (
	"context"
	"database/sql"
	"errors"
	"web_app/models"
)

// CreatePost 创建新帖子
func CreatePost(ctx context.Context, post *models.Post) (err error) {
	sqlStr := `insert into post (post_id,
		author_id,
		community_id,
		title,
		content)
		values(?,?,?,?,?)
	`
	_, err = db.ExecContext(ctx, sqlStr, post.PostID, post.AuthorID, post.CommunityID, post.Title, post.Content)
	return err
}

// GetPostIDs 获取所有帖子的id
func GetPostIDs() (ids []int64, err error) {
	sqlStr := `select post_id from post`
	err = db.Select(&ids, sqlStr)
	return ids, err
}

// GetPost 通过id获得帖子信息
func GetPost(ctx context.Context, postID int64) (post *models.Post, err error) {
	sqlStr := `select 
				post_id, author_id, community_id, title, content, create_time
				from
				post
				where post_id = ?
	`
	post = new(models.Post)
	err = db.GetContext(ctx, post, sqlStr, postID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorPostNotExist
		}
		return nil, err
	}
	return post, nil
}
