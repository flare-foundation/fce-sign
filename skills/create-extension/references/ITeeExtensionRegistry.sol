// SPDX-License-Identifier: MIT
// Bundled reference copy for tee-builder skill.
// Source: flare-smart-contracts-v2/contracts/userInterfaces/tee/ITeeExtensionRegistry.sol
pragma solidity >=0.7.6 <0.9;

// Minimal stub — full definition in ITeeMachineRegistry.sol
interface ITeeMachineRegistry {
    struct TeeMachine {
        address teeId;
        address teeProxyId;
        string url;
    }
}

/**
 * TeeExtensionRegistry interface.
 */
interface ITeeExtensionRegistry {

    event TeeInstructionsSent(
        uint256 indexed extensionId,
        bytes32 indexed instructionId,
        uint32 indexed rewardEpochId,
        ITeeMachineRegistry.TeeMachine[] teeMachines,
        bytes32 opType,
        bytes32 opCommand,
        bytes message,
        address[] cosigners,
        uint64 cosignersThreshold,
        uint256 fee
    );

    error NoTeeMachinesSpecified();
    error OperationTypeEmpty();
    error OperationCommandEmpty();
    error MessageEmpty();
    error ExtensionIdMismatch();
    error OnlyInstructionsSender();
    error FeeTooLow();
    error TeeMachineNotAvailable();

    /**
     * Send instructions to the TEE machines. Instruction ID will be generated internally and returned.
     * Emits a TeeInstructionsSent event.
     * @param _teeIds The TEE machine IDs to which the instructions are sent (must all belong to the same extension).
     * @param _opType The operation type.
     * @param _opCommand The operation command.
     * @param _message The message.
     * @param _cosigners The cosigners.
     * @param _cosignersThreshold The cosigners threshold.
     * @return _instructionId The generated instruction ID.
     * Can only be called by the TEE machines extension instructions sender.
     */
    function sendInstructions(
        address[] memory _teeIds,
        bytes32 _opType,
        bytes32 _opCommand,
        bytes memory _message,
        address[] memory _cosigners,
        uint64 _cosignersThreshold
    )
        external payable
        returns (bytes32 _instructionId);

    /**
     * Get number of registered TEE extensions.
     * @return The number of registered TEE extensions.
     */
    function extensionsCounter()
        external view
        returns (uint256);

    /**
     * Get the TEE extension instructions sender address.
     * @param _extensionId The id of the extension.
     * @return The TEE extension instructions sender address.
     */
    function getTeeExtensionInstructionsSender(uint256 _extensionId)
        external view
        returns (address);
}
