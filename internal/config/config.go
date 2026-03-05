package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config 全局配置结构
var Config *viper.Viper

// AppConfig 应用配置
type AppConfig struct {
	Name     string `mapstructure:"name"`
	Version  string `mapstructure:"version"`
	Debug    bool   `mapstructure:"debug"`
	LogLevel string `mapstructure:"log_level"`
	LogFile  string `mapstructure:"log_file"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type string `mapstructure:"type"`
	Path string `mapstructure:"path"`
}

// WebConfig Web 服务配置
type WebConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// CityCode 城市编码
type CityCode struct {
	Code string `mapstructure:"code"`
	Name string `mapstructure:"name"`
}

// SearchConfig 搜索配置
type SearchConfig struct {
	CityCodes  []CityCode `mapstructure:"city_codes"`
	Keywords   []string   `mapstructure:"keywords"`
	JobType    string     `mapstructure:"job_type"`
	Salary     string     `mapstructure:"salary"`
	Experience string     `mapstructure:"experience"`
	Degree     string     `mapstructure:"degree"`
}

// AIConfig AI 配置
type AIConfig struct {
	Enable    bool          `mapstructure:"enable"`
	ModelType string        `mapstructure:"model_type"`
	Qwen      QwenConfig    `mapstructure:"qwen"`
	Minimax   MinimaxConfig `mapstructure:"minimax"`
}

// QwenConfig 阿里 Qwen 配置
type QwenConfig struct {
	APIKey string `mapstructure:"api_key"`
	Model  string `mapstructure:"model"`
}

// MinimaxConfig Minimax 配置
type MinimaxConfig struct {
	APIKey string `mapstructure:"api_key"`
	Model  string `mapstructure:"model"`
}

// GreetingConfig 打招呼配置
type GreetingConfig struct {
	Default    string `mapstructure:"default"`
	EnableAI   bool   `mapstructure:"enable_ai"`
	AITemplate string `mapstructure:"ai_template"`
}

// DeliveryConfig 投递配置
type DeliveryConfig struct {
	SendImgResume bool   `mapstructure:"send_img_resume"`
	ImgResumePath string `mapstructure:"img_resume_path"`
	DailyLimit    int    `mapstructure:"daily_limit"`
}

// FilterConfig 过滤配置
type FilterConfig struct {
	FilterDeadHR     bool     `mapstructure:"filter_dead_hr"`
	CompanyBlacklist []string `mapstructure:"company_blacklist"`
	HRBlacklist      []string `mapstructure:"hr_blacklist"`
	JobBlacklist     []string `mapstructure:"job_blacklist"`
}

// BrowserConfig 浏览器配置
type BrowserConfig struct {
	Type         string `mapstructure:"type"`
	WindowWidth  int    `mapstructure:"window_width"`
	WindowHeight int    `mapstructure:"window_height"`
	UseExisting  bool   `mapstructure:"use_existing"`
	UserDataDir  string `mapstructure:"user_data_dir"`
}

// CookieConfig Cookie 配置
type CookieConfig struct {
	File string `mapstructure:"file"`
}

// BlacklistConfig 黑名单配置
type BlacklistConfig struct {
	Enable   bool     `mapstructure:"enable"`
	Keywords []string `mapstructure:"keywords"`
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) error {
	Config = viper.New()

	// 设置配置文件路径
	if configPath != "" {
		Config.SetConfigFile(configPath)
	} else {
		// 默认配置文件路径
		Config.SetConfigName("config")
		Config.AddConfigPath(".")
		Config.AddConfigPath("./config")
		Config.AddConfigPath("./conf")
	}

	// 设置环境变量前缀
	Config.SetEnvPrefix("APP")
	Config.AutomaticEnv()

	// 读取配置文件
	if err := Config.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 确保必要的目录存在
	if err := ensureDirectories(); err != nil {
		return fmt.Errorf("创建必要目录失败: %w", err)
	}

	// 初始化日志
	logLevel := Config.GetString("app.log_level")
	logFile := Config.GetString("app.log_file")
	if err := InitLogger(logLevel, logFile); err != nil {
		return fmt.Errorf("初始化日志失败: %w", err)
	}

	Info("配置文件加载成功: ", Config.ConfigFileUsed())
	return nil
}

// ensureDirectories 确保必要目录存在
func ensureDirectories() error {
	dirs := []string{
		"data",
		"logs",
		"resources",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

// GetString 获取字符串配置
func GetString(key string) string {
	return Config.GetString(key)
}

// GetInt 获取整数配置
func GetInt(key string) int {
	return Config.GetInt(key)
}

// GetBool 获取布尔配置
func GetBool(key string) bool {
	return Config.GetBool(key)
}

// GetStringSlice 获取字符串切片配置
func GetStringSlice(key string) []string {
	return Config.GetStringSlice(key)
}

// GetMap 获取 map 配置
func GetMap(key string) map[string]interface{} {
	return Config.GetStringMap(key)
}
