package params

import (
	"context"
	"math/big"
	"time"
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
	contextWithConfig
}

func newContextWithForkFlags(ctx context.Context) *contextWithForkFlags {
	return &contextWithForkFlags{contextWithConfig{ctx}}
}

func NewContextWithBlock(c *ChainConfig, blockNum *big.Int) ContextWithForkFlags {
	return c.WithEIPsFlags(context.Background(), blockNum)
}

func New(ctx context.Context, c *ChainConfig, blockNum *big.Int) ContextWithForkFlags {
	return c.WithEIPsFlags(ctx, blockNum)
}

func (c *ChainConfig) WithEIPsFlags(ctx context.Context, blockNum *big.Int) ContextWithForkFlags {
	checkers := c.getCheckers()

	for i, checker := range checkers {
		ctx = context.WithValue(ctx, eipFlags[i].eipFlag, checker(blockNum))
		ctx = context.WithValue(ctx, eipFlags[i].isEIPFlag, checker)
	}

	ctx = context.WithValue(ctx, BlockNumber, blockNum)
	ctx = context.WithValue(ctx, ChainID, getChainID(c.ChainID))

	return newContextWithForkFlags(ctx)
}

func (ctx contextWithForkFlags) GetForkFlag(name configKey) bool {
	b := ctx.Value(name)
	if b == nil {
		return false
	}
	if valB, ok := b.(bool); ok {
		return valB
	}
	return false
}

func (ctx contextWithForkFlags) GetBlockNumber() *big.Int {
	return getBigInt(ctx, BlockNumber)
}

func WithCancel(ctx contextWithForkFlags) (ContextWithForkFlags, context.CancelFunc) {
	ctxWithCancel, cancel := context.WithCancel(ctx.Context)

	return newContextWithForkFlags(ctxWithCancel), cancel
}

func WithValue(ctx contextWithForkFlags, key, val interface{}) ContextWithForkFlags {
	ctxWithValue := context.WithValue(ctx.Context, key, val)
	return newContextWithForkFlags(ctxWithValue)
}

func WithTimeout(ctx contextWithForkFlags, timeout time.Duration) (ContextWithForkFlags, context.CancelFunc) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx.Context, timeout)
	return newContextWithForkFlags(ctxWithTimeout), cancel
}

func WithDeadline(ctx contextWithForkFlags, d time.Time) (ContextWithForkFlags, context.CancelFunc) {
	ctxWithDeadline, cancel := context.WithDeadline(ctx.Context, d)
	return newContextWithForkFlags(ctxWithDeadline), cancel
}

func getBigInt(ctx context.Context, key configKey) *big.Int {
	b := ctx.Value(key)
	if b == nil {
		return nil
	}
	if valB, ok := b.(*big.Int); ok {
		return valB
	}
	return nil
}

func getForkFunc(ctx context.Context, name configKey) func(num *big.Int) bool {
	b := ctx.Value(name)
	if b == nil {
		return nil
	}
	if valB, ok := b.(func(num *big.Int) bool); ok {
		return valB
	}
	return nil
}
