package params

import (
	"context"
	"math/big"
)

type ContextWithConfig interface {
	context.Context
	WithEIPsBlockFlags(blockNum *big.Int) ContextWithForkFlags
	GetChainID() *big.Int
}

type contextWithConfig struct {
	context.Context
}

func NewContext(c *ChainConfig) ContextWithConfig {
	if c == nil {
		panic("NewContext(c *ChainConfig)")
	}
	return c.WithConfig(context.Background())
}

func (c *ChainConfig) WithConfig(ctx context.Context) ContextWithConfig {
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

	return contextWithConfig{ctx}
}

func (ctx contextWithConfig) WithEIPsBlockFlags(blockNum *big.Int) ContextWithForkFlags {
	for _, eipFlag := range eipFlags {
		ctx.Context = ctx.withUpdateEIPFlag(blockNum, eipFlag.eipFlag, eipFlag.isEIPFlag)
	}

	ctx.Context = context.WithValue(ctx, BlockNumber, blockNum)
	return contextWithForkFlags{ctx}
}

func (ctx contextWithConfig) withUpdateEIPFlag(blockNum *big.Int, eipFlag, isEIPFuncFlag configKey) contextWithConfig {
	isEIPFunc := getForkFunc(ctx, isEIPFuncFlag)
	if isEIPFunc == nil {
		// FIXME: we need to do something to not panic
		panic("not nil func expected")
	}

	ctx.Context =  context.WithValue(ctx, eipFlag, isEIPFunc(blockNum))

	return ctx
}

func (ctx contextWithConfig) GetChainID() *big.Int {
	return getBigInt(ctx, ChainID)
}

func getChainID(c *big.Int) *big.Int {
	chainID := big.NewInt(0)
	if c != nil {
		chainID.Set(c)
	}

	return chainID
}
