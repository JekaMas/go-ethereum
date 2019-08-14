package params

import (
	"context"
	"fmt"
	"math/big"
)

type configKey int

const (
	IsHomesteadEnabled configKey = iota
	isHomesteadEnabledFunc
	IsEIP150Enabled
	isEIP150EnabledFunc
	IsEIP155Enabled
	isEIP155EnabledFunc
	IsEIP158Enabled
	isEIP158EnabledFunc
	IsByzantiumEnabled
	isByzantiumEnabledFunc
	IsConstantinopleEnabled
	isConstantinopleEnabledFunc
	IsPetersburgEnabled
	isPetersburgEnabledFunc
	IsEWASM
	isEWASMFunc
	BlockNumber
	ChainID
)


type flag struct {
	eipFlag  configKey
	isEIPFlag configKey
}

var eipFlags = []flag{
	{IsHomesteadEnabled, isHomesteadEnabledFunc},
	{IsEIP150Enabled, isEIP150EnabledFunc},
	{IsEIP155Enabled, isEIP155EnabledFunc},
	{IsEIP158Enabled, isEIP158EnabledFunc},
	{IsByzantiumEnabled, isByzantiumEnabledFunc},
	{IsConstantinopleEnabled, isConstantinopleEnabledFunc},
	{IsPetersburgEnabled, isPetersburgEnabledFunc},
	{IsEWASM, isEWASMFunc},
}

type ContextWithForkFlags interface {
	ContextWithConfig
	GetForkFlag(name configKey) bool
	GetBlockNumber() *big.Int
}

type contextWithForkFlags struct {
	ContextWithConfig
}

func NewContextWithBlock(c *ChainConfig, blockNum *big.Int) ContextWithForkFlags {
	return c.WithEIPsFlags(context.Background(), blockNum)
}

func New(ctx context.Context, c *ChainConfig, blockNum *big.Int) ContextWithForkFlags {
	return c.WithEIPsFlags(ctx, blockNum)
}

func (c *ChainConfig) WithEIPsFlags(ctx context.Context, blockNum *big.Int) ContextWithForkFlags {
	checkers := []func(num *big.Int) bool{
		c.IsHomestead,
		c.IsEIP150,
		c.IsEIP155,
		c.IsEIP158,
		c.IsByzantium,
		c.IsConstantinople,
		c.IsPetersburg,
		c.IsEWASM,
	}

	for i, checker := range checkers {
		ctx = context.WithValue(ctx, eipFlags[i].eipFlag, checker(blockNum))
		ctx = context.WithValue(ctx, eipFlags[i].isEIPFlag, checker)
	}

	ctx = context.WithValue(ctx, BlockNumber, blockNum)
	ctx = context.WithValue(ctx, ChainID, getChainID(c.ChainID))

	return contextWithForkFlags{contextWithConfig{ctx}}
}

func (ctx contextWithForkFlags) GetForkFlag(name configKey) bool {
	b := ctx.Value(name)
	if b == nil {
		panic(fmt.Sprint("flag1", name))
		return false
	}
	if valB, ok := b.(bool); ok {
		return valB
	}
	panic("flag2")
	return false
}

func (ctx contextWithForkFlags) GetBlockNumber() *big.Int {
	return getBigInt(ctx, BlockNumber)
}

func getBigInt(ctx context.Context, key configKey) *big.Int {
	b := ctx.Value(key)
	if b == nil {
		panic("getBigInt1")
		return nil
	}
	if valB, ok := b.(*big.Int); ok {
		return valB
	}
	panic("getBigInt2")
	return nil
}

func getForkFunc(ctx context.Context, name configKey) func(num *big.Int) bool {
	b := ctx.Value(name)
	if b == nil {
		panic("getForkFunc1")
		return nil
	}
	if valB, ok := b.(func(num *big.Int) bool); ok {
		return valB
	}
	panic("getForkFunc1")
	return nil
}