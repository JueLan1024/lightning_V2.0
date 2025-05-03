package logic

import (
	"context"
	"strconv"
	"time"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/pkg/bloom"
	"web_app/pkg/snowflake"

	"go.uber.org/zap"
)

const (
	DefaultPageSize     int64         = 2             //默认每页显示数量
	DefaultCursor       string        = ""            //默认第一篇帖子
	PageTokenExpireTime time.Duration = 4 * time.Hour //pageToken过期时间
)

// CreatePost 创建帖子业务
func CreatePost(ctx context.Context, p *models.ParamPost, authorID int64) (err error) {
	postID := snowflake.GenID()

	post := &models.Post{
		PostID:      postID,
		AuthorID:    authorID,
		CommunityID: p.CommunityID,
		Title:       p.Title,
		Content:     p.Content,
	}
	// 存入数据库
	if err = mysql.CreatePost(ctx, post); err != nil {
		zap.L().Error("mysql.CreatePost(post) failed",
			zap.Int64("authorID", authorID),
			zap.Int64("communityID", p.CommunityID),
			zap.Error(err),
		)
		return err
	}
	// 将帖子ID存入布隆过滤器
	bloom.PostBloomFilter.AddString(strconv.FormatInt(postID, 10))
	return nil
}

// GetPostDetail 获得帖子信息业务
func GetPostDetail(ctx context.Context, postID int64) (postDetail *models.ApiPostDetail, err error) {
	// 用布隆过滤器判断帖子id是否存在
	if !bloom.IsPostIDExist(postID) {
		return nil, ErrorPostNotExist
	}
	// 用singleFlight防止缓存击穿
	post, err := getPostDetailSingleFlight(ctx, postID)
	if err != nil {
		zap.L().Error("getPostDetailSingleFlight failed", zap.Error(err))
		return nil, err
	}
	// 获取帖子的社区信息
	community, err := GetCommunityDetail(ctx, post.CommunityID)
	if err != nil {
		zap.L().Error("GetCommunityDetail failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err),
		)
	}
	// 获取作者用户名
	authorName, err := mysql.GetUserName(ctx, post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserName failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err),
		)
	}
	// 获取投票数
	voteNum, err := redis.GetVoteNum(ctx, postID)
	if err != nil {
		zap.L().Error("redis.GetVoteNum failed",
			zap.Int64("post_id", postID),
			zap.Error(err),
		)
	}
	// 合并数据
	postDetail = &models.ApiPostDetail{
		AuthorName:      authorName,
		VoteNum:         voteNum,
		Post:            post,
		CommunityDetail: community,
	}
	return postDetail, nil
}

// getPostDetailSingleFlight 使用singleFlight获得帖子信息
func getPostDetailSingleFlight(ctx context.Context, postID int64) (post *models.Post, err error) {
	key := redis.GetKeyPostHash(postID)
	v, err, _ := g.Do(key, func() (interface{}, error) {
		// 查缓存
		post, err = redis.GetPost(ctx, key)
		if err == nil {
			return post, nil
		}
		// 缓存没数据查数据库
		if err == redis.ErrorDataNotFound {
			zap.L().Warn("post not found in redis",
				zap.Int64("post_id", postID),
			)
			post, err := mysql.GetPost(ctx, postID)
			if err == nil { // 查到数据设置缓存
				err = redis.InsertPost(ctx, post)
				if err != nil {
					zap.L().Error("redis.CreatePost failed",
						zap.Int64("post_id", postID),
						zap.Error(err),
					)
				}
				return post, nil
			}
			if err == mysql.ErrorPostNotExist {
				zap.L().Error("post not exists",
					zap.Int64("post_id", postID),
					zap.Error(err),
				)
				return nil, err
			}
			zap.L().Error("mysql.GetPost failed",
				zap.Int64("post_id", postID),
				zap.Error(err),
			)
			return nil, err
		}
		// 缓存出错直接返回，防止灾难传递至DB
		zap.L().Error("redis.GetPost failed",
			zap.Int64("post_id", postID),
			zap.Error(err),
		)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	// 格式转换
	post, ok := v.(*models.Post)
	if !ok {
		zap.L().Error("parse interface{} failed")
		return nil, err
	}
	return post, nil
}

// GetPostList 获得帖子列表业务
func GetPostList(ctx context.Context, p *models.ParamGetPostsInOrder) (postsAndToken *models.PostsAndToken, err error) {
	// 默认第一页开始
	pageSize := DefaultPageSize
	cursor := DefaultCursor
	// 解析Token
	if len(p.Token) > 0 {
		pageInfo := models.Token(p.Token).Decode()
		if pageInfo.InValid() { //解析结果无效返回错误
			return nil, ErrorInvalidPageToken
		}
		pageSize = pageInfo.PageSize
		cursor = pageInfo.NextID
	}
	// 在redis中查询帖子ID列表
	postIDStrs, err := redis.GetPostList(ctx, p, cursor, pageSize+1)
	if err != nil {
		zap.L().Error("redis.GetPostList failed",
			zap.Int64("community_id", p.CommunityID),
			zap.String("order", p.Order),
			zap.String("cursor", cursor),
			zap.Error(err),
		)
		return nil, err
	}

	// 判断是否还有下一页
	var nextPageToken string
	realPageSize := int(pageSize)
	if len(postIDStrs) > int(pageSize) {
		nextPageInfo := &models.Page{
			NextID:        postIDStrs[pageSize],
			NextTimeAtUTC: time.Now().Add(PageTokenExpireTime).Unix(),
			PageSize:      DefaultPageSize,
		}
		nextPageToken = string(nextPageInfo.Encode())
	} else {
		realPageSize = len(postIDStrs)
	}
	postIDStrs = postIDStrs[:realPageSize]
	//通过帖子ID列表查询帖子信息
	postList := make([]*models.ApiPostDetail, 0, realPageSize)
	for _, postIDStr := range postIDStrs {
		postID, err := strconv.ParseInt(postIDStr, 10, 64)
		if err != nil {
			zap.L().Error("strconv.ParseInt failed",
				zap.String("postIDStr", postIDStr),
				zap.Error(err),
			)
			continue
		}
		postDetial, err := GetPostDetail(ctx, postID)
		if err != nil {
			zap.L().Error("GetPostDetail failed",
				zap.Int64("postID", postID),
				zap.Error(err),
			)
			continue
		}
		postList = append(postList, postDetial)
	}
	postsAndToken = &models.PostsAndToken{
		Token:    nextPageToken,
		PostList: postList,
	}
	return postsAndToken, nil
}
