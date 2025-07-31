package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type BlogConfig struct {
	SqliteName  string        `mapstructure:"sqlite_name"`
	ServerPort  string        `mapstructure:"server_port"`
	ExpireHours time.Duration `mapstructure:"expire_hours"`
	JwtSecret   string        `mapstructure:"jwt_secret"`
}

var Cfg BlogConfig

func (c *BlogConfig) Load(path string) error {
	addr, _ := os.Getwd()
	fmt.Println("当前路径：", addr)
	// 设置配置文件路径和名称
	viper.SetConfigFile(path)

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		return err
	}

	// 解析配置到结构体
	if err := viper.Unmarshal(&c); err != nil {
		return err
	}
	c.ExpireHours *= time.Hour

	return nil
}
