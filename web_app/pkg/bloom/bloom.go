package bloom

import (
	"strconv"
	"web_app/dao/mysql"

	"github.com/bits-and-blooms/bloom/v3"
	"go.uber.org/zap"
)

var CommunityBloomFilter *bloom.BloomFilter // 社区布隆过滤器
var PostBloomFilter *bloom.BloomFilter      // 帖子布隆过滤器

// 初始化布隆过滤器
func InitBloomFilter() (err error) {
	if err = initCommunityBloom(); err != nil {
		zap.L().Error("initCommunityBloom() failed", zap.Error(err))
		return err
	}
	if err = initPostBloom(); err != nil {
		zap.L().Error("initPostBloom() failed", zap.Error(err))
		return err
	}
	return nil
}

// initCommunityBloom 初始社区布隆过滤器
func initCommunityBloom() (err error) {
	// 创建布隆过滤器
	CommunityBloomFilter = bloom.New(1000000, 5)

	// 从MySQL中加载历史数据
	ids, err := mysql.GetCommunityIDs()
	if err != nil {
		zap.L().Error("failed to get community ids in mysql", zap.Error(err))
		return err
	}

	for _, id := range ids {
		idStr := strconv.FormatInt(id, 10)
		CommunityBloomFilter.AddString(idStr)
	}
	zap.L().Info("Community Bloom Filter initialized with historical data")
	return nil
}

// initPostBloom 初始帖子布隆过滤器
func initPostBloom() (err error) {
	// 创建布隆过滤器
	PostBloomFilter = bloom.New(1000000, 5)
	// 从mysql中加载历史数据
	ids, err := mysql.GetPostIDs()
	if err != nil {
		zap.L().Error("failed to get post ids in mysql", zap.Error(err))
		return err
	}
	// 将Ids存入布隆过滤器
	for _, id := range ids {
		idStr := strconv.FormatInt(id, 10)
		PostBloomFilter.AddString(idStr)
	}
	zap.L().Info("Post Bloom Filter initialized with historical data")
	return nil
}

// IsCommunityIDExist 通过ID判断社区是否存在
func IsCommunityIDExist(id int64) bool {
	idStr := strconv.FormatInt(id, 10)
	return CommunityBloomFilter.TestString(idStr)
}

// IsPostIDExist 通过ID判断帖子是否存在
func IsPostIDExist(id int64) bool {
	idStr := strconv.FormatInt(id, 10)
	return PostBloomFilter.TestString(idStr)
}
