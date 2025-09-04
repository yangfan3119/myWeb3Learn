// SPDX-License-Identifier: UNLICENSED
pragma solidity >0.8.0 <0.9.0;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/Console.sol";
import {nftAuctionV1} from "../src/nftAuctionV1.sol";
import {nftAuctionFactory} from "../src/nftAuctionFactory.sol";
import {ERC1967Utils} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Utils.sol";

contract NftAuctionV1Deploy is Script {
    function DeployAuctionV1() public returns(address){
        console.log("block.chainid is : ", block.chainid);
        
        // 获取部署者私钥
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        vm.startBroadcast(deployerPrivateKey);

        // 1. 部署V1逻辑合约
        nftAuctionV1 logicV1_Templat = new nftAuctionV1();
        bytes4 initSelector = nftAuctionV1.initialize.selector;
        console.log("nftAuctionV1 initialize selector: ");
        console.logBytes4(initSelector);
        console.log("nftAuctionV1 logic contract address: ", address(logicV1_Templat));

        // 2. 部署工厂合约
        nftAuctionFactory factory = new nftAuctionFactory();
        console.log("NFT Auction Factory contract address: ", address(factory));

        // 3. 设置工厂合约中的实现地址
        factory._setAuctionTemplate(address(logicV1_Templat), initSelector);
        (address template, bytes4 selector) = factory._getAuctionTemplate();
        require(template == address(logicV1_Templat), "template address not match");
        require(selector == initSelector, "initialize selector not match");

        // 4. 预创建一个拍卖合约作为测试
        address testOwner = msg.sender;
        address feeAccount = msg.sender;
        uint256 feePercent = 2; // 2%手续费
        
        address newAuction = factory.createNftAuction(
            testOwner,
            feeAccount,
            feePercent
        );
        console.log("create New NftAuction contract address : ", newAuction);

        // 5. 验证初始化是否成功
        nftAuctionV1 auction = nftAuctionV1(newAuction);
        require(auction.owner() == testOwner, "create failed: owner not match");
        require(auction.feeAccount() == feeAccount, "create failed: fee account not match");
        require(auction.feePercent() == feePercent, "create failed: fee percent not match");

        console.log("Deploy contract is Done!");
        vm.stopBroadcast();

        return address(factory);
    }

    function run() external{
        DeployAuctionV1();
    }
}