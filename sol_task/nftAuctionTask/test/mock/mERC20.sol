// SPDX-License-Identifier: UNLICENSED
pragma solidity >0.8.0 <0.9.0;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {ERC20Permit} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Permit.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

contract TestERC20 is ERC20, ERC20Permit,Ownable {

    constructor(string memory name_, string memory symbol_, address _admin) ERC20(name_, symbol_) ERC20Permit(name_) Ownable(_admin) {
    }

    function mint(address to, uint256 amount) public onlyOwner {
        _mint(to, amount);
    }

    function getOwner() external view returns(address){
        return owner();
    }

    function balanceOf(address account) public view override returns (uint256) {
        return super.balanceOf(account);
    }

}