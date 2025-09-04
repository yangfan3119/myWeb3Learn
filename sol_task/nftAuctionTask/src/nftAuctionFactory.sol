// SPDX-License-Identifier: UNLICENSED
pragma solidity >0.8.0 <0.9.0;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

contract nftAuctionFactory is Ownable, ReentrancyGuard {
    address private auctionTemplate;
    bytes4 private auctionInitSelector;

    address[] private allNftAuctions;
    uint256 private nftAuctionCount = 1;

    event NftAuctionCreated(address indexed auctionAddress, address indexed creator, uint256 auctionId);

    constructor() Ownable(msg.sender) {
    }

    function _setAuctionTemplate(address newImplementation, bytes4 newInitSelector) external onlyOwner {
        require(newImplementation != address(0), "nftAuctionFactory: new implementation is the zero address.");
        auctionTemplate = newImplementation; 
        auctionInitSelector = newInitSelector;
    }

    function _getAuctionTemplate() external view returns (address,bytes4) {
        return (auctionTemplate, auctionInitSelector);
    }

    function createNftAuction(address _owner, address feeAccount, uint256 feePercent) public nonReentrant returns (address) {
        // 使用 ERC1967Proxy 替代自定义代理逻辑
        bytes memory initData = abi.encodeWithSelector(auctionInitSelector, _owner, feeAccount, feePercent);
        ERC1967Proxy newAuctionProxy = new ERC1967Proxy(auctionTemplate, initData);
        address newAuctionAddress = address(newAuctionProxy);

        // require(auctionTemplate==getImplementation(newAuctionAddress), "create error auction ERC1967Proxy");

        // 4. 记录新拍卖合约地址
        allNftAuctions.push(newAuctionAddress);
        nftAuctionCount++;
        
        emit NftAuctionCreated(newAuctionAddress, msg.sender, nftAuctionCount);
        return newAuctionAddress;
    }

    function getAllNftAuctions() external view returns (address[] memory) {
        return allNftAuctions;
    }

    function getNftAuctionById(uint256 nftAuctionId) public view returns (address) {
        require((nftAuctionId < allNftAuctions.length)&&(nftAuctionId != 0),  "getNftAuctionById get Error Auction Id.");
        return allNftAuctions[nftAuctionId];
    }

    function getAuctionVersion(uint256 nftAuctionId) external returns(string memory){
        address auctionContract = getNftAuctionById(nftAuctionId);
        string memory version = "";
        // 1. 编码函数签名：getVersion() 的函数选择器（前4字节）
        // 函数签名为 "getVersion()"，通过 keccak256 哈希后取前4字节
        bytes4 auctionSelector = bytes4(keccak256("getVersion()"));

        // 2. 发起 call 调用：目标地址 + 函数选择器（无参数）
        (bool success, bytes memory returnData) = auctionContract.call(abi.encodePacked(auctionSelector));

        // 3. 验证调用成功，并解码返回数据（string 类型）
        if (success) {
            // 关键：将返回的 bytes 解码为 string
            version = abi.decode(returnData, (string));
        }
        return version;
    }
}