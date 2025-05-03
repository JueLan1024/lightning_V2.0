package controller

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUsernameExist
	CodeUsernameOrPasswordWrong
	CodeInvalidToken
	CodeNeedLogin
	CodeNewToken
	CodeCommunityExists
	CodeCommunityNotExists
	CodePostNotExists
	CodeVoteRepeated
	CodeInvalidPageToken
	CodeServerBusy
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:                 "success",
	CodeInvalidParam:            "参数错误",
	CodeUsernameExist:           "用户名已存在",
	CodeUsernameOrPasswordWrong: "用户名或密码错误",
	CodeInvalidToken:            "invalid token",
	CodeNeedLogin:               "未登录",
	CodeNewToken:                "新Token",
	CodeCommunityExists:         "社区已存在",
	CodeCommunityNotExists:      "社区不存在",
	CodePostNotExists:           "帖子不存在",
	CodeVoteRepeated:            "重复投票",
	CodeInvalidPageToken:        "invalid page token",
	CodeServerBusy:              "服务繁忙",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
