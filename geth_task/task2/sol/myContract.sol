// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

contract myCounter {
    event IncreNumber(uint256 number);

    uint256 public number;
    constructor(uint256 _num){
        number = _num;
    }

    function setNumber(uint256 newNumber) public {
        number = newNumber;
    }

    function increment() public {
        number++;
        emit IncreNumber(number);
    }
}