# APT/5.0 下一道 Gate（闸门）

当前策略：每完成一个 Gate，先验证，再立即进入下一个目标；未通过不得标记 PASS。

## Gate 1：TCP Core E2E

必须全部通过：

- TLS 1.3 + authentication（认证）
- 8/100/1000 multiplexed streams（多路复用流）
- DATA integrity（数据完整性）
- Flow Control（流量控制）
- Backpressure（反压）
- Half-close（半关闭）
- FIN/CLOSE
- RST
- Race test（竞争检测）
- netem delay/loss/reorder（网络仿真）

## Gate 2：UDP Datagram

验证 UDP echo、并发、超时、资源限制和大规模数据报。

## Gate 3：QUIC Transport（传输层）

验证 QUIC 建连、APT Session 映射、stream/datagram、连接迁移。

## Gate 4：HTTP/3 Transport

验证 HTTP/3 承载 APT、重连和异常恢复。

## Gate 5：CONNECT-UDP / CONNECT-IP

验证 UDP 与 IP packet tunnel（IP 数据包隧道）。

## Gate 6：Security

Replay Protection（重放保护）、公钥身份、密钥轮换、Fuzz Testing（模糊测试）。

## Gate 7：Real-network

真实 VPS、IPv4/IPv6、高 RTT、丢包、带宽限制和长时间稳定性。

## Gate 8：Performance / Interop

Benchmark（性能基准）、1000+ streams、CPU/内存、跨平台和互操作测试。

## Gate 9：Release

Protocol Freeze（协议冻结）、兼容性承诺、独立安全审查，然后发布 APT/5.0.0。
