# APT Development Notes

## Current implementation boundary

APT/0.1 MVP currently uses one TLS 1.3 connection per SOCKS5 CONNECT. The wire format already carries a stream ID so the protocol can evolve to true multiplexing, but the MVP deliberately avoids concurrent reads from a shared stream decoder.

This is intentional: correctness and a deterministic test target come before multiplexing complexity.

## Next implementation milestone

Introduce a single connection reader/dispatcher:

```text
TLS connection
      |
      v
Frame decoder
      |
      +---- stream 1 -> channel 1
      +---- stream 3 -> channel 3
      +---- stream 5 -> channel 5
```

Each logical stream then gets independent flow control and half-close semantics.

## Transport roadmap

- `transport/tcp_tls`: current
- `transport/h2`: HTTP/2 carrier
- `transport/quic`: QUIC carrier
- `transport/h3`: HTTP/3 carrier
- `transport/masque`: CONNECT-UDP / CONNECT-IP adapters
