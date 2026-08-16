package proxy

import (
    "encoding/binary"
    "fmt"
    "io"
    "net"
)

func ServeSOCKS5(l net.Listener, connect func(net.Conn, string) error) error {
    for {
        c, err := l.Accept()
        if err != nil { return err }
        go func() { _ = handleSOCKS5(c, connect) }()
    }
}

func handleSOCKS5(c net.Conn, connect func(net.Conn, string) error) error {
    defer c.Close()
    var h [4]byte
    if _, err := io.ReadFull(c, h[:2]); err != nil { return err }
    if h[0] != 5 { return fmt.Errorf("not socks5") }
    methods := make([]byte, h[1])
    if _, err := io.ReadFull(c, methods); err != nil { return err }
    if _, err := c.Write([]byte{5, 0}); err != nil { return err }
    if _, err := io.ReadFull(c, h[:4]); err != nil { return err }
    if h[0] != 5 || h[1] != 1 { return fmt.Errorf("only CONNECT supported") }
    addr, err := readAddr(c, h[3])
    if err != nil { return err }
    if err = connect(c, addr); err != nil {
        _, _ = c.Write([]byte{5, 1, 0, 1, 0, 0, 0, 0, 0, 0})
        return err
    }
    _, err = c.Write([]byte{5, 0, 0, 1, 0, 0, 0, 0, 0, 0})
    return err
}

func readAddr(r io.Reader, atyp byte) (string, error) {
    switch atyp {
    case 1:
        b := make([]byte, 4); if _, e := io.ReadFull(r, b); e != nil { return "", e }
        var p [2]byte; if _, e := io.ReadFull(r, p[:]); e != nil { return "", e }
        return net.JoinHostPort(net.IP(b).String(), fmt.Sprint(binary.BigEndian.Uint16(p[:]))), nil
    case 3:
        var n [1]byte; if _, e := io.ReadFull(r, n[:]); e != nil { return "", e }
        b := make([]byte, n[0]); if _, e := io.ReadFull(r, b); e != nil { return "", e }
        var p [2]byte; if _, e := io.ReadFull(r, p[:]); e != nil { return "", e }
        return net.JoinHostPort(string(b), fmt.Sprint(binary.BigEndian.Uint16(p[:]))), nil
    case 4:
        b := make([]byte, 16); if _, e := io.ReadFull(r, b); e != nil { return "", e }
        var p [2]byte; if _, e := io.ReadFull(r, p[:]); e != nil { return "", e }
        return net.JoinHostPort(net.IP(b).String(), fmt.Sprint(binary.BigEndian.Uint16(p[:]))), nil
    default:
        return "", fmt.Errorf("unsupported address type %d", atyp)
    }
}
