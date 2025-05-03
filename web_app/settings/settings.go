package settings

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(AppConf)

type AppConf struct {
	Name                 string        `mapstructure:"name"`
	Version              string        `mapstructure:"version"`
	Mode                 string        `mapstructure:"mode"`
	Host                 string        `mapstructure:"host"`
	StartTime            string        `mapstructure:"start_time"`
	MachineID            int64         `mapstructure:"machine_id"`
	Port                 int           `mapstructure:"port"`
	AccessTokenDuration  time.Duration `mapstructure:"access_token_duration"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
	*LogConfig           `mapstructure:"log"`
	*MysqlConfig         `mapstructure:"mysql"`
	*RedisConfig         `mapstructure:"redis"`
	*KafkaConfig         `mapstructure:"kafka"`
	*RatelimitConfig     `mapstructure:"ratelimit"`
}

type LogConfig struct {
	Filename   string `mapstructure:"filename"`
	Level      string `mapstructure:"level"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
}

type MysqlConfig struct {
	Host         string `mapstructure:"host"`
	Dbname       string `mapstructure:"dbname"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	Db       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type KafkaConfig struct {
	Brokers          []string `mapstructure:"brokers"`
	GroupIDCommunity string   `mapstructure:"group_id_community"`
	GroupIDPost      string   `mapstructure:"group_id_post"`
	GroupIDVotePost  string   `mapstructure:"group_id_vote_post"`
	TopicCommunity   string   `mapstructure:"topic_community"`
	TopicPost        string   `mapstructure:"topic_post"`
	TopicVotePost    string   `mapstructure:"topic_vote_post"`
}
type RatelimitConfig struct {
	FillInterval time.Duration `mapstructure:"fill_interval"`
	Cap          int64         `mapstructure:"cap"`
}

func Init() (err error) {
	viper.SetConfigFile("./conf/config.yaml") //指定准确的文件路径、文件名和文件类型
	// viper.SetConfigName("config") //指定配置文件名称
	// viper.SetConfigType("yaml")    //指定文件类型
	// viper.AddConfigPath("./conf/") //指定查找配置文件的路径
	err = viper.ReadInConfig() //读取配置文件信息
	if err != nil {
		// 读取配置文件失败
		fmt.Printf("viper.ReadInConfig() failed,err:%v\n", err)
		return
	}
	// 将配置信息反序列化到结构体Conf中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal(Conf) failed,err:%v\n", err)
	}
	viper.WatchConfig() //监视配置文件变化
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("配置文件修改了: %s, 操作类型: %v\n", in.Name, in.Op)
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal(Conf) failed,err:%v\n", err)
		}
	})
	return
}
