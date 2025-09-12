package main

import (
	"fmt"
	"geth_task/task2/mycounter"
	"log"
	"os"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

/*
使用 abigen 工具自动生成 Go 绑定代码，用于与 Sepolia 测试网络上的智能合约进行交互。
具体任务
1.编写智能合约
- 使用 Solidity 编写一个简单的智能合约，例如一个计数器合约。
- 编译智能合约，生成 ABI 和字节码文件。
2.使用 abigen 生成 Go 绑定代码
- 安装 abigen 工具。
- 使用 abigen 工具根据 ABI 和字节码文件生成 Go 绑定代码。
3.使用生成的 Go 绑定代码与合约交互
- 编写 Go 代码，使用生成的 Go 绑定代码连接到 Sepolia 测试网络上的智能合约。
- 调用合约的方法，例如增加计数器的值。
- 输出调用结果。
*/

type BlogConfig struct {
	RawUrl             string `mapstructure:"raw_url"`
	Account1           string `mapstructure:"account1"`
	Account6           string `mapstructure:"account6"`
	Account6PrivateKey string `mapstructure:"account6_private_key"`
	ContractAddress    string `mapstructure:"contract_address"`
}

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

	return nil
}

func (c *BlogConfig) Save(path string) error {
	// 清除现有配置，确保删除的字段会从配置文件中移除
	viper.Reset()

	// 设置配置文件路径
	viper.SetConfigFile(path)

	// 将结构体转换为map，便于设置到viper
	configMap := make(map[string]interface{})
	if err := mapstructure.Decode(c, &configMap); err != nil {
		return fmt.Errorf("结构体转换失败: %v", err)
	}

	// 将结构体字段设置到viper
	for key, value := range configMap {
		viper.Set(key, value)
	}

	// 尝试写入配置文件
	err := viper.WriteConfig()
	if err != nil {
		log.Println("写入配置文件失败: ", err)
	} else {
		log.Println("配置文件已更新")
	}

	return nil
}

var cPath = "./dev.config.yaml"

// 0xAfa5E2E91C255022aCe7C7D734758F8a59Ee72a9
// 0x31a207EFa91B140046eFc0D5bc3418451e9D45f2

func main() {
	var Cfg BlogConfig
	Cfg.Load(cPath)

	if Cfg.ContractAddress == "" {
		Cfg.ContractAddress = mycounter.DeployMyCounter(Cfg.RawUrl, Cfg.Account6PrivateKey)
		Cfg.Save(cPath)
	}
	// 获取当前的counter值，然后+13重新写入，再重新读取验证是否正确
	cNumber := mycounter.GetMyCounter(Cfg.RawUrl, Cfg.ContractAddress)
	log.Println("当前合约中的number值为: ", cNumber)
	cNumber += 13
	if err := mycounter.SetMyCounter(Cfg.RawUrl, Cfg.Account6PrivateKey, Cfg.ContractAddress, cNumber); err != nil {
		log.Fatal(err)
	}
	cNumber = mycounter.GetMyCounter(Cfg.RawUrl, Cfg.ContractAddress)
	log.Println("当前合约中的number值为: (预期与上次获取值相差13)", cNumber)
}
