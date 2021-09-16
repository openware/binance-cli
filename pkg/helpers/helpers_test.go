package helpers

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValuePrecision(t *testing.T) {
	assert.Equal(t, 3, ValuePrecision(json.Number("0.00100000")))
	assert.Equal(t, 1, ValuePrecision(json.Number("0.1000000")))
	assert.Equal(t, 0, ValuePrecision(json.Number("1.0")))
	assert.Equal(t, -1, ValuePrecision(json.Number("10.0")))
	assert.Equal(t, -2, ValuePrecision(json.Number("100")))
	assert.Equal(t, -5, ValuePrecision(json.Number("100000")))
}
