# APT v0.1.0 Release Notes

## Status

APT v0.1.0 is an **Experimental（实验性）** development release.

It is intended to establish a reproducible baseline for the protocol and implementation. It is not a production security product and has not received an independent security audit.

## Included

- TLS 1.3（传输层安全）
- APT/0.1 framed protocol（协议帧）
- Token Authentication（令牌认证）
- Stream ID（数据流标识）
- TCP CONNECT forwarding（TCP 目标转发）
- Local SOCKS5 interface（本地 SOCKS5 接口）
- Go unit tests
- race detector CI
- go vet CI
- Linux build CI
- Chinese server/client documentation

## Known limitations

The following are deliberately **not** part of the v0.1.0 stability claim:

- Full Multiplexing（完整多路复用）
- UDP Datagram（UDP 数据报）
- HTTP/2 transport（HTTP/2 传输）
- QUIC / HTTP/3 transport（QUIC / HTTP/3 传输）
- CONNECT-UDP
- CONNECT-IP
- Public-key authentication（公钥认证）
- Connection Migration（连接迁移）
- Adaptive Transport（自适应传输）
- TUN/IP tunnel（TUN/IP 隧道）
- independent security audit（独立安全审计）

## Chinese documentation

See:

`docs/使用说明_CN.md`

## Recommended use

Use v0.1.0 for local development, protocol experiments, interoperability work and benchmarks. Do not use it as a security boundary for sensitive production traffic.
