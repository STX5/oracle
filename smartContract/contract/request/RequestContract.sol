// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.5.17;

contract RequestContract {

    event RequestEvent(address from, string value);

    constructor() {}

    function request(string memory url, string memory pattern) public  {

        // bytes memory tempUrl = bytes(url);
        // bytes memory tempPattern = bytes(pattern);
        // uint lenUrl = tempUrl.length;
        // uint lenPattern = tempPattern.length;
        // uint lenValue = lenUrl + lenPattern;
        // bytes memory tempValue = new bytes(lenValue);

        // for (uint i = 0; i < lenValue; i++) {

        //     if (i < lenUrl) {
        //         tempValue[i] = tempUrl[i];
        //     } else if (i >= lenUrl) {
        //         tempValue[i] = tempPattern[i - lenUrl];
        //     }
        // }
        string memory pre = '\'{"url":"';
        string memory mid = '","pattern":"';
        string memory post = '"}\'';
        string memory value = string(abi.encodePacked(pre, url, mid, pattern, post));

        emit RequestEvent(msg.sender, value);
    }
}