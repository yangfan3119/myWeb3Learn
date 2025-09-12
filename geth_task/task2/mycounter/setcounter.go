package mycounter

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func SetMyCounter(rawurl string, deployPrivateKey string, contractAddr string, newCounter int64) error {
	client, err := ethclient.Dial(rawurl)
	if err != nil {
		log.Fatal(err)
	}
	// 创建私钥（在实际应用中，您应该使用更安全的方式来管理私钥）
	privateKey, err := crypto.HexToECDSA(deployPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 合约初始化参数设置
	contractABI, err := abi.JSON(strings.NewReader(MycounterABI))
	if err != nil {
		log.Fatal(err)
	}
	contract_setData, err := contractABI.Pack("setNumber", big.NewInt(newCounter))
	if err != nil {
		log.Fatal(err)
	}
	contractAddress := common.HexToAddress(contractAddr)
	chainID := getChainId(client)
	nonce := getNonce(client, fromAddress)
	gasPrice := getGasPrice(client)
	gasTipCap := getGasTipCap(client)
	gasFeeCap := new(big.Int).Add(gasPrice, big.NewInt(2000000))
	// 创建交易
	// tx := types.NewContractCreation(nonce, big.NewInt(0), 3000000, gasPrice, data)
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:    chainID,
		Nonce:      nonce,
		GasTipCap:  gasTipCap,
		GasFeeCap:  gasFeeCap,
		Gas:        uint64(3000000),
		To:         &contractAddress,
		Value:      &big.Int{},
		Data:       contract_setData,
		AccessList: types.AccessList{},
	})

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

	// 等待交易被挖矿
	receipt, err := waitForReceipt(client, signedTx.Hash())
	if err != nil {
		log.Fatal(err)
	}
	_ = receipt

	return err
}

func GetMyCounter(rawurl string, contractAddr string) int64 {
	client, err := ethclient.Dial(rawurl)
	if err != nil {
		log.Fatal(err)
	}
	// 合约初始化参数设置
	contractABI, err := abi.JSON(strings.NewReader(MycounterABI))
	if err != nil {
		log.Fatal(err)
	}
	contractAddress := common.HexToAddress(contractAddr)
	data, err := contractABI.Pack("number")
	if err != nil {
		log.Fatal("Failed to pack 'number()': ", err)
	}

	callMsg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}
	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatal(err)
	}
	var contractNumber *big.Int
	if err := contractABI.UnpackIntoInterface(&contractNumber, "number", result); err != nil {
		log.Fatalf("解码结果失败: %v", err)
	}

	// 8. 输出结果
	log.Println("合约中的number值为: ", contractNumber.String())

	return contractNumber.Int64()
}
