// SPDX-License-Identifier: UNLICENSED
pragma solidity >0.8.0 <0.9.0;

import {Test} from "forge-std/Test.sol";
import {console} from "forge-std/Console.sol";
import {HelperConfig} from "../script/HelperConfig.s.sol";
import {nftAuctionFactory} from "../src/nftAuctionFactory.sol";
import {nftAuctionV1, IERC721, IERC20} from "../src/nftAuctionV1.sol";
import {nftAuctionV2} from "../src/nftAuctionV2.sol";
import {NftAuctionV1Deploy} from "../script/nftAuctionV1Deploy.s.sol";

contract nftAuctionV1Test is Test {
    nftAuctionV1 private v1Auction;
    HelperConfig private hc;
    HelperConfig.UserInfo public user1;
    HelperConfig.UserInfo public user2;
    HelperConfig.UserInfo public user3;

    address public constant USER1 = address(0x1);
    address public constant FEE_ACCOUNT = address(0x100);
    uint256 public constant FEE_PERCENT = 3;

    
    address public constant USER2 = address(0x2);
    address public constant FEE_ACCOUNT_2 = address(0x200);
    uint256 public constant FEE_PERCENT_2 = 5;

    uint256 private constant NEW_STYLE = 3; // 1: v1拍卖合约, 2：Factory生成v1合约，, 3: Deploy 获取Factory ，然后创建v1

    function setUp() public {
        // Setup code if needed
        if (NEW_STYLE == 1){
            v1Auction = nftAuctionV1(newNftAuctionV1());
            v1Auction.initialize(USER1, FEE_ACCOUNT, FEE_PERCENT);
        }else if (NEW_STYLE == 2){
            v1Auction = nftAuctionV2(newNftAuctionFactory());
        }else {
            v1Auction = nftAuctionV2(newNftAuctionV1Deploy());
        }
        
        vm.startBroadcast();
        hc = new HelperConfig();
        vm.stopBroadcast();

        hc.dataInit();

        HelperConfig.UserInfo memory userData = hc.getUser(1);
        user1 = userData;
        assertEq(user1.userAddress, 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266, "user1 is not Match.");

        userData = hc.getUser(2);
        user2 = userData;
        assertEq(user2.userAddress, 0x70997970C51812dc3A010C7d01b50e0d17dc79C8, "user2 is not Match.");

        userData = hc.getUser(3);
        user3 = userData;
        assertEq(user3.userAddress, 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC, "user3 is not Match.");
    }

    function newNftAuctionV1() public returns(address) {
        vm.startBroadcast();
        nftAuctionV1 v1 = new nftAuctionV1();
        vm.stopBroadcast();
        return address(v1);
    }

    function newNftAuctionFactory() public returns(address) {
        vm.startBroadcast();
        nftAuctionFactory nftAuctionFacory = new nftAuctionFactory();
        nftAuctionV1 template = new nftAuctionV1();
        nftAuctionFacory._setAuctionTemplate(address(template), template.initialize.selector);
        nftAuctionFacory.createNftAuction(USER2, FEE_ACCOUNT_2, FEE_PERCENT_2);
        nftAuctionFacory.createNftAuction(USER2, FEE_ACCOUNT_2, FEE_PERCENT_2);
        nftAuctionFacory.createNftAuction(USER2, FEE_ACCOUNT_2, FEE_PERCENT_2);

        vm.stopBroadcast();
        return nftAuctionFacory.getNftAuctionById(1);
    }

    function newNftAuctionV1Deploy() public returns(address) {
        vm.startBroadcast();
        NftAuctionV1Deploy v1d = new NftAuctionV1Deploy();
        vm.stopBroadcast();

        nftAuctionFactory nftAuctionFacory = nftAuctionFactory(v1d.DeployAuctionV1());
        
        vm.startBroadcast();
        nftAuctionFacory.createNftAuction(USER2, FEE_ACCOUNT_2, FEE_PERCENT_2);
        nftAuctionFacory.createNftAuction(USER2, FEE_ACCOUNT_2, FEE_PERCENT_2);
        vm.stopBroadcast();

        string memory version = nftAuctionFacory.getAuctionVersion(1);
        console.log("newNftAuctionV1Deploy auction 1 version is:");
        console.logString(version);

        return nftAuctionFacory.getNftAuctionById(1);
    }

    function getAdmin() public pure returns(address){
        if(NEW_STYLE==1){
            return USER1;
        } else {
            return USER2;
        }
    }

    function getFeeAccount() public pure returns(address){
        if(NEW_STYLE==1){
            return FEE_ACCOUNT;
        } else {
            return FEE_ACCOUNT_2;
        }
    }

    function getFeePercent() public pure returns(uint256){
        if(NEW_STYLE==1){
            return FEE_PERCENT;
        } else {
            return FEE_PERCENT_2;
        }
    }

    function testV1AuctionVersion() public view{
        assertEq("v1.0.0", v1Auction.getVersion(), "nftAuction version is v1.0.0.");
    }

    function testInitializeV1Auction() public view {
        address auctionOwner = v1Auction.owner();
        assertEq(auctionOwner, getAdmin(), "v1Auction owner is USER1.");

        address feeAccount = v1Auction.feeAccount();
        assertEq(feeAccount, getFeeAccount(), "v1Auction feeAccount is FEE_ACCOUNT_2");

        uint256 feePercent = v1Auction.feePercent();
        assertEq(feePercent, getFeePercent(), "v1Auction feePercent is FEE_PERCENT_2");
    }

    function testV1SetPriceFeedAndGetPriceUSD() public {
        // 设置ETH和ET1的喂价器预言机
        vm.startPrank(getAdmin());
        v1Auction.setPriceFeed(address(0), user1.mETHtoUSDPriceFeed);
        v1Auction.setPriceFeed(user1.eT1Token,user1.mET1toUSDPriceFeed);
        vm.stopPrank();

        // 计算ETH价格是否正确
        uint256 ethPrice = v1Auction.getPriceUSD(address(0), 1);
        console.log("ETH to USD priceFeed is: ", hc.ANVIL_PRICEFEED_ETH_USD());
        assertEq(ethPrice, 200, "ETH to USD priceFeed is wrong");

        // 计算ET1价格是否正确
        uint256 et1Price = v1Auction.getPriceUSD(user2.eT1Token, 100);
        console.log("ETH to USD priceFeed is: ", hc.ANVIL_PRICEFEED_ETH_USD());
        assertEq(et1Price, 10, "ETH to USD priceFeed is wrong");
    }

    modifier auctionPrepare() {
        // 设置ETH和ET1的喂价器预言机
        vm.startPrank(getAdmin());
        v1Auction.setPriceFeed(address(0), user1.mETHtoUSDPriceFeed);
        v1Auction.setPriceFeed(user1.eT1Token,user1.mET1toUSDPriceFeed);
        vm.stopPrank();

        vm.warp(1000);
        uint256 startTime = block.timestamp + 10;
        uint256 endTime = block.timestamp + 40;

        address nftOwner = IERC721(user3.nftAddress).ownerOf(user3.nftTokenId);
        console.log("Auction create before NFT belong of:", nftOwner);

        vm.prank(user3.userAddress);
        IERC721(user3.nftAddress).approve(address(v1Auction), user3.nftTokenId);

        address approved = IERC721(user3.nftAddress).getApproved(user3.nftTokenId);
        console.log("IERC721(nftAddress) approved address: ", approved);

        vm.prank(user3.userAddress);
        uint256 auctionId = v1Auction.createAuction(
            user3.nftAddress, 
            user3.nftTokenId, 
            1, startTime, endTime);
        assertEq(auctionId, 1, "create an auction succ.");

        nftOwner = IERC721(user3.nftAddress).ownerOf(user3.nftTokenId);
        console.log("Auction create after NFT belong of:", nftOwner);
        _;
    }

    function testV1AuctionCreateStatus() public auctionPrepare{
        address nftOwner = IERC721(user3.nftAddress).ownerOf(user3.nftTokenId);
        assertEq(nftOwner, address(v1Auction), "IERC721 nft has transfer Auction contract.");
    }

    function testWarpTo30SecondsLater() public {
        // 1. 获取当前区块链时间戳
        uint256 currentTime = block.timestamp;
        console.log("time of now::", currentTime);

        // 2. 计算 30 秒后的时间戳
        uint256 targetTime = currentTime + 30; // 当前时间 + 30 秒
        console.log("time after 30s:", targetTime);

        // 3. 使用 vm.warp() 跳转到目标时间
        vm.warp(targetTime);

        // 4. 验证时间是否正确跳转
        assertEq(block.timestamp, targetTime, "time is not warp 30s.");
        console.log("time after 30s:", block.timestamp);
    }

    function testV1AuctionRevertBidBeforeStart() public auctionPrepare{
        uint256 warpTime = block.timestamp - 30;
        console.log("block.timestamp is", block.timestamp);
        console.log("warpTime is", warpTime);
        vm.warp(warpTime);
        assertEq(block.timestamp, warpTime, "time is not warp.");

        vm.prank(user1.userAddress);
        vm.expectRevert(bytes("nftAuctionV1: auction not started."));
        v1Auction.placeBidETH{value: 0.1 ether}(user1.nftTokenId);
    }

    function testV1AuctionUser1BidETHCheckStatus() public auctionPrepare{
        vm.warp(block.timestamp + 30);

        uint256 auctionId = 1;
        uint256 user1Balance = payable(user1.userAddress).balance;
        console.log("user1Balance before bid: ", user1Balance);

        vm.prank(user1.userAddress);
        v1Auction.placeBidETH{value: 0.1 ether}(auctionId);

        user1Balance = payable(user1.userAddress).balance;
        console.log("user1Balance after bid: ", user1Balance);

        (address bidder,address bidToken, uint256 bidAmount,uint256 bidUsd) = v1Auction.getHighestBidInfo(auctionId);
        console.log("bidder address: ", bidder);
        console.log("bidToken address: ", bidToken);
        console.log("bidAmount : ", bidAmount);
        console.log("bidUsd : ", bidUsd);
    }

    function testV1AuctionBidSendBackETH() public auctionPrepare{
        vm.warp(block.timestamp + 30);

        uint256 auctionId = 1;
        uint256 user1Balance = payable(user1.userAddress).balance;
        console.log("user1Balance before bid: ", user1Balance);

        vm.prank(user1.userAddress);
        v1Auction.placeBidETH{value: 0.1 ether}(auctionId);
        console.log("user1Balance after bid: ", payable(user1.userAddress).balance);

        vm.prank(user2.userAddress);
        v1Auction.placeBidETH{value: 0.2 ether}(auctionId);

        console.log("after user2 Bid, send back user1 ETH, user1 balance : ", payable(user1.userAddress).balance);
        assertEq(user1Balance, payable(user1.userAddress).balance, "nftAuctionV1 has not send back user1.");

        (address bidder,address bidToken, uint256 bidAmount,uint256 bidUsd) = v1Auction.getHighestBidInfo(auctionId);
        console.log("bidder address: ", bidder);
        console.log("bidToken address: ", bidToken);
        console.log("bidAmount : ", bidAmount);
        console.log("bidUsd : ", bidUsd);
    }

    function testV1AuctionUser1BidERC20CheckStatus() public auctionPrepare{
        vm.warp(block.timestamp + 30);

        uint256 auctionId = 1;
        uint256 user1BalanceErc20 = IERC20(user1.eT1Token).balanceOf(user1.userAddress);
        console.log("user1BalanceErc20 before bid: ", user1BalanceErc20);

        uint256 user1_bidprice_ERC20 = 100 * 10 **18;
        vm.prank(user1.userAddress);
        IERC20(user1.eT1Token).approve(address(v1Auction), user1_bidprice_ERC20);

        vm.prank(user1.userAddress);
        v1Auction.placeBidERC20(auctionId, user1.eT1Token, user1_bidprice_ERC20);

        user1BalanceErc20 = IERC20(user1.eT1Token).balanceOf(user1.userAddress);
        console.log("user1BalanceErc20 after bid: ", user1BalanceErc20);

        (address bidder,address bidToken, uint256 bidAmount,uint256 bidUsd) = v1Auction.getHighestBidInfo(auctionId);
        console.log("bidder address: ", bidder);
        console.log("bidToken address: ", bidToken);
        console.log("bidAmount : ", bidAmount);
        console.log("bidUsd : ", bidUsd);
    }

    function testV1AuctionBidSendBackERC20() public auctionPrepare{
        vm.warp(block.timestamp + 30);

        uint256 auctionId = 1;
        address eT1Token = user1.eT1Token;

        // user1出价
        uint256 user1BalanceErc20 = IERC20(eT1Token).balanceOf(user1.userAddress);
        console.log("user1BalanceErc20 before bid: ", user1BalanceErc20);
        
        uint256 user1_bidprice_ERC20 = 100 * 10 **18;
        vm.prank(user1.userAddress);
        IERC20(eT1Token).approve(address(v1Auction), user1_bidprice_ERC20);

        vm.prank(user1.userAddress);
        v1Auction.placeBidERC20(auctionId, eT1Token, user1_bidprice_ERC20);
        console.log("user1BalanceErc20 after bid: ", IERC20(eT1Token).balanceOf(user1.userAddress));

        // user2出价
        uint256 user2_bidprice_ERC20 = 200 * 10 **18;
        vm.prank(user2.userAddress);
        IERC20(eT1Token).approve(address(v1Auction), user2_bidprice_ERC20);

        vm.prank(user2.userAddress);
        v1Auction.placeBidERC20(auctionId, eT1Token, user2_bidprice_ERC20);

        (address bidder,address bidToken, uint256 bidAmount,uint256 bidUsd) = v1Auction.getHighestBidInfo(auctionId);
        console.log("bidder address: ", bidder);
        console.log("bidToken address: ", bidToken);
        console.log("bidAmount : ", bidAmount);
        console.log("bidUsd : ", bidUsd);

        console.log("after user2 Bid, send back user1 ERC20, user1 balance : ", IERC20(eT1Token).balanceOf(user1.userAddress));
        assertEq(
            user1BalanceErc20, 
            IERC20(eT1Token).balanceOf(user1.userAddress), 
            "nftAuctionV1 has not send back user1.");
    }

    function testV1AuctionUser1BidEthUser2BidERC20() public auctionPrepare{
        vm.warp(block.timestamp + 30);
        uint256 auctionId = 1;
        address eT1Token = user1.eT1Token;

        // user1出价
        uint256 user1Balance = payable(user1.userAddress).balance;
        console.log("user1Balance before bid: ", user1Balance);

        vm.prank(user1.userAddress);
        v1Auction.placeBidETH{value: 0.1 ether}(auctionId);
        console.log("user1Balance after bid: ", payable(user1.userAddress).balance);

        // user2出价
        uint256 user2_bidprice_ERC20 = 2001 * 10 **18;
        vm.prank(user2.userAddress);
        IERC20(eT1Token).approve(address(v1Auction), user2_bidprice_ERC20);

        vm.prank(user2.userAddress);
        v1Auction.placeBidERC20(auctionId, eT1Token, user2_bidprice_ERC20);

        console.log("user1Balance after user2 bid: ", payable(user1.userAddress).balance);
        assertEq(user1Balance, payable(user1.userAddress).balance, "ETH bid sendback right amount.");

        (address bidder,address bidToken, uint256 bidAmount,uint256 bidUsd) = v1Auction.getHighestBidInfo(auctionId);
        console.log("bidder address: ", bidder);
        console.log("bidToken address: ", bidToken);
        console.log("bidAmount : ", bidAmount);
        console.log("bidUsd : ", bidUsd);
    }

    function testV1AuctionUser1BidERC20User2BidEth() public auctionPrepare{
        vm.warp(block.timestamp + 30);
        uint256 auctionId = 1;
        address eT1Token = user1.eT1Token;
        
        uint256 user1_bidprice_ERC20 = 100 * 10 **18; // 10 USD
        uint256 user1_bidprice_ERC20_USD = v1Auction.getPriceUSD(eT1Token, user1_bidprice_ERC20);
        console.log("user1_bidprice_ERC20_USD : ", user1_bidprice_ERC20_USD);

        uint256 user2_bidprice_ETH = 1 * 10 **17; // 20 USD
        uint256 user2_bidprice_ETH_USD = v1Auction.getPriceUSD(address(0), user2_bidprice_ETH);
        console.log("user2_bidprice_ETH_USD : ", user2_bidprice_ETH_USD);

        // user1出价
        uint256 user1BalanceErc20 = IERC20(eT1Token).balanceOf(user1.userAddress);
        console.log("user1BalanceErc20 before bid: ", user1BalanceErc20);

        vm.prank(user1.userAddress);
        IERC20(eT1Token).approve(address(v1Auction), user1_bidprice_ERC20);

        vm.prank(user1.userAddress);
        v1Auction.placeBidERC20(auctionId, eT1Token, user1_bidprice_ERC20);
        console.log("user1BalanceErc20 after bid: ", IERC20(eT1Token).balanceOf(user1.userAddress));

        // user2出价
        vm.prank(user2.userAddress);
        v1Auction.placeBidETH{value: user2_bidprice_ETH * 1 wei}(auctionId);

        console.log("user1Balance after user2 bid ETH: ", IERC20(eT1Token).balanceOf(user1.userAddress));
        assertEq(user1BalanceErc20, IERC20(eT1Token).balanceOf(user1.userAddress), "ETH bid sendback right amount.");

        (address bidder,address bidToken, uint256 bidAmount,uint256 bidUsd) = v1Auction.getHighestBidInfo(auctionId);
        console.log("bidder address: ", bidder);
        console.log("bidToken address: ", bidToken);
        console.log("bidAmount : ", bidAmount);
        console.log("bidUsd : ", bidUsd);
    }

    function testV1AuctionUser1BidOutTime() public auctionPrepare {
        vm.warp(block.timestamp + 130);
        // user1出价
        uint256 auctionId = 1;
        uint256 user1Balance = payable(user1.userAddress).balance;
        console.log("user1Balance before bid: ", user1Balance);

        vm.prank(user1.userAddress);
        vm.expectRevert(bytes("nftAuctionV1: auction ended."));
        v1Auction.placeBidETH{value: 0.1 ether}(auctionId);

        user1Balance = payable(user1.userAddress).balance;
        console.log("user1Balance after bid: ", user1Balance);

        (address bidder,address bidToken, uint256 bidAmount,uint256 bidUsd) = v1Auction.getHighestBidInfo(auctionId);
        console.log("bidder address: ", bidder);
        console.log("bidToken address: ", bidToken);
        console.log("bidAmount : ", bidAmount);
        console.log("bidUsd : ", bidUsd);
    }

    function testV1AuctionUser1BidETHEndAuction() public auctionPrepare {
        vm.warp(block.timestamp + 30);
        // user1出价
        uint256 auctionId = 1;
        uint256 user1Balance = payable(user1.userAddress).balance;
        console.log("user1Balance before bid: ", user1Balance);

        uint256 user2_bidprice_ETH = 1 * 10 **17; // 20 USD
        uint256 afterAuction_accountFee = (getFeePercent() * user2_bidprice_ETH) / 100;
        uint256 afterAuction_sellerFee = user2_bidprice_ETH - afterAuction_accountFee;

        vm.prank(user1.userAddress);
        v1Auction.placeBidETH{value: user2_bidprice_ETH * 1 wei}(auctionId);

        user1Balance = payable(user1.userAddress).balance;
        console.log("user1Balance after bid: ", user1Balance);

        (address bidder,address bidToken, uint256 bidAmount,uint256 bidUsd) = v1Auction.getHighestBidInfo(auctionId);
        console.log("bidder address: ", bidder);
        console.log("bidToken address: ", bidToken);
        console.log("bidAmount : ", bidAmount);
        console.log("bidUsd : ", bidUsd);

        vm.warp(block.timestamp + 130);
        /* user1 结算，结算后：
        1. NFT owner 转为user1;
        2. user1 出价20 USD, 手续费3%, 平台收取0.6 USD, user3收取19.4 USD;
        3. 折算 ETH即，平台收取：3e15, user3收取：97e15;
        4. user1 获得NFT;
        */
        uint256 auctionAccountBalance = payable(getFeeAccount()).balance;
        uint256 user3Balance = payable(user3.userAddress).balance;

        vm.prank(user3.userAddress);
        v1Auction.endAuction(auctionId);

        address nftOwner = IERC721(user3.nftAddress).ownerOf(user3.nftTokenId);
        assertEq(nftOwner, user1.userAddress, "IERC721 nft has transfer Auction contract.");
        
        uint256 auctionAccountFee = payable(getFeeAccount()).balance - auctionAccountBalance;
        assertEq(auctionAccountFee, afterAuction_accountFee, "auctionAccountFee is wrong.");

        uint256 user3AuctionFee = payable(user3.userAddress).balance - user3Balance;
        assertEq(user3AuctionFee, afterAuction_sellerFee, "user3AuctionFee is wrong.");
    }

    function testV1AuctionUser1BidErc20EndAuction() public auctionPrepare {
        vm.warp(block.timestamp + 30);
        // user1出价
        uint256 auctionId = 1;
        address eT1Token = user1.eT1Token;

        uint256 user1_bidprice_ERC20 = 100 * 10 **18; // 10 USD
        uint256 afterAuction_accountFee = (getFeePercent() * user1_bidprice_ERC20) / 100;
        uint256 afterAuction_sellerFee = user1_bidprice_ERC20 - afterAuction_accountFee;

        uint256 user1_bidprice_ERC20_USD = v1Auction.getPriceUSD(eT1Token, user1_bidprice_ERC20);
        console.log("user1_bidprice_ERC20_USD : ", user1_bidprice_ERC20_USD);

        uint256 user1Balance = IERC20(eT1Token).balanceOf(user1.userAddress);
        console.log("user1Balance before bid: ", user1Balance);

        vm.prank(user1.userAddress);
        IERC20(user1.eT1Token).approve(address(v1Auction), user1_bidprice_ERC20);

        vm.prank(user1.userAddress);
        v1Auction.placeBidERC20(auctionId, eT1Token, user1_bidprice_ERC20);

        user1Balance = IERC20(eT1Token).balanceOf(user1.userAddress);
        console.log("user1Balance after bid: ", user1Balance);

        (address bidder,address bidToken, uint256 bidAmount,uint256 bidUsd) = v1Auction.getHighestBidInfo(auctionId);
        console.log("bidder address: ", bidder);
        console.log("bidToken address: ", bidToken);
        console.log("bidAmount : ", bidAmount);
        console.log("bidUsd : ", bidUsd);

        vm.warp(block.timestamp + 130);
        /* user1 结算，结算后：
        1. NFT owner 转为user1;
        2. user1 出价100 ERC20 代币，平台收取：3e18, user3收取：97e18;
        3. user1 获得NFT;
        */
        uint256 auctionAccountBalance = IERC20(eT1Token).balanceOf(getFeeAccount());
        uint256 user3Balance = IERC20(eT1Token).balanceOf(user3.userAddress);
        address nftOwner = IERC721(user3.nftAddress).ownerOf(user3.nftTokenId);
        console.log("Auction End before NFT belong of:", nftOwner);

        vm.prank(user3.userAddress);
        v1Auction.endAuction(auctionId);

        nftOwner = IERC721(user3.nftAddress).ownerOf(user3.nftTokenId);
        console.log("Auction End after NFT belong of:", nftOwner);

        assertEq(nftOwner, user1.userAddress, "IERC721 nft has transfer Auction contract.");
        
        uint256 auctionAccountFee = IERC20(eT1Token).balanceOf(getFeeAccount()) - auctionAccountBalance;
        assertEq(auctionAccountFee, afterAuction_accountFee, "auctionAccountFee is wrong.");

        uint256 user3AuctionFee = IERC20(eT1Token).balanceOf(user3.userAddress) - user3Balance;
        assertEq(user3AuctionFee, afterAuction_sellerFee, "user3AuctionFee is wrong.");
    }
}