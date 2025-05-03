package controller

import "web_app/models"

// _ResponseCommunityList 返回社区列表
type _ResponseCommunityList struct {
	Code    ResCode             `json:"code"`    // 业务响应状态码
	Message string              `json:"message"` // 提示信息
	Data    []*models.Community `json:"data"`    // 社区列表data
}

// _ResponseCommunityDetail 返回社区信息详情
type _ResponseCommunityDetail struct {
	Code    ResCode                `json:"code"`    // 业务响应状态码
	Message string                 `json:"message"` // 提示信息
	Data    models.CommunityDetail `json:"data"`    // 社区详细信息data
}

// _ResponsePostDetail 返回帖子详情
type _ResponsePostDetail struct {
	Code    ResCode              `json:"code"`    // 业务响应状态码
	Message string               `json:"message"` // 提示信息
	Data    models.ApiPostDetail `json:"data"`    // 帖子详情data
}

// _Response 基本响应参数
type _Response struct {
	Code    ResCode     `json:"code"`    // 业务响应状态码
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // data
}

// _ResponsePosts 返回帖子列表和pageToken
type _ResponsePosts struct {
	Code    ResCode               `json:"code"`    // 业务响应状态码
	Message string                `json:"message"` // 提示信息
	Data    *models.PostsAndToken `json:"data"`    // 帖子列表和pageToken
}
