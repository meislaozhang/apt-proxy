# APT Roadmap

APT (Adaptive Proxy Transport) is an experimental transport architecture.

## Milestones

- [x] APT/0.1 frame format
- [x] TLS 1.3 transport
- [x] Token authentication
- [x] SOCKS5 CONNECT
- [ ] APT/0.2 multiplexed streams
- [ ] UDP datagrams
- [ ] HTTP/2 transport
- [ ] QUIC transport
- [ ] HTTP/3 transport
- [ ] CONNECT-UDP
- [ ] CONNECT-IP
- [ ] Public-key authentication
- [ ] replay protection
- [ ] connection migration
- [ ] adaptive transport selection
- [ ] TUN/IP tunnel
- [ ] fuzzing and race tests
- [ ] interoperability test suite
- [ ] cross-platform release artifacts
- [ ] protocol freeze for APT/1.0

## Design rule

Implement standard transports instead of inventing look-alike cryptographic or browser-fingerprint layers. APT should compose TLS 1.3, HTTP/2, HTTP/3, QUIC and MASQUE where appropriate.

This roadmap is intentionally explicit: an unchecked feature is not considered implemented merely because the architecture anticipates it.
