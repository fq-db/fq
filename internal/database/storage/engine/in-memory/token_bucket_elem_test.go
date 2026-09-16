package inmemory

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fq-db/fq/internal/database"
)

func TestTokenBucketRefillSaturatesWithoutOverflow(t *testing.T) {
	const capacity = database.ValueType(math.MaxInt64)
	const refillAmount = database.ValueType(math.MaxInt64 / 2)

	e := NewTokenBucketElem(1)
	e.Restore(database.DumpElem{Value: 1, TxAt: 1, Tx: 5})

	result, err := e.RLimit(
		database.TxContext{Tx: 6, DumpTx: database.NoTx, CurrTime: 4_000_000_000},
		capacity,
		refillAmount,
		nil,
	)
	require.NoError(t, err)
	require.True(t, result.Allowed)
	require.Equal(t, database.ValueType(1), result.Current)
	require.Equal(t, capacity-1, result.Remaining)
}

func TestTokenBucketRefillAddsPartialPeriods(t *testing.T) {
	e := NewTokenBucketElem(10)
	e.Restore(database.DumpElem{Value: 2, TxAt: 100, Tx: 5})

	result, err := e.RLimit(
		database.TxContext{Tx: 6, DumpTx: database.NoTx, CurrTime: 130},
		100,
		5,
		nil,
	)
	require.NoError(t, err)
	require.True(t, result.Allowed)
	require.Equal(t, database.ValueType(16), result.Remaining)
}
