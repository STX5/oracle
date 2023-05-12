// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.8.2 <0.9.0;


contract OracleRequestContract {
    event SentEvent(address sender, string taskId, string value);

    function request(string memory taskId, string memory value) public {
        emit SentEvent(msg.sender, taskId, value);
    }
}