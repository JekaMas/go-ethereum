package params

import (
	"context"
	"math/big"
	"time"
)

type ContextWithConfig interface {
	context.Context
	GetContext() context.Context
	GetChainID() *big.Int
	WithEIPsBlockFlags(blockNum *big.Int) ContextWithForkFlags
}

type contextWithConfig struct {
	context.Context
}

func newContextWithConfig(ctx context.Context) *contextWithConfig {
	return &contextWithConfig{ctx}
}

func NewContext(c *ChainConfig) ContextWithConfig {
	return c.WithConfig(context.Background())
}

func (c *ChainConfig) WithConfig(ctx context.Context) *contextWithConfig {
	checkers := c.getCheckers()

	for i, checker := range checkers {
		ctx = context.WithValue(ctx, eipFlags[i].isEIPFlag, checker)
	}

	ctx = context.WithValue(ctx, ChainID, getChainID(c.ChainID))

	return &contextWithConfig{ctx}
}

func (c *ChainConfig) getCheckers() []func(num *big.Int) bool {
	return []func(num *big.Int) bool{
		c.IsHomestead,
		c.IsEIP150,
		c.IsEIP155,
		c.IsEIP158,
		c.IsByzantium,
		c.IsConstantinople,
		c.IsPetersburg,
		c.IsEWASM,
	}
}

func (ctx contextWithConfig) WithEIPsBlockFlags(blockNum *big.Int) ContextWithForkFlags {
	for _, eipFlag := range eipFlags {
		ctx.Context = ctx.withUpdateEIPFlag(blockNum, eipFlag.eipFlag, eipFlag.isEIPFlag)
	}

	ctx.Context = context.WithValue(ctx, BlockNumber, blockNum)
	return &contextWithForkFlags{ctx}
}

func (ctx contextWithConfig) withUpdateEIPFlag(blockNum *big.Int, eipFlag, isEIPFuncFlag configKey) context.Context {
	isEIPFunc := getForkFunc(ctx, isEIPFuncFlag)
	if isEIPFunc == nil {
		return ctx
	}

	return context.WithValue(ctx, eipFlag, isEIPFunc(blockNum))
}

func (ctx contextWithConfig) GetChainID() *big.Int {
	return getBigInt(ctx, ChainID)
}

func (ctx contextWithConfig) GetContext() context.Context {
	return ctx.Context
}

func ConfigWithCancel(ctx ContextWithConfig) (ContextWithConfig, context.CancelFunc) {
	ctxWithCancel, cancel := context.WithCancel(ctx.GetContext())

	return newContextWithConfig(ctxWithCancel), cancel
}

func ConfigWithValue(ctx ContextWithConfig, key, val interface{}) ContextWithConfig {
	ctxWithValue := context.WithValue(ctx.GetContext(), key, val)
	return newContextWithConfig(ctxWithValue)
}

func ConfigWithTimeout(ctx ContextWithConfig, timeout time.Duration) (ContextWithConfig, context.CancelFunc) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx.GetContext(), timeout)
	return newContextWithConfig(ctxWithTimeout), cancel
}

func ConfigWithDeadline(ctx ContextWithConfig, d time.Time) (ContextWithConfig, context.CancelFunc) {
	ctxWithDeadline, cancel := context.WithDeadline(ctx.GetContext(), d)
	return newContextWithConfig(ctxWithDeadline), cancel
}

func getChainID(c *big.Int) *big.Int {
	chainID := big.NewInt(0)
	if c != nil {
		chainID.Set(c)
	}

	return chainID
}
