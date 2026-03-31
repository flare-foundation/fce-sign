package app

// To regenerate contract bindings after modifying the Solidity source:
//
//   1. Compile the contract:
//        cd ../../contract && forge build --root . --contracts . --out out
//
//   2. Extract ABI and BIN:
//        jq -r '.abi' contract/out/InstructionSender.sol/InstructionSender.json > go/tools/app/contract/InstructionSender.abi
//        jq -r '.bytecode.object' contract/out/InstructionSender.sol/InstructionSender.json > go/tools/app/contract/InstructionSender.bin
//
//   3. Run go generate from go/tools/:
//        go generate ./...

//go:generate abigen --abi contract/InstructionSender.abi --bin contract/InstructionSender.bin --pkg contract --type InstructionSender --out contract/autogen.go
