package auth

import "testing"

func TestEqualToken(t *testing.T) {
    if !EqualToken([]byte("abc"), []byte("abc")) { t.Fatal("equal token rejected") }
    if EqualToken([]byte("abc"), []byte("abd")) { t.Fatal("different token accepted") }
}
