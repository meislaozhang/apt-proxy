package protocol

import (
    "encoding/binary"
    "fmt"
)

// Capabilities are negotiated during HELLO/HELLO_OK. Unknown bits are ignored
// by peers, allowing additive extension without changing the wire major version.
type Capabilities struct {
    Bits uint64
}

func (c Capabilities) Has(bit uint64) bool { return c.Bits&bit != 0 }
func (c *Capabilities) Enable(bit uint64) { c.Bits |= bit }

func (c Capabilities) Encode() []byte {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, c.Bits)
    return b
}

func DecodeCapabilities(b []byte) (Capabilities, error) {
    if len(b) != 8 {
        return Capabilities{}, fmt.Errorf("invalid capabilities payload length %d", len(b))
    }
    return Capabilities{Bits: binary.BigEndian.Uint64(b)}, nil
}
