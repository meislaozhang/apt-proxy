package protocol

import (
    "bufio"
    "encoding/binary"
    "fmt"
    "io"
)

const (
    Version    byte = 1
    TypeAuth   byte = 1
    TypeAuthOK byte = 2
    TypeOpen   byte = 3
    TypeOpenOK byte = 4
    TypeData   byte = 5
    TypeClose  byte = 6
    TypeError  byte = 7
)

const MaxPayload = 1 << 20

type Frame struct {
    Type     byte
    Flags    byte
    StreamID uint32
    Payload  []byte
}

func (f Frame) Write(w io.Writer) error {
    if len(f.Payload) > MaxPayload {
        return fmt.Errorf("payload too large: %d", len(f.Payload))
    }
    h := make([]byte, 12)
    h[0] = Version
    h[1] = f.Type
    h[2] = f.Flags
    binary.BigEndian.PutUint32(h[4:8], f.StreamID)
    binary.BigEndian.PutUint32(h[8:12], uint32(len(f.Payload)))
    if _, err := w.Write(h); err != nil { return err }
    _, err := w.Write(f.Payload)
    return err
}

func ReadFrame(r *bufio.Reader) (Frame, error) {
    var h [12]byte
    if _, err := io.ReadFull(r, h[:]); err != nil { return Frame{}, err }
    if h[0] != Version { return Frame{}, fmt.Errorf("unsupported version %d", h[0]) }
    n := binary.BigEndian.Uint32(h[8:12])
    if n > MaxPayload { return Frame{}, fmt.Errorf("payload too large: %d", n) }
    p := make([]byte, n)
    if _, err := io.ReadFull(r, p); err != nil { return Frame{}, err }
    return Frame{Type: h[1], Flags: h[2], StreamID: binary.BigEndian.Uint32(h[4:8]), Payload: p}, nil
}
