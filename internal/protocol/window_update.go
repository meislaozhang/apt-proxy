package protocol

import "encoding/binary"

// EncodeWindowUpdate encodes a positive byte-credit increment.
func EncodeWindowUpdate(increment uint64) []byte {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, increment)
    return b
}

// DecodeWindowUpdate decodes a WINDOW_UPDATE increment.
func DecodeWindowUpdate(b []byte) (uint64, bool) {
    if len(b) != 8 { return 0, false }
    n := binary.BigEndian.Uint64(b)
    return n, n != 0
}
