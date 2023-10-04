pragma solidity ^0.8.0;

contract IPaymaster {
    uint256[] public a;

    struct Transaction {
        address from;
    }

    struct ExecutionResult {
        bool success;
    }

    constructor(){}
    receive() payable external {}

    function validateAndPayForPaymasterTransaction(
        bytes32 _txHash,
        Transaction calldata _transaction
    ) external payable returns (bytes4 magic, bytes memory context) {
        a.push(a.length);
        require(_txHash != "", "empty txHash");
        a.push(a.length);
        magic = bytes4(abi.encodePacked(a.length));
        return (magic, context);
    }

    function postTransaction(
        bytes calldata _context,
        Transaction calldata _transaction,
        bytes32 _txHash,
        ExecutionResult calldata _txResult,
        uint256 _maxRefundedGas
    ) external payable {
        a.push(a.length);
        require(_txHash != "", "empty txHash");
        a.push(a.length);
    }
}
