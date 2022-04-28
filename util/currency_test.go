package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSupportedCurrency(t *testing.T) {
	currency := "IDR"
	require.True(t, IsSupportedCurrency(currency))

	currency2 := "ASD"
	require.False(t, IsSupportedCurrency(currency2))
}
