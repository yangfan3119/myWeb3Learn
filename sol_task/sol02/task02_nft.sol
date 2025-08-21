// SPDX-License-Identifier: UNLICENSED
pragma solidity >0.8.0 <0.9.0;

// 导入OpenZeppelin的ERC721标准实现和所有权管理库
// ERC721是NFT的基础标准，Ownable用于权限控制（仅所有者可执行特定操作）
// 用于安全铸造NFT（检查接收地址是否支持ERC721）
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/utils/ERC721Holder.sol";

/*
作业2：在测试网上发行一个图文并茂的 NFT
任务目标:
1. 使用 Solidity 编写一个符合 ERC721 标准的 NFT 合约。
2. 将图文数据上传到 IPFS，生成元数据链接。
3. 将合约部署到以太坊测试网（如 Goerli 或 Sepolia）。
4. 铸造 NFT 并在测试网环境中查看。

任务步骤
1. 编写 NFT 合约
- 使用 OpenZeppelin 的 ERC721 库编写一个 NFT 合约。
- 合约应包含以下功能：
- 构造函数：设置 NFT 的名称和符号。
- mintNFT 函数：允许用户铸造 NFT，并关联元数据链接（tokenURI）。
- 在 Remix IDE 中编译合约。
2. 准备图文数据
- 准备一张图片，并将其上传到 IPFS（可以使用 Pinata 或其他工具）。
- 创建一个 JSON 文件，描述 NFT 的属性（如名称、描述、图片链接等）。
- 将 JSON 文件上传到 IPFS，获取元数据链接。
- JSON文件参考 https://docs.opensea.io/docs/metadata-standards
3. 部署合约到测试网
- 在 Remix IDE 中连接 MetaMask，并确保 MetaMask 连接到 Goerli 或 Sepolia 测试网。
- 部署 NFT 合约到测试网，并记录合约地址。
4. 铸造 NFT
- 使用 mintNFT 函数铸造 NFT：
- 在 recipient 字段中输入你的钱包地址。
- 在 tokenURI 字段中输入元数据的 IPFS 链接。
- 在 MetaMask 中确认交易。
5. 查看 NFT
- 打开 OpenSea 测试网 或 Etherscan 测试网。
- 连接你的钱包，查看你铸造的 NFT。

0xa513E6E4b8f2a923D98304ec87F64353C4D5C853
*/
contract task02_nft is ERC721 {
    // 存储每个NFT的元数据URI（tokenId => 元数据地址）
    mapping (uint256 => string) private  _tokenURIs;

    uint256 private _nextTokenId;

    address private _owner;

    constructor(string memory _name, string memory _symble) ERC721(_name, _symble){
        _owner = msg.sender;
        _nextTokenId++;
    }

    modifier mOnlyOwner(){
        require(msg.sender == _owner, "myERC20: only owner can call this function.");
        _;
    }

    function mintNFT(address recipient, string memory uri)external mOnlyOwner returns(uint256){
        require(recipient != address(0), "recipient address is not 0.");
        require(bytes(uri).length > 0, "tokenURI is not null.");

        // 获取tokenId，使用ERC721safeMint铸造绑定address和tokenId
        uint256 tokenId = _nextTokenId;
        _nextTokenId++;
        _safeMint(recipient, tokenId);

        // 关联address和URI地址
        _tokenURIs[tokenId] = uri;

        return tokenId;
    }

    function transferNFT(address from, address to, uint256 tokenId) external {
        // 调用ERC721的转账函数，会自动检查权限
        transferFrom(from, to, tokenId);
    }

    function tokenURI(uint256 tokenId) public view override returns (string memory) {
        require(tokenExists(tokenId), "NFT tokenId is unexist.");
        return _tokenURIs[tokenId];
    }

    function tokenExists(uint256 tokenId)private view returns(bool){
        return (bytes(_tokenURIs[tokenId]).length != 0);
    }
}