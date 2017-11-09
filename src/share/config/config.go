package config

import (
	microConfig "github.com/micro/go-os/config"
	"github.com/micro/go-os/config/source/consul"
)

// ConfigPrefix 配置的前缀
const ConfigPrefix = "wolf/"

// 唯一配置实例
var instance *Config

// Instance 唯一配置实例
func Instance() *Config {
	return instance
}

// Init 初始化
func Init(sourceNames ...string) *Config {
	InitEnv()
	instance = createConfig(sourceNames...)
	return instance
}

// Config 配置对象
type Config struct {
	microConfig.Config
}

// CreateConfig 创建一个配置对象
func createConfig(sourceNames ...string) *Config {
	options := []microConfig.Option{}
	for _, sourceName := range sourceNames {
		options = append(options, microConfig.WithSource(consul.NewSource(
			microConfig.SourceName(ConfigPrefix+CurrentEnv+"/"+CurrentZoneArea+"/"+sourceName),
		)))
	}
	return &Config{
		Config: microConfig.NewConfig(options...),
	}
}

// CfgMustGetStringSlice 获取字符串数组, 为空则Panic
func (c *Config) CfgMustGetStringSlice(key string) []string {
	strs := c.Get(key).StringSlice([]string{})
	if len(strs) == 0 {
		panic("empty " + key)
	}
	return strs
}

// CfgMustGetString 获取字符串, 为空则Panic
func (c *Config) CfgMustGetString(key string) string {
	str := c.Get(key).String("")
	if len(str) == 0 {
		panic("empty " + key)
	}
	return str
}

// CfgMustGetInt 获取数字, 为0则Panic
func (c *Config) CfgMustGetInt(key string) int {
	integer := c.Get(key).Int(0)
	if integer == 0 {
		panic("empty " + key)
	}
	return integer
}

// IsInMaintenance 检查是否正在维护中
func (c *Config) IsInMaintenance(platform string) bool {
	if len(platform) == 0 || (platform != "ios" && platform != "android") {
		return true
	}
	return c.Get("isInMaintenance_" + platform).Bool(false)
}

// IsUnMaintenanceOfOpenID 是否在白名单中
func (c *Config) IsUnMaintenanceOfOpenID(openID string) bool {
	if len(openID) < 1 {
		return false
	}
	unMainOpenIDs := c.Get("unMaintenanceOpenIDs").StringSlice([]string{})
	for _, unMOpenID := range unMainOpenIDs {
		if openID == unMOpenID {
			return true
		}
	}
	return false
}

// RedisConfig RedisConfig
type RedisConfig struct {
	Address  string `json:"addr"`
	Password string `json:"password"`
	DBNum    int    `json:"db"`
}

// CfgMustGetRedisConfig 获取结构, addredd为空则Panic
func (c *Config) CfgMustGetRedisConfig(key string) RedisConfig {
	redisConf := RedisConfig{}
	err := c.Get(key).Scan(&redisConf)
	if err != nil {
		panic(key + " error:" + err.Error())
	}
	if len(redisConf.Address) == 0 {
		panic("redis address empty")
	}
	return redisConf
}
