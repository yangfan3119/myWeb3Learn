package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

/*
使用 Sepolia 测试网络实现基础的区块链交互，包括查询区块和发送交易。
具体任务:
1.环境搭建
- 安装必要的开发工具，如 Go 语言环境、 go-ethereum 库。
- 注册 Infura 账户，获取 Sepolia 测试网络的 API Key。
2.查询区块
- 编写 Go 代码，使用 ethclient 连接到 Sepolia 测试网络。
- 实现查询指定区块号的区块信息，包括区块的哈希、时间戳、交易数量等。
- 输出查询结果到控制台。
3.发送交易
- 准备一个 Sepolia 测试网络的以太坊账户，并获取其私钥。
- 编写 Go 代码，使用 ethclient 连接到 Sepolia 测试网络。
- 构造一笔简单的以太币转账交易，指定发送方、接收方和转账金额。
- 对交易进行签名，并将签名后的交易发送到网络。
- 输出交易的哈希值。
*/
type BlogConfig struct {
	RawUrl             string `mapstructure:"raw_url"`
	Account1           string `mapstructure:"account1"`
	Account6           string `mapstructure:"account6"`
	Account6PrivateKey string `mapstructure:"account6_private_key"`
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
func main() {
	var Cfg BlogConfig
	Cfg.Load("../dev.config.yaml")

	client, err := ethclient.Dial(Cfg.RawUrl)
	if err != nil {
		log.Fatalln("ethclient err. ", err)
	}

	privateKey, err := crypto.HexToECDSA(Cfg.Account6PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	log.Println("fromAddress is equal account6. ", fromAddress.String() == Cfg.Account6)
	toAddress := common.HexToAddress(Cfg.Account1)
	log.Println("toAddress is equal account1. ", toAddress.String() == Cfg.Account1)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	gasTipCap, _ := client.SuggestGasTipCap(context.Background())
	gasFeeCap := new(big.Int).Add(gasPrice, big.NewInt(2000000))
	value := new(big.Int).Mul(big.NewInt(1), big.NewInt(1e13))
	// 创建交易
	// tx := types.NewContractCreation(nonce, big.NewInt(0), 3000000, gasPrice, data)
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:    chainID,
		Nonce:      nonce,
		GasTipCap:  gasTipCap,
		GasFeeCap:  gasFeeCap,
		Gas:        uint64(300000),
		To:         &toAddress,
		Value:      value,
		Data:       nil,
		AccessList: types.AccessList{},
	})
	// tx := types.NewTransaction(nonce, contractAddress, big.NewInt(0), 300000, gasPrice, input)
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())
	receipt, err := waitForReceipt(client, signedTx.Hash())
	if err != nil {
		log.Fatal(err)
	}
	if receipt.Status == 0 {
		log.Fatal("Transaction failed. (reverted)")
	} else {
		log.Println("Transaction success.")
	}
}
func waitForReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			return receipt, nil
		}
		if err != ethereum.NotFound {
			return nil, err
		}
		// 等待一段时间后再次查询
		time.Sleep(1 * time.Second)
	}
}
