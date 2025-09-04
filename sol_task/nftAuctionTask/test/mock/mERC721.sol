// SPDX-License-Identifier: UNLICENSED
pragma solidity >0.8.0 <0.9.0;

import {ERC721, IERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import {ERC721URIStorage} from "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

contract task02_nft is ERC721URIStorage, Ownable {
    uint256 private _nextTokenId;

    constructor(string memory _name, string memory _symble, address _admin) ERC721(_name, _symble) Ownable(_admin) {
        _nextTokenId++;
    }

    function supportsInterface(bytes4 interfaceId) public view virtual override(ERC721URIStorage) returns (bool) {
        return super.supportsInterface(interfaceId);
    }

    function mintNFT(address recipient, string memory uri)external onlyOwner returns(uint256){
        require(recipient != address(0), "recipient address is not 0.");
        require(bytes(uri).length > 0, "tokenURI is not null.");

        // 获取tokenId，使用ERC721safeMint铸造绑定address和tokenId
        uint256 tokenId = _nextTokenId;
        _nextTokenId++;
        _safeMint(recipient, tokenId);
        _setTokenURI(tokenId, uri);

        return tokenId;
    }

    function transferNFT(address from, address to, uint256 tokenId) external onlyOwner{
        // 调用ERC721的转账函数，会自动检查权限
        transferFrom(from, to, tokenId);
    }

    function isOwnerTokenId(uint256 tokenId) external view returns(bool){
        address owner = ownerOf(tokenId);
        return owner == msg.sender;
    }

    function ownerOf(uint256 tokenId) public view override(IERC721, ERC721) returns (address) {
        return super.ownerOf(tokenId);
    }

    function tokenURI(uint256 tokenId) public view override(ERC721URIStorage) returns (string memory) {
        return super.tokenURI(tokenId);
    }

    function getOwner() external view returns(address){
        return owner();
    }
}