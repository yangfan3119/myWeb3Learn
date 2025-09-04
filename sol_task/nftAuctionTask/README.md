## Foundry


### 1. lib下载
```bash
forge install OpenZeppelin/openzeppelin-contracts
forge install OpenZeppelin/openzeppelin-contracts-upgradeable
forge install smartcontractkit/chainlink-brownie-contracts
forge install foundry-rs/forge-std
```

### 2. 部署合约
```bash
# .env文件中需要添加PRIVATE_KEY和PRC_URL
# 编译
forge build
# 测试
forge test
# 部署
forge script nftAuctionV1Deploy.s.sol --rpc-url $PRC_URL --private-key $PRIVATE_KEY --broadcast 
# 升级，需要把之前部署的Factory地址写入env中。NFT_AUCTION_FACTORY
# NFT_AUCTION_FACTORY=0x1234567890abcdef1234567890abcdef12345678默认假地址，用于测试时可读取Test创建的虚拟Factory地址
forge script nftAuctionV2Deploy.s.sol --rpc-url $PRC_URL --private-key $PRIVATE_KEY --broadcast 
```




