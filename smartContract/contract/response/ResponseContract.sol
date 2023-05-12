// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.5.17;

contract ResponseContract {

    mapping (address => string) mapValue;

    function setValue(address from, string memory value) public {

        mapValue[from] = value;
    }

    function getValue() public view returns (string memory) {

        return mapValue[msg.sender];
    }
}