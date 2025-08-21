// SPDX-License-Identifier: UNLICENSED
pragma solidity >0.8.0 <0.9.0;

/*
合并两个有序数组 (Merge Sorted Array)
题目描述：将两个有序数组合并为一个有序数组。
*/

contract MergeSortedArray {

    function mergeArray(uint[] calldata aArray, uint[] calldata bArray)public pure returns(uint[] memory){
        uint alen = aArray.length;
        uint blen = bArray.length;
        
        uint[] memory res = new uint[](alen+blen);

        uint ai = 0;
        uint bi = 0;
        uint resi = 0;
        while (true) {
            uint a = aArray[ai];
            uint b = bArray[bi];
            if (a > b) {
                res[resi] = b;
                bi++;
            }else{
                res[resi] = a;
                ai++;
            }
            resi++;
            if (ai >= alen) {
                for (;bi < blen;bi++){
                    res[resi] = bArray[bi];
                    resi++;
                }
                break ;
            }
            if (bi >= blen) {
                for (;ai < alen;ai++){
                    res[resi] = aArray[ai];
                    resi++;
                }
                break ;
            }
        }
        return res;
    }
}