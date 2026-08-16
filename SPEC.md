# APT/0.1 Protocol Draft

## Frame

All frames use a 12-byte header:

- version: 1 byte
- type: 1 byte
- flags: 1 byte
- reserved: 1 byte
- stream_id: 4 bytes, network byte order
- payload_length: 4 bytes, network byte order

Payload is limited to 1 MiB in the MVP.

## Types

- 1 AUTH
- 2 AUTH_OK
- 3 OPEN
- 4 OPEN_OK
- 5 DATA
- 6 CLOSE
- 7 ERROR

## Session

The client establishes TLS 1.3, sends AUTH, and waits for AUTH_OK. Each OPEN creates a logical stream identified by stream_id. DATA carries bytes for that stream.

This draft deliberately avoids inventing cryptography. Confidentiality and server authentication are delegated to TLS 1.3. The token is an application-level credential and is not a replacement for TLS.

## Design principles

APT is intended to separate a stable proxy/session layer from multiple standard transports. Future versions may carry the same logical session over TCP/TLS, HTTP/2, QUIC/HTTP/3 and standardized HTTP proxy mechanisms such as CONNECT-UDP and CONNECT-IP.
