package redis

import "fmt"

const (
	Prefix         = "lightning:"          // 全局前缀
	KeyUserPF      = Prefix + "user:"      // 用户模块前缀
	KeyCommunityPF = Prefix + "community:" //社区模块前缀
	KeyPostPF      = Prefix + "post:"      //帖子模块前缀
	KeyVotePostPF  = Prefix + "vote:post:" //帖子投票模块前缀
)

// KeyUserRefreshToken 获取用户RefreshToken的Key,键值对存储方式
// lightning:user:<user_id>
func GetKeyUserRefreshToken(userID int64) string {
	return fmt.Sprintf("%s%d:refresh_token", KeyUserPF, userID)
}

// GetKeyCommunityListHash 获取社区详细信息的Key,Hash存储方式
// lightning:community:<community_id>
func GetKeyCommunityHash(communityID int64) string {
	return fmt.Sprintf("%s%d", KeyCommunityPF, communityID)
}

// GetKeyCommunityIDsZSet 获取社区IDs的Key,ZSet存储方式
// lightning:community:list
func GetKeyCommunityIDsZSet() string {
	return KeyCommunityPF + "list"
}

// GetKeyPostHash 获取帖子信息的Key,Hash存储方式
// lightning:post:<post_id>
func GetKeyPostHash(postID int64) string {
	return fmt.Sprintf("%s%d", KeyPostPF, postID)
}

// GetKeyPostTimeZSet 获取帖子按创建时间排序的key,ZSet存储方式
// lightning:post:time
func GetKeyPostTimeZSet() string {
	return KeyPostPF + "time"
}

// GetKeyPostScoreZSet 获取帖子按分数排序的key,ZSet存储方式
// lightning:post:score
func GetKeyPostScoreZSet() string {
	return KeyPostPF + "score"
}

// GetKeyCommunityPostsSet 获取存储社区中所有帖子的key,Set存储方式
// lightning:community:<community_id>:posts
func GetKeyCommunityPostsSet(communityID int64) string {
	return fmt.Sprintf("%s%d:posts", KeyCommunityPF, communityID)
}

// GetKeyVotePostHash 获取用户给帖子投票的Key,Hash存储方式,键为user_id,值为vote_type
//  lightning:vote:post:<post_id>
func GetKeyVotePostHash(postID int64) string {
	return fmt.Sprintf("%s%d", KeyVotePostPF, postID)
}

// GetKeyCommunityPostScoreZSet 获取按分数排序的社区帖子的Key,ZSet存储方式
// lightning:community:<community_id>:post:score
func GetKeyCommunityPostScoreZSet(communityID int64) string {
	return fmt.Sprintf("%s%d:post:score", KeyCommunityPF, communityID)
}

// GetKeyCommunityPostTimeZSet 获取按时间排序的社区帖子的Key,ZSet存储方式
// lightning:community:<community_id>:post:time
func GetKeyCommunityPostTimeZSet(communityID int64) string {
	return fmt.Sprintf("%s%d:post:time", KeyCommunityPF, communityID)
}
