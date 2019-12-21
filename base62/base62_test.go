package base62

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	assert.Equal(t, "0", Encode(uint64(0)))
	assert.Equal(t, "a", Encode(uint64(10)))
	assert.Equal(t, "1C", Encode(uint64(100)))
	assert.Equal(t, "aUKYOz", Encode(uint64(9999999999)))
}

func TestDecode(t *testing.T) {
	assert.Equal(t, uint64(0), Decode("0"))
	assert.Equal(t, uint64(10), Decode("a"))
	assert.Equal(t, uint64(100), Decode("1C"))
	assert.Equal(t, uint64(0x333f0966f7f), Decode("ZZZZZZZ"))
}
