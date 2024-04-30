package types

import (
	storetypes "cosmossdk.io/store/types"
)

var _ storetypes.GasMeter = &ContractSDKGasMeter{}

type ContractSDKGasMeter struct {
	actualGasConsumed          storetypes.Gas
	originalGas                storetypes.Gas
	underlyingGasMeter         storetypes.GasMeter
	contractAddress            string
	contractOperation          uint64
	contractGasCalculationFunc func(operationId uint64, info GasConsumptionInfo) GasConsumptionInfo
}

func NewContractGasMeter(gasLimit uint64, gasCalculationFunc func(uint64, GasConsumptionInfo) GasConsumptionInfo, contractAddress string, contractOperation uint64) ContractSDKGasMeter {
	return ContractSDKGasMeter{
		actualGasConsumed:          0,
		originalGas:                0,
		contractGasCalculationFunc: gasCalculationFunc,
		underlyingGasMeter:         storetypes.NewGasMeter(gasLimit),
		contractAddress:            contractAddress,
		contractOperation:          contractOperation,
	}
}

func (c *ContractSDKGasMeter) GetContractAddress() string {
	return c.contractAddress
}

func (c *ContractSDKGasMeter) GetContractOperation() uint64 {
	return c.contractOperation
}

func (c *ContractSDKGasMeter) GetOriginalGas() storetypes.Gas {
	return c.originalGas
}

func (c *ContractSDKGasMeter) GetActualGas() storetypes.Gas {
	return c.actualGasConsumed
}

func (c *ContractSDKGasMeter) GasConsumed() storetypes.Gas {
	return c.underlyingGasMeter.GasConsumed()
}

func (c *ContractSDKGasMeter) GasConsumedToLimit() storetypes.Gas {
	return c.underlyingGasMeter.GasConsumedToLimit()
}

func (c *ContractSDKGasMeter) Limit() storetypes.Gas {
	return c.underlyingGasMeter.Limit()
}

func (c *ContractSDKGasMeter) ConsumeGas(amount storetypes.Gas, descriptor string) {
	updatedGasInfo := c.contractGasCalculationFunc(c.contractOperation, GasConsumptionInfo{SDKGas: amount})
	c.underlyingGasMeter.ConsumeGas(updatedGasInfo.SDKGas, descriptor)
	c.originalGas += amount
	c.actualGasConsumed += updatedGasInfo.SDKGas
}

func (c *ContractSDKGasMeter) RefundGas(amount storetypes.Gas, descriptor string) {
	updatedGasInfo := c.contractGasCalculationFunc(c.contractOperation, GasConsumptionInfo{SDKGas: amount})
	c.underlyingGasMeter.RefundGas(updatedGasInfo.SDKGas, descriptor)
	c.originalGas -= amount
	c.actualGasConsumed -= updatedGasInfo.SDKGas
}

func (c *ContractSDKGasMeter) IsPastLimit() bool {
	return c.underlyingGasMeter.IsPastLimit()
}

func (c *ContractSDKGasMeter) IsOutOfGas() bool {
	return c.underlyingGasMeter.IsOutOfGas()
}

func (c *ContractSDKGasMeter) String() string {
	return c.underlyingGasMeter.String()
}

func (c *ContractSDKGasMeter) CloneWithNewLimit(gasLimit uint64, description string) *ContractSDKGasMeter {
	newContractGasMeter := NewContractGasMeter(gasLimit, c.contractGasCalculationFunc, c.contractAddress, c.contractOperation)
	return &newContractGasMeter
}

func (c *ContractSDKGasMeter) GasRemaining() uint64 {
	return c.underlyingGasMeter.GasRemaining()
}
