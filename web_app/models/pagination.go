package models

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type Page struct {
	NextID        string `json:"next_id"`          // 下一篇帖子
	NextTimeAtUTC int64  `json:"next_time_at_utc"` //  token过期时间
	PageSize      int64  `json:"page_size"`        // 查询的帖子数量
}

type Token string

// Encode 返回分页token
func (p Page) Encode() Token {
	b, err := json.Marshal(p)
	if err != nil {
		return Token("")
	}
	return Token(base64.StdEncoding.EncodeToString(b))
}

// InValid判断page是否无效
func (p Page) InValid() bool {
	return p.NextID == "" || p.NextTimeAtUTC < time.Now().Unix() || p.PageSize <= 0
}

// Decode 解析分页信息
func (t Token) Decode() Page {
	var result Page
	if len(t) == 0 {
		return result
	}

	bytes, err := base64.StdEncoding.DecodeString(string(t))
	if err != nil {
		return result
	}

	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return result
	}

	return result
}
