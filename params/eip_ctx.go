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

func NewContext(c *ChainConfig) context.Context {
	if c == nil {
		panic("NewContext(c *ChainConfig)")
	}
	return c.WithConfig(context.Background())
}

func NewContextWithBlock(c *ChainConfig, blockNum *big.Int) context.Context {
	return c.WithEIPsFlags(context.Background(), blockNum)
}

func (c *ChainConfig) WithConfig(ctx context.Context) context.Context {
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
		ctx = context.WithValue(ctx, eipFlags[i].isEIPFlag, checker)
	}

	ctx = context.WithValue(ctx, ChainID, getChainID(c.ChainID))

	return ctx
}

func (c *ChainConfig) WithEIPsFlags(ctx context.Context, blockNum *big.Int) context.Context {
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

	return ctx
}

func getChainID(c *big.Int) *big.Int {
	chainID := big.NewInt(0)
	if c != nil {
		chainID.Set(c)
	} else {
		//panic("getChainID")
	}

	return chainID
}

func WithEIPsBlockFlags(ctx context.Context, blockNum *big.Int) context.Context {
	for _, eipFlag := range eipFlags {
		ctx = withUpdateEIPFlag(ctx, blockNum, eipFlag.eipFlag, eipFlag.isEIPFlag)
	}

	ctx = context.WithValue(ctx, BlockNumber, blockNum)
	return ctx
}

func withUpdateEIPFlag(ctx context.Context, blockNum *big.Int, eipFlag, isEIPFuncFlag configKey) context.Context {
	isEIPFunc := getForkFunc(ctx, isEIPFuncFlag)
	if isEIPFunc == nil {
		// FIXME: we need to do something to not panic
		panic("not nil func expected")
	}
	return context.WithValue(ctx, eipFlag, isEIPFunc(blockNum))
}

func GetForkFlag(ctx context.Context, name configKey) bool {
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

func GetBlockNumber(ctx context.Context) *big.Int {
	return getBigInt(ctx, BlockNumber)
}

func GetChainID(ctx context.Context) *big.Int {
	return getBigInt(ctx, ChainID)
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