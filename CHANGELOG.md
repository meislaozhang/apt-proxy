# Changelog

## [Unreleased]

### Added

- 中文服务端/客户端使用说明：`docs/使用说明_CN.md`
- 明确当前版本能力边界与安全状态。

### Current MVP

- TLS 1.3（传输层安全）
- APT/0.1 Frame（协议帧）
- Token Authentication（令牌认证）
- Stream ID（数据流标识）
- TCP CONNECT（TCP 连接代理）
- Local SOCKS5（本地 SOCKS5 接口）

### Not yet released

以下能力仍处于开发/设计阶段，未在当前版本中宣称完成：

- 完整 Multiplexing（多路复用）
- UDP Datagram（UDP 数据报）
- HTTP/2 Transport（HTTP/2 传输）
- QUIC / HTTP/3 Transport（QUIC / HTTP/3 传输）
- CONNECT-UDP
- CONNECT-IP
- Public-key Authentication（公钥认证）
- Connection Migration（连接迁移）
- Adaptive Transport（自适应传输）
- TUN/IP Tunnel（TUN/IP 隧道）
- Fuzz Testing（模糊测试）与独立安全审计

## Release policy

APT 在完成协议兼容性测试、资源限制、模糊测试、跨平台构建和安全审查之前，不发布生产级稳定版承诺。
