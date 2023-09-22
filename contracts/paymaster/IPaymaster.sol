pragma solidity ^0.5.4;

interface IPaymaster {
    struct Transaction {
        address from;
    }

    struct ExecutionResult {
        bool success;
    }

    function validateAndPayForPaymasterTransaction(
        bytes32 _txHash,
        Transaction calldata _transaction
    ) external payable returns (bytes4 magic, bytes memory context);

    function postTransaction(
        bytes calldata _context,
        Transaction calldata _transaction,
        bytes32 _txHash,
        ExecutionResult _txResult,
        uint256 _maxRefundedGas
    ) external payable;
}
