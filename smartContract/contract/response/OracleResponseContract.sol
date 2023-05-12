pragma solidity >=0.8.2 <0.9.0;

/**
 * @title 预言机响应合约
 * @dev 预言机响应合约，接收来自worker的数据
 */
contract OracleResponseContract {
    mapping(string => string) public taskResultMap;

    // taskId是当前任务Id
    // value是当前任务的结果
    function write(string memory taskId, string memory taskResult) public {
        taskResultMap[taskId] = taskResult;
    }

    // taskId是当前任务的id
    // 读取任务的结果
    function read(string memory taskId) public returns(string) {
        return taskResultMap[taskId];
    }
}
