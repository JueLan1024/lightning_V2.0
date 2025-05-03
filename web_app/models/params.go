package models

// ParamSignUp 注册请求的结构体
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 登录请求的结构体
type ParamLogin struct {
	Username string `json:"username" binding:"required" example:"juelan"`
	Password string `json:"password" binding:"required" example:"123"`
}

// ParamCommunity 创建社区请求的参数结构体
type ParamCommunity struct {
	CommuntiyID   int64  `json:"community_id,string" binding:"required"`
	CommunityName string `json:"community_name" binding:"required"`
	Introduction  string `json:"introduction" binding:"required"`
}

// ParamPost 创建帖子请求的参数结构体
type ParamPost struct {
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content" binding:"required"`
	CommunityID int64  `json:"community_id,string" binding:"required"`
}

// ParamVoteForPost 给帖子投票的参数
type ParamVoteForPost struct {
	VoteType int8 `json:"vote_type,string" binding:"oneof=1 0 -1" example:"1"` //投票参数{1,0,-1}
	// UserID string `json:"user_id"` 从当前登录用户获取ID
	PostID int64 `json:"post_id,string" binding:"required" example:"7549250837680128"` //帖子ID
}

// ParamGetPosts获取帖子列表参数
type ParamGetPostsInOrder struct {
	CommunityID int64  `json:"community_id,string" form:"community_id" binding:"required" example:"1"`
	Token       string `json:"token" form:"token"`
	Order       string `json:"order" form:"order" example:"score"`
}
