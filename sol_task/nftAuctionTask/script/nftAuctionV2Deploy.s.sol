//SPDX-License-Identifier: UNLICENSED
pragma solidity >0.8.0 <0.9.0;

import {Script, console} from "forge-std/Script.sol";
import {nftAuctionV2} from "../src/nftAuctionV2.sol";
import {nftAuctionFactory, ERC1967Proxy} from "../src/nftAuctionFactory.sol";
import {NftAuctionV1Deploy, nftAuctionV1} from "../script/NftAuctionV1Deploy.s.sol";

contract NftAuctionV2Deploy is Script {
    address public factoryAddress;
    nftAuctionFactory public factory;

    function DeployAuctionV2() public returns(address){
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        
        // 部署新版本 V2 逻辑合约
        vm.broadcast(deployerPrivateKey);
        nftAuctionV2 logicV2_Templat = new nftAuctionV2();

        bytes4 initSelector = nftAuctionV1.initialize.selector;
        console.log("nftAuctionV1 initialize selector: ");
        console.logBytes4(initSelector);

        console.log("nftAuctionV2 logic contract address: ", address(logicV2_Templat));

        // 获取工厂合约实例
        factoryAddress = address(vm.envAddress("NFT_AUCTION_FACTORY"));

        if (factoryAddress == address(0x1234567890abcdef1234567890abcdef12345678)){
            vm.broadcast(deployerPrivateKey);
            NftAuctionV1Deploy v1D = new NftAuctionV1Deploy();
            factoryAddress = v1D.DeployAuctionV1();
        }

        factory = nftAuctionFactory(factoryAddress);
        
        // 更新工厂合约的模板地址
        vm.broadcast(deployerPrivateKey);
        factory._setAuctionTemplate(address(logicV2_Templat), initSelector);

        (address template, bytes4 selector) = factory._getAuctionTemplate();
        require(template == address(logicV2_Templat), "template address not match");
        require(selector == initSelector, "initialize selector not match");
        
        // 升级工厂中已有的拍卖合约
        address[] memory allAuctions = factory.getAllNftAuctions();
        console.log("allAuctions length: ", allAuctions.length);

        bytes32 v2_template_UUID = logicV2_Templat.proxiableUUID();
        console.log("v2_template_UUID uuid :");
        console.logBytes32(v2_template_UUID);

        for (uint256 i = 1;i <= allAuctions.length; i++){
            address auctionAddress = allAuctions[i-1];
            console.log("auctionAddress : ", auctionAddress);
            console.log("before Upgrade V1, auctionVersion is: ", nftAuctionV1(auctionAddress).getVersion());

            // 升级到 V2
            nftAuctionV1(auctionAddress).upgradeToAndCall(address(logicV2_Templat), "");
            console.log("after Upgrade V1, auctionVersion is: ", nftAuctionV2(auctionAddress).getVersion());
        }

        return address(factory);
    }

    function run() external {
        DeployAuctionV2();
    }
}