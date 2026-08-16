# APT Protocol Notes

APT separates the session protocol from its transport.

## Session concepts

- `OPEN`: create a logical stream toward a target.
- `DATA`: carry stream bytes.
- `CLOSE`: terminate a logical stream.
- `DATAGRAM`: carry unreliable datagrams when the selected transport supports them.
- `PING` / `PONG`: liveness and path validation primitives.
- `AUTH`: authenticate a client before opening application streams.

## Transport mapping

APT is designed to run over:

1. TCP + TLS 1.3
2. HTTP/2 + TLS 1.3
3. QUIC + TLS 1.3
4. HTTP/3
5. MASQUE mechanisms such as CONNECT-UDP and CONNECT-IP where standards and implementations permit them.

The transport is an implementation detail of the session. APT does not require a custom TLS implementation or browser-fingerprint impersonation.

## Security goals

The protocol must provide confidentiality and integrity through a standard authenticated transport. Application authentication is separate from transport authentication. Replay protection, public-key authentication, and key rotation are future protocol requirements and must not be claimed as implemented until interoperable tests exist.

## Compatibility

The local interface can expose SOCKS5 and HTTP CONNECT. Future releases may expose TUN/IP tunnel mode without changing the application-facing session model.
