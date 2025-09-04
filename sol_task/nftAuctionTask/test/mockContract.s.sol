// SPDX-License-Identifier: UNLICENSED
pragma solidity >0.8.0 <0.9.0;

import {Test} from "forge-std/Test.sol";
import {TestERC20} from "../test/mock/mERC20.sol";
import {task02_nft, IERC721} from "../test/mock/mERC721.sol";
import {MockV3Aggregator} from "../test/mock/mV3Aggregator.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";


contract mockContract is Test {
    TestERC20 public eT1;
    task02_nft public nft;
    address public constant DEPLOY_USER = address(0x10);
    address public constant USER1 = address(0x1);
    address public constant USER2 = address(0x2);
    string public constant NFT_TEST_URI = "https://prominent-moccasin-thrush.myfilebase.com/ipfs/QmSh6DUx9dSgtcA6Tn9GNVqCjsxWt9ptz1BHns94Udhx37";

    function setUp() public {
        // Setup code if needed
        vm.startBroadcast();
        
        eT1 = new TestERC20("myTestErc20", "ET1", DEPLOY_USER);
        nft = new task02_nft("myTestNft", "MNFT", DEPLOY_USER);
        
        vm.stopBroadcast();
    }

    function testErc20Owner() public {
        address owner = eT1.getOwner();
        assertEq(owner, DEPLOY_USER, "eT1 owner is not DEPLOY_USER");
    }

    function testNftOwner() public {
        address owner = nft.getOwner();
        assertEq(owner, DEPLOY_USER, "nft owner is not DEPLOY_USER");
    }

    function testErc20Mint() public {
        vm.startPrank(DEPLOY_USER);
        eT1.mint(USER1, 1000 * 10**18);
        uint256 balance = eT1.balanceOf(USER1);
        assertEq(balance, 1000 * 10**18, "USER1 balance is not correct");
        vm.stopPrank();
    }

    function testErc20Transfer() public {
        vm.startPrank(DEPLOY_USER);
        eT1.mint(DEPLOY_USER, 100 * 10**18);
        eT1.approve(DEPLOY_USER, 100 * 10**18);
        vm.stopPrank();

        uint256 balance = eT1.balanceOf(USER1);
        vm.startPrank(DEPLOY_USER);
        eT1.transferFrom(DEPLOY_USER, USER1, 10 * 10**18);
        vm.stopPrank();

        balance = eT1.balanceOf(USER1) - balance;
        assertEq(balance, 10 * 10**18, "USER1 balance is not correct");
    }

    function testIerc20Transfer() public {
        vm.prank(DEPLOY_USER);
        eT1.mint(USER1, 1000 * 10**18);
        
        vm.prank(USER1);
        eT1.approve(USER1, 100 * 10**18);

        vm.prank(USER1);
        IERC20(address(eT1)).transferFrom(USER1, USER2, 10 * 10**18);

        uint256 balance = eT1.balanceOf(USER2);
        assertEq(balance, 10 * 10**18, "USER2 balance is not correct");
    }

    function testNftMit() public {
        vm.prank(DEPLOY_USER);
        uint256 tokenId = nft.mintNFT(DEPLOY_USER, NFT_TEST_URI);
        assertEq(tokenId, 1, "the first Nft tokenId is 1.");

        bool isOwner;
        vm.prank(DEPLOY_USER);
        isOwner = nft.isOwnerTokenId(tokenId);
        assertEq(isOwner, true, "tokenId 1 is own Deploy_User");

        vm.prank(USER1);
        isOwner = nft.isOwnerTokenId(tokenId);
        assertEq(isOwner, false, "tokenId 1 is not own USER1");

        vm.prank(USER2);
        isOwner = nft.isOwnerTokenId(tokenId);
        assertEq(isOwner, false, "tokenId 1 is own USER2");
    }

    function testNftTransferNFT() public {
        vm.startPrank(DEPLOY_USER);

        uint256 tokenId = nft.mintNFT(DEPLOY_USER, NFT_TEST_URI);
        assertEq(tokenId, 1, "the first Nft tokenId is 1.");

        bool isOwner = nft.isOwnerTokenId(tokenId);
        assertEq(isOwner, true, "tokenId 1 is own Deploy_User");

        nft.approve(DEPLOY_USER, tokenId);
        nft.transferNFT(DEPLOY_USER, USER1, tokenId);
        
        vm.stopPrank();

        vm.prank(USER1);
        isOwner = nft.isOwnerTokenId(tokenId);
        assertEq(isOwner, true, "tokenId 1 is own USER1");
    }

    function testNftTransferFrom() public {
        vm.prank(DEPLOY_USER);
        uint256 tokenId = nft.mintNFT(DEPLOY_USER, NFT_TEST_URI);
        assertEq(tokenId, 1, "the first Nft tokenId is 1.");

        vm.prank(DEPLOY_USER);
        nft.approve(USER1, tokenId);

        vm.prank(USER1);
        nft.transferFrom(DEPLOY_USER, USER1, tokenId);
        
        vm.prank(USER1);
        bool isOwner = nft.isOwnerTokenId(tokenId);
        assertEq(isOwner, true, "tokenId 1 is own USER1");

        vm.prank(USER1);
        nft.approve(USER1, tokenId);

        vm.prank(USER1);
        IERC721(address(nft)).transferFrom(USER1, USER2, tokenId);

        vm.prank(USER2);
        isOwner = nft.isOwnerTokenId(tokenId);
        assertEq(isOwner, true, "tokenId 1 is own USER2");
    }

    function testAggregatorV3() public {
        vm.startBroadcast();
        MockV3Aggregator pricefee = new MockV3Aggregator(8, 1e8);
        vm.stopBroadcast();
        
        (,int256 answer,,,) = pricefee.latestRoundData();

        assertEq(answer, int256(1e8), "get right price fee.");
    }
}