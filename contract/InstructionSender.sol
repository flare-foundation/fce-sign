// SPDX-License-Identifier: MIT
pragma solidity >=0.7.6 <0.9;

import { ITeeExtensionRegistry } from "./interface/ITeeExtensionRegistry.sol";
import { ITeeMachineRegistry } from "./interface/ITeeMachineRegistry.sol";

contract InstructionSender {
    ITeeExtensionRegistry public immutable teeExtensionRegistry;
    ITeeMachineRegistry public immutable teeMachineRegistry;

    uint256 public _extensionId;

    constructor(
        address _teeExtensionRegistry,
        address _teeMachineRegistry
    ) {
        teeExtensionRegistry = ITeeExtensionRegistry(_teeExtensionRegistry);
        teeMachineRegistry = ITeeMachineRegistry(_teeMachineRegistry);
    }

    /// @notice Discover and store this contract's extension ID.
    function setExtensionId() external {
        require(_extensionId == 0, "extension ID already set");

        uint256 count = teeExtensionRegistry.extensionsCounter();
        for (uint256 i = 1; i <= count; i++) {
            if (
                teeExtensionRegistry.getTeeExtensionInstructionsSender(i) ==
                address(this)
            ) {
                _extensionId = i;
                return;
            }
        }
        revert("extension ID not found");
    }

    /// @notice Update the stored private key by sending an encrypted key to the TEE.
    function updateKey(bytes calldata _encryptedKey) external payable returns (bytes32) {
        require(_extensionId != 0, "extension ID not set");

        address[] memory teeIds = teeMachineRegistry.getRandomTeeIds(_extensionId, 1);

        ITeeExtensionRegistry.TeeInstructionParams memory params;
        params.opType = bytes32("KEY");
        params.opCommand = bytes32("UPDATE");
        params.message = _encryptedKey;

        return teeExtensionRegistry.sendInstructions{value: msg.value}(teeIds, params);
    }

    /// @notice Request the TEE to sign a message with the stored private key.
    function sign(bytes calldata _message) external payable returns (bytes32) {
        require(_extensionId != 0, "extension ID not set");

        address[] memory teeIds = teeMachineRegistry.getRandomTeeIds(_extensionId, 1);

        ITeeExtensionRegistry.TeeInstructionParams memory params;
        params.opType = bytes32("KEY");
        params.opCommand = bytes32("SIGN");
        params.message = _message;

        return teeExtensionRegistry.sendInstructions{value: msg.value}(teeIds, params);
    }
}
