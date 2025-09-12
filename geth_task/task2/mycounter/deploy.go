package mycounter

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func getChainId(client *ethclient.Client) *big.Int {
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return chainID
}

func getNonce(client *ethclient.Client, fromAddress common.Address) uint64 {
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	return nonce
}

func getGasPrice(client *ethclient.Client) *big.Int {
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return gasPrice
}

func getGasTipCap(client *ethclient.Client) *big.Int {
	gasTipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return gasTipCap
}

func DeployMyCounter(rawurl string, deployPrivateKey string) string {
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
	data_constructor, err := contractABI.Pack("", big.NewInt(10))
	if err != nil {
		log.Fatal(err)
	}
	data := append(common.FromHex(MycounterBin), data_constructor...)

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
		To:         nil,
		Value:      &big.Int{},
		Data:       data,
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
	contractAddress := receipt.ContractAddress.Hex()
	log.Printf("Contract deployed at: %s\n", contractAddress)
	return contractAddress
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
