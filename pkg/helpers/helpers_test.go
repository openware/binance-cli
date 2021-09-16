package helpers

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestValuePrecision(t *testing.T) {
	assert.Equal(t, int64(3), ValuePrecision(decimal.RequireFromString("0.00100000")))
	assert.Equal(t,  int64(1), ValuePrecision(decimal.RequireFromString("0.1000000")))
	assert.Equal(t,  int64(0), ValuePrecision(decimal.RequireFromString("1.0")))
	assert.Equal(t,  int64(-1), ValuePrecision(decimal.RequireFromString("10.0")))
	assert.Equal(t,  int64(-2), ValuePrecision(decimal.RequireFromString("100")))
	assert.Equal(t,  int64(-5), ValuePrecision(decimal.RequireFromString("100000")))
}
