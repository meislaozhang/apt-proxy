# APT/5.0 开发进度

## Gate 1 — TCP Core

### 已实现
- TLS 1.3 client/server
- Token Authentication（令牌认证）
- TCP multiplexing（TCP 多路复用）
- Stream open/close
- HALF_CLOSE frame
- RESET frame
- Client `CloseWrite()` / `Reset()`
- Server half-close/reset handling
- E2E multiplexing test
- E2E half-close/reset tests
- Client close ordering fix: buffered DATA is consumed before EOF after CLOSE/RESET

### 尚未 PASS
- Flow Control E2E
- Backpressure E2E
- FIN/CLOSE complete state machine
- Race test
- netem delay/loss/reorder
- 100/1000 stream stress
- large-transfer integrity

## Gate 2 — UDP Datagram

未开始。Gate 1 必须通过后立即进入。

## Gate 3 — QUIC

未开始。

## Gate 4 — HTTP/3

未开始。

## Gate 5 — CONNECT-UDP / CONNECT-IP

未开始。

## Gate 6 — Security

Replay Protection（重放保护）、public-key identity（公钥身份）、key rotation（密钥轮换）、fuzzing（模糊测试）。

## Gate 7 — Real Network

真实 VPS、IPv4/IPv6、延迟、丢包、长期运行。

## Gate 8 — Performance / Interop

吞吐、CPU、内存、1000+ streams、跨实现互操作。

## Gate 9 — Release

Protocol Freeze（协议冻结）、兼容性保证、APT/5.0.0。

> 规则：只有有可重复的测试证据才标记 PASS；提交代码本身不算 PASS。
