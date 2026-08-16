# APT 开发路线（中文）

APT（Adaptive Proxy Transport，自适应代理传输）采用分阶段实现，避免把未验证功能直接标成稳定版。

## 已完成的基础能力

- TLS 1.3（传输层安全）
- APT Frame（协议帧）
- Token Authentication（令牌认证）
- Stream ID（数据流标识）
- 服务端多 Stream 并发处理
- 客户端共享 Session（会话）的 SOCKS5 连接入口
- `go test`、`go test -race`、`go vet` 基础 CI

## 下一阶段

### 1. Multiplexing（多路复用）完善

- Stream 生命周期状态机
- FIN / RST（正常结束 / 强制复位）
- Flow Control（流量控制）
- Backpressure（反压）
- 最大并发 Stream 限制
- Session 和 Stream 超时

### 2. UDP Datagram（UDP 数据报）

增加独立 Datagram Frame（数据报帧），不把 UDP 强行模拟成 TCP 字节流。

### 3. QUIC / HTTP/3

使用成熟 QUIC 实现承载 APT，复用 QUIC 的拥塞控制、可靠 Stream 和 Datagram 能力。

### 4. HTTP/2

增加标准 HTTP/2 Extended CONNECT（扩展 CONNECT）传输。

### 5. MASQUE

实现 CONNECT-UDP 和 CONNECT-IP，并分别进行互操作测试。

### 6. Authentication（认证）

在 TLS Server Authentication（TLS 服务端认证）之外增加可轮换的公钥身份认证和 Replay Protection（重放保护）。

### 7. Connection Migration（连接迁移）

优先依赖 QUIC 的 Path Validation（路径验证）和连接迁移机制，避免重复实现一个不成熟的 UDP 迁移协议。

### 8. Adaptive Transport（自适应传输）

根据 RTT（往返时延）、Loss（丢包）、Jitter（抖动）、Throughput（吞吐量）和连接建立时间选择传输方式。

### 9. TUN / IP Tunnel（TUN / IP 隧道）

最后才实现完整 IP 层代理，配合 CONNECT-IP 做端到端测试。

## 发布标准

APT 在每个阶段都必须通过：

- 单元测试
- Race Detector（数据竞争检测）
- 集成测试
- Fuzz Testing（模糊测试）
- 资源耗尽测试
- 跨平台构建
- Benchmark（性能基准测试）
- 安全审查

未经这些验证的功能不会进入 Stable（稳定版）声明。
