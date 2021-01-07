package base62

import (
	"testing"
)

func TestEncode(t *testing.T) {
	assertEqual(t, "0", Encode(uint64(0)))
	assertEqual(t, "a", Encode(uint64(10)))
	assertEqual(t, "1C", Encode(uint64(100)))
	assertEqual(t, "aUKYOz", Encode(uint64(9999999999)))
}

func TestDecode(t *testing.T) {
	assertEqual(t, uint64(0), Decode("0"))
	assertEqual(t, uint64(10), Decode("a"))
	assertEqual(t, uint64(100), Decode("1C"))
	assertEqual(t, uint64(0x333f0966f7f), Decode("ZZZZZZZ"))
}

func assertEqual(t *testing.T, want, got interface{}) {
	t.Helper()
	if want != got {
		t.Errorf("expected %v, got %v", want, got)
	}
}
