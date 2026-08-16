package protocol

import (
    "encoding/binary"
    "fmt"
)

// WindowUpdatePayload encodes a stream-level receive-window increment.
func WindowUpdatePayload(n uint32) []byte {
    b := make([]byte, 4)
    binary.BigEndian.PutUint32(b, n)
    return b
}

func ParseWindowUpdate(b []byte) (uint32, error) {
    if len(b) != 4 {
        return 0, fmt.Errorf("APT: invalid WINDOW_UPDATE payload length %d", len(b))
    }
    n := binary.BigEndian.Uint32(b)
    if n == 0 {
        return 0, fmt.Errorf("APT: WINDOW_UPDATE increment must be non-zero")
    }
    return n, nil
}
