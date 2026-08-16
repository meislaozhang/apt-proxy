# APT — Adaptive Proxy Transport

APT (Adaptive Proxy Transport，自适应代理传输) is an experimental open-source proxy transport protocol.

## Status

**Experimental / research prototype.** The current MVP implements:

- TLS 1.3 transport
- APT/0.1 framed session
- token authentication
- stream IDs for protocol evolution
- TCP CONNECT forwarding
- local SOCKS5 interface

Planned phases:

1. True multiplexing
2. HTTP/2 transport
3. QUIC / HTTP/3 transport
4. UDP datagrams
5. CONNECT-UDP (RFC 9298)
6. CONNECT-IP (RFC 9484)
7. adaptive transport selection
8. connection migration
9. TUN/IP tunnel

## Build

```bash
go test ./...
go vet ./...
go build ./cmd/apt-server
go build ./cmd/apt-client
```

## Development test

Generate a local certificate:

```bash
openssl req -x509 -newkey rsa:2048 -nodes -days 7 \
  -keyout server.key -out server.crt -subj '/CN=localhost'
```

Start the server:

```bash
./apt-server -listen 127.0.0.1:8443 -cert server.crt -key server.key -token dev-token
```

Start the client:

```bash
./apt-client -listen 127.0.0.1:1080 -server 127.0.0.1:8443 -server-name localhost -token dev-token -insecure-skip-verify
```

Then configure a program to use SOCKS5 at `127.0.0.1:1080`.

## Security warning

This is not production-ready. Do not use it as a security boundary until the protocol, authentication, concurrency model, replay resistance, traffic handling, resource limits, fuzzing and independent cryptographic/security review are complete.

## License

Apache-2.0.
