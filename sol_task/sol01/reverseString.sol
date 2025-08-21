// SPDX-License-Identifier: UNLICENSED
pragma solidity >0.8.0 <0.9.0;

/*
✅ 反转字符串 (Reverse String)
题目描述：反转一个字符串。输入 "abcde"，输出 "edcba"
*/

contract reverseString{

    function t1(string memory str) public pure returns(bytes memory, uint){
        bytes memory strBys = bytes(str);
        return (strBys, strBys.length);
    }

    function resverseStr(string memory str) public pure returns(string memory){
        bytes memory strBys = bytes(str);
        uint slen = strBys.length;
        bytes memory reverBys = new bytes(slen);

        for(uint i=0;i<slen;i++){
            uint j = slen - 1 - i;
            reverBys[j] = strBys[i];
        }
        return string(reverBys);
    }
}