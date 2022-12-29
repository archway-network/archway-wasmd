<!--
order: 0
title: Tracking Overview
parent:
  title: "Wasm"
-->


# Wasm

## Abstract
Wasm provides a fully compatible CosmWasm runtime.

## Archway Spec

### Wasmer Engine
The WasmerEngine defines the WASM contract runtime engine. Archway has made modifications to allow for tracking of gas through an execution thread of contracts.

#### Execution Thread
An execution thread happens when a smart contract calls a secondary smart contract in it's execution. The ability to correctly identify execution threads is paramount to being able to gauge complex contract interactions and how they affect gas consumption on the block.

#### QuerierWithCtx
Querier with Context represents an Archway unique interface that retains context awareness from the begining of the execution thread to the end.

#### PrefixStoreInfo
A custom store structure that prefixes store information for the smart contract.

### Tracking Wasmer Engine
The Archway implementation of the Wasmer Engine. it wraps around the official wasm engine and extends it's functionality using a custom `ContractGasProcessor`

The Tracking Wasmer attempts to only extend the base engine with a couple of functionalities and decorating gas processor utilities around the base engine calls.

Such that tracking happens in the following order: 

1- Retrieve Calculation from [Gas Processor](README.md#ContractGasProcessor)
2- Get Contract Meter 
3- Get Store
4- Interact with Base VM
5- Calculate Gas
5- Add [VM Records] (README.md#VMRecord)



### Gas Tracking
Gas tracking involves the combination of [VM Records](README.md#VMRecord), [Session Records](README.md#SessionRecord) & [Contract Gas Processor](README.md#SessionRecord) 


#### ContractGasProcessor
Defines the contract for which the engine will gauge gas records.

```
type ContractGasProcessor interface {
	IngestGasRecord(ctx sdk.Context, records []ContractGasRecord) error
	CalculateUpdatedGas(ctx sdk.Context, record ContractGasRecord) (GasConsumptionInfo, error)
	GetGasCalculationFn(ctx sdk.Context, contractAddress string) (func(operationId uint64, gasInfo GasConsumptionInfo) GasConsumptionInfo, error)
}
```

Note: It's up to the implementation side to determine how the gas will be processed if at all, as well as what to do with the gas records.

#### VMRecord
VMRecords are used to capture information of execution thread retaining awareness invoker gas and it's invokee, this is placed into an active session, once the active session is over and no longer needed
The VMRecord is created with a contrast within the gas consumed in the base VM and the actual gas consumed after the gas gauge calculation from the gas processor.

```
type VMRecord struct {
 OriginalVMGas // Gas consumed by the BaseVM
 ActualVMGas // Gas consumed by the gas processor calculations
}
```

#### SessionRecord
Session Records are stored in state, they consolidate all VMRecords as well as SDK gas info and what operation caused the Session.
