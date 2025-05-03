package models

import "time"

type Post struct {
	PostID      int64     `json:"post_id,string" db:"post_id"`
	AuthorID    int64     `json:"author_id,string" db:"author_id"`
	CommunityID int64     `json:"community_id,string" db:"community_id"`
	Status      int32     `json:"status,string" db:"status"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	CreatTime   time.Time `json:"create_time" db:"create_time"`
}

type ApiPostDetail struct {
	AuthorName string `json:"author_name"`
	VoteNum    int64  `json:"vote_num"`
	*Post
	*CommunityDetail `json:"community"`
}

type PostsAndToken struct {
	Token    string           `json:"token"`
	PostList []*ApiPostDetail `json:"post_list"`
}
