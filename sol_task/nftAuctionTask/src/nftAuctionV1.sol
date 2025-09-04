// SPDX-License-Identifier: UNLICENSED
pragma solidity >0.8.0 <0.9.0;

import {IERC721} from "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {AggregatorV3Interface} from "lib/chainlink-brownie-contracts/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

contract nftAuctionV1 is Initializable, UUPSUpgradeable, OwnableUpgradeable, ReentrancyGuard {

    struct Auction {
        address seller;          // 卖家地址
        address nftAddress;     // NFT合约地址
        uint256 tokenId;        // NFT的tokenId

        uint256 startBid;       // 起始竞拍价
        uint256 startTime;      // 竞拍开始时间
        uint256 endTime;        // 竞拍结束时间
    
        address highestBidder;   // 最高出价者
        address highestBidToken; // 最高出价的币种地址
        uint256 highestBidAmount; // 最高出价金额
        uint256 highestBidUsd;     // 最高竞拍价
        
        bool ended;             // 竞拍是否结束    
    }

    mapping(address token => AggregatorV3Interface priceFeed) public priceFeeds; // 价格预言机映射
    mapping(uint256 auctionId => Auction) public auctionMap; // 通过ID映射拍卖
    uint256 public auctionCount; // 拍卖计数器
    address public feeAccount; // 平台手续费接收地址
    uint256 public feePercent; // 平台手续费百分比

    uint256[50] private __gap; // OpenZeppelin 推荐的存储填充

    function initialize(address _admin, address _feeAccount, uint256 _feePercent) public initializer {
        require(_feeAccount != address(0), "nftAuctionV1: fee account is the zero address.");
        require(_feePercent <= 100, "nftAuctionV1: fee percent must be between 0 and 100.");
        auctionCount = 1;
        feeAccount = _feeAccount;
        feePercent = _feePercent;

        __Ownable_init(_admin);
        __UUPSUpgradeable_init();
    }

    function getVersion() external pure virtual returns (string memory) {
        return "v1.0.0";
    }

    function setPriceFeed(address token, address priceFeed) external onlyOwner {
        priceFeeds[token] = AggregatorV3Interface(priceFeed);
    }

    function getPriceUSD(address token, uint256 amount) public view returns (uint256) {
        AggregatorV3Interface priceFeed = priceFeeds[token];

        (, int256 price, , , ) = priceFeed.latestRoundData();
        uint8 decimals = priceFeed.decimals();

        // 计算美元价值，调整小数位
        return (amount * uint256(price)) / (10 ** decimals);
    }

    function createAuction(
        address nftAddress,
        uint256 tokenId,
        uint256 startBid,
        uint256 startTime,
        uint256 endTime
    ) external nonReentrant returns(uint256){
        require(endTime > startTime, "nftAuctionV1: endTime must be greater than startTime.");
        require(startTime >= block.timestamp, "nftAuctionV1: startTime must be greater than current time.");

        // 转移NFT到合约地址
        IERC721(nftAddress).transferFrom(msg.sender, address(this), tokenId);

        uint256 auctionId = auctionCount;
        auctionMap[auctionId] = Auction({
            seller: msg.sender,
            nftAddress: nftAddress,
            tokenId: tokenId,
            startBid: startBid,
            startTime: startTime,
            endTime: endTime,
            highestBidder: address(0),
            highestBidToken: address(0),
            highestBidAmount: 0,
            highestBidUsd: 0,
            ended: false
        });

        auctionCount++;

        return auctionId;
    }

    function placeBidETH(uint256 auctionId) external payable nonReentrant {
        Auction storage auction = auctionMap[auctionId];

        require(auction.ended == false, "nftAuctionV1: auction already ended.");
        require(block.timestamp >= auction.startTime, "nftAuctionV1: auction not started.");
        require(block.timestamp <= auction.endTime, "nftAuctionV1: auction ended.");
        
        uint256 bidUsd = getPriceUSD(address(0), msg.value);
        require(bidUsd >= auction.startBid, "nftAuctionV1: bid must be greater than startBid.");
        require(bidUsd > auction.highestBidUsd, "nftAuctionV1: bid must be greater than highestBidUsd.");

        // 退还之前的最高出价者
        highestBidSendBack(auction);

        // 更新最高出价信息
        auction.highestBidder = msg.sender;
        auction.highestBidToken = address(0); // ETH
        auction.highestBidAmount = msg.value;
        auction.highestBidUsd = bidUsd;
    }

    function getHighestBidInfo(uint256 auctionId)public view returns(address,address,uint256,uint256){
        Auction storage auction = auctionMap[auctionId];
        return (auction.highestBidder, auction.highestBidToken, auction.highestBidAmount, auction.highestBidUsd);
    }

    function placeBidERC20(uint256 auctionId, address bidToken, uint256 bidAmount) external nonReentrant{
        Auction storage auction = auctionMap[auctionId];

        require(auction.ended == false, "nftAuctionV1: auction already ended.");
        require(block.timestamp >= auction.startTime, "nftAuctionV1: auction not started.");
        require(block.timestamp <= auction.endTime, "nftAuctionV1: auction ended.");

        uint256 bidUsd = getPriceUSD(bidToken, bidAmount);
        require(bidUsd >= auction.startBid, "nftAuctionV1: bid must be greater than startBid.");
        require(bidUsd > auction.highestBidUsd, "nftAuctionV1: bid must be greater than highestBidUsd.");

        // 从竞拍者转移代币到合约
        bool succ = IERC20(bidToken).transferFrom(msg.sender, address(this), bidAmount);
        require(succ, "nftAuctionV1: placeBidERC20 ERC20 transferFrom msg.sender to address(this) failed.");

        // 退还之前的最高出价者
        highestBidSendBack(auction);

        // 更新最高出价信息
        auction.highestBidder = msg.sender;
        auction.highestBidToken = bidToken;
        auction.highestBidAmount = bidAmount;
        auction.highestBidUsd = bidUsd;
    }

    function endAuction(uint256 auctionId) external {
        Auction storage auction = auctionMap[auctionId];

        require(auction.ended == false, "nftAuctionV1: auction already ended.");
        require(block.timestamp > auction.endTime, "nftAuctionV1: auction not yet ended.");
        require(msg.sender == auction.seller || msg.sender == owner(), "nftAuctionV1: only seller or owner can end the auction.");

        auction.ended = true;
        bool success;

        if (auction.highestBidder != address(0)) {
            // 计算平台手续费
            uint256 feeAmount = (auction.highestBidAmount * feePercent) / 100;

            // 转移手续费到平台账户
            if (auction.highestBidToken == address(0)) {
                // ETH
                payable(feeAccount).transfer(feeAmount);
                // 转移剩余金额给卖家
                payable(auction.seller).transfer(auction.highestBidAmount - feeAmount);
            } else {
                IERC20(auction.highestBidToken).approve(address(this), auction.highestBidAmount);
                // ERC20
                success = IERC20(auction.highestBidToken).transferFrom(address(this), feeAccount, feeAmount);
                require(success, "nftAuctionV1: IERC20 transferFrom address(this) to feeAccount failed.");
                // 转移剩余金额给卖家
                success = IERC20(auction.highestBidToken).transferFrom(address(this), auction.seller, auction.highestBidAmount - feeAmount);
                require(success, "nftAuctionV1: IERC20 transferFrom address(this) to feeAccount failed.");
            }

            // 转移NFT给最高出价者
            IERC721(auction.nftAddress).transferFrom(address(this), auction.highestBidder, auction.tokenId);
        } else {
            // 无人出价，退还NFT给卖家
            IERC721(auction.nftAddress).transferFrom(address(this), auction.seller, auction.tokenId);
        }
    }

    function highestBidSendBack(Auction memory auction) internal {
        // 退还之前的最高出价者,首先判断是否有最高价，有则退款
        if (auction.highestBidder != address(0)) {
            if (auction.highestBidToken == address(0)) {
                payable(auction.highestBidder).transfer(auction.highestBidAmount);
            } else {
                bool succ = IERC20(auction.highestBidToken).transfer(auction.highestBidder, auction.highestBidAmount);
                require(succ, "nftAuctionV1: placeBidERC20 ERC20 transfer auction.highestBidder failed.");
            }
        }
    }

    function getAuction(uint256 auctionId) external view returns (Auction memory) {
        return auctionMap[auctionId];
    }

    function withdraw() external onlyOwner {
        payable(owner()).transfer(address(this).balance);
    }

    function owner() public view override(OwnableUpgradeable) returns (address) {
        return super.owner();
    }

    function getFeeAccount() public view returns (address) {
        return feeAccount;
    }

    function getFeePercent() public view returns (uint256) {
        return feePercent;
    }

    function _authorizeUpgrade(address newImplementation) internal virtual override onlyOwner {
    }
    
    // 删除自定义的proxiableUUID函数，使用UUPSUpgradeable的默认实现

    // function upgradeToAndCall(address newImplementation, bytes calldata data) external onlyOwner {
    //     _authorizeUpgrade(newImplementation); // 权限验证
    //     _upgradeToAndCallUUPS(newImplementation, data, false); // UUPS 标准升级逻辑
    // }
}