// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8;

/*
二分查找 (Binary Search)
题目描述：在一个有序数组中查找目标值。
*/

contract BinarySearch {

    uint[] searchArray = [1,2,3,4,6,9,8,10,15,20,24,28,29,30,33,36,37,48,49,57];

    function binarySearch(uint a) public view returns(uint,uint) {
        uint iStart = 0;
        uint iEnd = searchArray.length;

        uint index = 0;
        do {
            index = iStart + (iEnd - iStart)/2;
            if(searchArray[index] > a){
                iEnd = index -1;
            }
            if (searchArray[index] < a){
                iStart = index +1;
            }
        }while (searchArray[index] != a);

        return (index,searchArray[index]);
    }
}