// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/Console.sol";
import {TestERC20} from "../test/mock/mERC20.sol";
import {task02_nft} from "../test/mock/mERC721.sol";
import {MockV3Aggregator} from "../test/mock/mV3Aggregator.sol";

contract HelperConfig is Script{
    
    struct UserInfo {
        address userAddress;

        address eT1Token;
        uint256 eT1Balance;
        address mET1toUSDPriceFeed;

        uint256 ethBalance;
        address mETHtoUSDPriceFeed;

        address nftAddress;
        uint256 nftTokenId;
    }
    
    UserInfo public user1;
    UserInfo public user2;
    UserInfo public user3;

    uint256 public constant STARTING_USER_BALANCE = 100 ether;

    uint256 public constant INITIAL_ERC20_BALANCE = 100000 * 10**18;
    int256 public constant ANVIL_PRICEFEED_ETH_USD = 200e18;
    int256 public constant ANVIL_PRICEFEED_ET1_USD = 1e17;

    string public constant NFT_TEST_URI = "https://prominent-moccasin-thrush.myfilebase.com/ipfs/QmSh6DUx9dSgtcA6Tn9GNVqCjsxWt9ptz1BHns94Udhx37";
    address public constant DEPLOY_USER = address(10);    

    constructor() {

    }

    function getUser(uint8 index) public view returns(UserInfo memory){
        if (index == 1){
            return user1;
        }
        if (index == 2){
            return user2;
        }
        return user3;
    }

    function dataInit() external {
        // 创建用户,使用anvil提供的测试地址01，02，03.
        // 01: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
        // 02: 0x70997970C51812dc3A010C7d01b50e0d17dc79C8
        // 03: 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC
        user1.userAddress = address(0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266);
        user2.userAddress = address(0x70997970C51812dc3A010C7d01b50e0d17dc79C8);
        user3.userAddress = address(0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC);
        
        vm.deal(user1.userAddress, STARTING_USER_BALANCE);
        vm.deal(user2.userAddress, STARTING_USER_BALANCE);
        vm.deal(user3.userAddress, STARTING_USER_BALANCE);

        vm.startBroadcast(DEPLOY_USER);
        console.log("Deploy user address: ", DEPLOY_USER);
        // 创建ERC20代币 eT1
        TestERC20 eT1 = new TestERC20("myTestErc20", "ET1", DEPLOY_USER);
        eT1.mint(DEPLOY_USER, INITIAL_ERC20_BALANCE*10);

        address eT1Address = address(eT1);
        console.log("eT1 address: ", eT1Address);
        console.log("eT1 Deploy user address: ", eT1.getOwner());
        
        // 给user123转账eT1代币
        eT1.approve(DEPLOY_USER, INITIAL_ERC20_BALANCE*3);
        eT1.transferFrom(DEPLOY_USER, user1.userAddress, INITIAL_ERC20_BALANCE);
        eT1.transferFrom(DEPLOY_USER, user2.userAddress, INITIAL_ERC20_BALANCE);
        eT1.transferFrom(DEPLOY_USER, user3.userAddress, INITIAL_ERC20_BALANCE);

        // 创建NFT合约 
        task02_nft nft = new task02_nft("myTestNft", "MFT", DEPLOY_USER);

        address nftAddress = address(nft);
        console.log("nft address: ", nftAddress);
        console.log("nft Deploy user address: ", eT1.getOwner());

        // 合约所有者，给自己铸造NFT
        uint256 nftTokenId = nft.mintNFT(DEPLOY_USER, NFT_TEST_URI);
        console.log("nft tokenId 1 owner: ", nft.ownerOf(nftTokenId));

        nft.approve(DEPLOY_USER, nftTokenId);

        // 铸造的NFT转账给user1
        nft.transferNFT(DEPLOY_USER, user3.userAddress, nftTokenId);

        // 创建ETH/USD价格预言机合约, 创建eT1/USD价格预言机合约
        MockV3Aggregator ethUsdPriceFeed = new MockV3Aggregator(18, ANVIL_PRICEFEED_ETH_USD);
        MockV3Aggregator eT1UsdPriceFeed = new MockV3Aggregator(18, ANVIL_PRICEFEED_ET1_USD);

        vm.stopBroadcast();

        address ethUsdPriceFeedAddress = address(ethUsdPriceFeed);
        console.log("ethUsdPriceFeed address: ", ethUsdPriceFeedAddress);

        address eT1UsdPriceFeedAddress = address(eT1UsdPriceFeed);
        console.log("eT1UsdPriceFeed address: ", eT1UsdPriceFeedAddress);
        

        // 设置用户信息
        user1.eT1Token = eT1Address;
        user1.eT1Balance = INITIAL_ERC20_BALANCE;
        user1.ethBalance = STARTING_USER_BALANCE;
        user1.mET1toUSDPriceFeed = eT1UsdPriceFeedAddress;
        user1.mETHtoUSDPriceFeed = ethUsdPriceFeedAddress;
        user1.nftAddress = nftAddress;
        user1.nftTokenId = nftTokenId;

        user2.eT1Token = eT1Address;
        user2.eT1Balance = INITIAL_ERC20_BALANCE;
        user2.ethBalance = STARTING_USER_BALANCE;
        user2.mET1toUSDPriceFeed = eT1UsdPriceFeedAddress;
        user2.mETHtoUSDPriceFeed = ethUsdPriceFeedAddress;
        user2.nftAddress = nftAddress;
        user2.nftTokenId = nftTokenId;

        user3.eT1Token = eT1Address;
        user3.eT1Balance = INITIAL_ERC20_BALANCE;
        user3.ethBalance = STARTING_USER_BALANCE;
        user3.mET1toUSDPriceFeed = eT1UsdPriceFeedAddress;
        user3.mETHtoUSDPriceFeed = ethUsdPriceFeedAddress;
        user3.nftAddress = nftAddress;
        user3.nftTokenId = nftTokenId;
    }
}