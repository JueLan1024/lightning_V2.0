package models

import "time"

// VotePost 帖子投票结构体
type VotePost struct {
	PostID     int64     `json:"post_id" db:"post_id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	VoteType   int8      `json:"vote_type" db:"vote_type"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
}
