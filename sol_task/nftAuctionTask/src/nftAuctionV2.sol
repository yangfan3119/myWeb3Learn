// SPDX-License-Identifier: UNLICENSED
pragma solidity >0.8.0 <0.9.0;

import {nftAuctionV1} from "./nftAuctionV1.sol";

contract nftAuctionV2 is nftAuctionV1 {

    // 删除自定义的proxiableUUID函数，使用UUPSUpgradeable的默认实现

    function getVersion() external pure override returns (string memory) {
        return "v2.0.0";
    }

    // 必须重写 _authorizeUpgrade（UUPS 强制要求）
    function _authorizeUpgrade(address newImplementation) internal virtual override onlyOwner {
        // 可以添加额外逻辑，比如黑名单检查
        super._authorizeUpgrade(newImplementation);
    }
}