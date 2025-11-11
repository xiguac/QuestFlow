// Package config 负责加载和管理应用配置
package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// DBConfig 封装了数据库配置项
type DBConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Params   string `mapstructure:"params"`
}

// DSN 方法根据配置动态生成数据库连接字符串
func (db *DBConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		db.User, db.Password, db.Host, db.Port, db.DBName, db.Params)
}

// AppConfig 是应用配置的结构体
type AppConfig struct {
	App struct {
		Port string `mapstructure:"port"`
		Mode string `mapstructure:"mode"`
	} `mapstructure:"app"`
	Database DBConfig `mapstructure:"database"`
	Redis    struct {
		Addr                string `mapstructure:"addr"`
		Password            string `mapstructure:"password"`
		DB                  int    `mapstructure:"db"`
		SubmissionStreamKey string `mapstructure:"submission_stream_key"`
		SubmissionGroupName string `mapstructure:"submission_group_name"`
	} `mapstructure:"redis"`
	JWT struct {
		Secret      string `mapstructure:"secret"`
		Issuer      string `mapstructure:"issuer"`
		ExpireHours int    `mapstructure:"expire_hours"`
	} `mapstructure:"jwt"`
}

// Cfg 是一个全局的配置实例
var Cfg AppConfig

// Init 加载配置文件到全局变量 Cfg
func Init(configPath string) {
	// --- 新增环境变量处理逻辑 ---
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 允许 Viper 从环境变量读取值
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		panic(fmt.Errorf("unable to decode into struct, %w", err))
	}
}
