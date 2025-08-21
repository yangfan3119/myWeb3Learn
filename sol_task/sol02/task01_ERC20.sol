// SPDX-License-Identifier: UNLICENSED
pragma solidity >0.8.0 <0.9.0;

/*
作业 1：ERC20 代币
任务：参考 openzeppelin-contracts/contracts/token/ERC20/IERC20.sol实现一个简单的 ERC20 代币合约。要求：
合约包含以下标准 ERC20 功能：
1. balanceOf：查询账户余额。
2. transfer：转账。
3. approve 和 transferFrom：授权和代扣转账。
4. 使用 event 记录转账和授权操作。
5. 提供 mint 函数，允许合约所有者增发代币。
提示：
- 使用 mapping 存储账户余额和授权信息。
- 使用 event 定义 Transfer 和 Approval 事件。
- 部署到sepolia 测试网，导入到自己的钱包
*/
contract task01_ERC20 {
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);

    mapping(address account => uint256) private _balances;
    mapping(address owner => mapping(address spender => uint256)) private _allowances;

    string private _name;

    string private _symbol;

    address private _owner;

    constructor(string memory name_, string memory symbol_){
        _name = name_;
        _symbol = symbol_;
        _owner = msg.sender;
    }

    modifier mERC20onlyOwner(){
        require(msg.sender == _owner, "myERC20: only owner can call this function.");
        _;
    }

    function changeOwner(address newOwner) public mERC20onlyOwner{
        require(newOwner != address(0), "myERC20: new owner is the zero address.");
        _owner = newOwner;
    }

    function balanceOf() public view returns(uint256){
        return _balances[msg.sender];
    }

    function transfer(address to, uint256 amount) public returns(bool){
        require(msg.sender != address(0), "ERC20: transfer from the zero address");
        require(to != address(0), "ERC20: transfer to the zero address");

        require(_balances[msg.sender] >= amount, "Not enough banlance.");
        _balances[msg.sender] = _balances[msg.sender] - amount;
        _balances[to] = _balances[to] + amount;

        emit Transfer(msg.sender, to, amount);
        return true;
    }

    function approve(address spender, uint256 amount)public returns(bool){
        require(msg.sender != address(0), "ERC20: approve from the zero address");
        require(spender != address(0), "ERC20: approve to the zero address");

        _allowances[msg.sender][spender] = amount;
        emit Approval(msg.sender, spender, amount);
        return true;
    }

    function transferFrom(address from, address to, uint256 amount) public returns(bool){
        require(from != address(0), "myERC20: transfer from the zero address.");
        require(to != address(0), "myERC20: transfer to the zero address.");

        require(_allowances[from][to] >= amount, "myERC20: transfer amount exceeds allowance.");

        _allowances[from][to] = _allowances[from][to] - amount;
        _balances[from] = _balances[from] - amount;
        _balances[to] = _balances[to] + amount;

        emit Transfer(from, to, amount);
        return true;
    }

    function mint(address to, uint256 amount) public mERC20onlyOwner returns(bool){
        require(to != address(0), "myERC20: mint to the zero address.");

        _balances[to] = _balances[to] + amount;
        emit Transfer(address(0), to, amount);
        return true;
    }
}
