// SPDX-License-Identifier: UNLICENSED
pragma solidity >0.8.0 <0.9.0;

/*
✅ 创建一个名为Voting的合约，包含以下功能：
一个mapping来存储候选人的得票数
一个vote函数，允许用户投票给某个候选人
一个getVotes函数，返回某个候选人的得票数
一个resetVotes函数，重置所有候选人的得票数
*/

contract Voting {
    address private immutable owner;  // 创建合约所有者
    constructor() {
        owner = msg.sender;
    }

    modifier onlyOwner(){
        require(msg.sender==owner, "This operation is restricted to the contract owner only.");
        _;
    }

    mapping (address => address) OnesVote;  // 存储投票人地址和候选人地址
    address[] votedOnes;    // 存储投票人地址

    struct Candidate {
        address addr;
        uint256 voteCount;
        uint8 age;
        string name;
        string details;
    }

    mapping (address => Candidate) Candidates; // 存储候选人和投票数
    address[] registeredCandidates;

    // 候选人信息导入
    function setCandidates(address cAddr, uint8 cAge, string calldata cName, string calldata cDetails) external onlyOwner {
        require(cAge > 18, "Candidates must be at least 18 years old.");
        require(bytes(cName).length > 4, "Name length must be greater than 4.");

        if (bytes(Candidates[cAddr].name).length > 4) {        // 已存在的候选人，更新资料
            Candidate storage c = Candidates[cAddr];
            c.age = cAge;
            c.addr = cAddr;
            c.name = cName;
            c.details = cDetails;
        } else {
            require(registeredCandidates.length < 100, "Candidates must not exceed 100 in number.");
            Candidates[cAddr] = Candidate({addr: cAddr, voteCount: 0, age: cAge, name: cName, details: cDetails});
            registeredCandidates.push(cAddr);
        }
    }
    // 获取所有候选人信息
    function getAllCandidates() external  view returns(Candidate[] memory){
        uint8 cLen = uint8(registeredCandidates.length);
        Candidate[] memory cs = new Candidate[](cLen);
        
        for(uint8 i = 0; i < cLen; i++){
            address iAddr = registeredCandidates[i];
            Candidate memory c = Candidates[iAddr];
            cs[i] = Candidate({addr: c.addr, voteCount: c.voteCount, age: c.age, name: c.name, details: c.details});
        }
        return cs;
    }

    // 投票,一个vote函数，允许用户投票给某个候选人
    function vote(address cAddr) external {
        require(OnesVote[msg.sender] == address(0), "Duplicate voting is not allowed.");
        require(bytes(Candidates[cAddr].name).length > 0, "invalid candidate.");

        OnesVote[msg.sender] = cAddr;
        votedOnes.push(msg.sender);
        Candidate storage c = Candidates[cAddr];
        c.voteCount = c.voteCount+1;
    } 

    // 一个getVotes函数，返回某个候选人的得票数
    function getVotes(address cAddr) external view returns (uint256){
        require(bytes(Candidates[cAddr].name).length > 0, "invalid candidate.");
        return Candidates[cAddr].voteCount;
    }

    // 一个resetVotes函数，重置所有候选人的得票数
    function resetVotes() external onlyOwner {
        for(uint i = votedOnes.length-1; i >= 0; i--){
            delete OnesVote[votedOnes[i]];
            votedOnes.pop();
        }

        uint8 cLen = uint8(registeredCandidates.length);
        for(uint8 i = 0; i < cLen; i++){
            address iAddr = registeredCandidates[i];
            Candidate storage c = Candidates[iAddr];
            c.voteCount = 0;
        }
    }
}