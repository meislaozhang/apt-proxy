# APT 目标进度（2026-08-17）

## 已完成/已有原型

- TLS 1.3（传输层安全）基础传输
- APT Frame（协议帧）
- Token Authentication（令牌认证）
- TCP CONNECT（TCP 目标连接）
- SOCKS5 本地接口
- Stream ID（数据流标识）
- 基础 Session/Stream 多路复用模型
- WINDOW_UPDATE（窗口更新）协议负载
- SendWindow（发送窗口）基础状态
- 单元测试与 race detector（数据竞争检测）CI
- 服务端/客户端构建 CI
- 中文使用说明与发布流程文档

## 尚未完成

### APT/0.2：可靠 TCP 多路复用

- DATA 路径接入 SendWindow
- Backpressure（反压）
- ReceiveWindow（接收窗口）
- Session Window（会话级窗口）
- Half-close（半关闭）
- FIN/RST 完整状态机
- Keepalive（保活）/PING/PONG
- Stream/Session timeout（超时）
- 并发、慢消费者和资源耗尽测试

### APT/0.3：UDP

- Datagram（数据报）帧
- UDP 生命周期
- UDP idle timeout（空闲超时）
- UDP 资源限制
- TCP/UDP unified session（统一会话）

### APT/0.4：标准传输

- HTTP/2 Transport（HTTP/2 传输）
- QUIC Transport（QUIC 传输）
- HTTP/3 Transport（HTTP/3 传输）
- 传输层抽象与互操作测试

### APT/0.5：IP/UDP 隧道

- CONNECT-UDP
- CONNECT-IP
- MASQUE 互操作
- IP packet（IP 数据包）封装与 MTU/PMTU 处理

### APT/0.6：身份与安全

- Public-key Authentication（公钥认证）
- Ed25519 身份
- Key rotation（密钥轮换）
- Replay Protection（重放保护）
- Credential lifecycle（凭据生命周期）

### APT/0.7：连接迁移

- Connection Migration（连接迁移）
- Path Validation（路径验证）
- 网络切换恢复
- NAT rebinding（NAT 重新绑定）处理

### APT/0.8：自适应传输

- RTT/loss/jitter measurement（延迟/丢包/抖动测量）
- Transport selection（传输选择）
- 动态 fallback（回退）
- 传输切换一致性测试

### APT/0.9：完整 IP 代理

- TUN interface（TUN 虚拟网卡）
- IPv4/IPv6 packet forwarding（数据包转发）
- DNS handling（DNS 处理）
- MTU/fragmentation（分片）处理

### APT/1.0：发布质量

- Fuzz Testing（模糊测试）
- Benchmark（性能基准）
- 多平台构建：Linux/Windows/macOS，AMD64/ARM64
- 双端 integration test（集成测试）
- interoperability test（互操作测试）
- failure injection（故障注入）
- memory/resource exhaustion tests（资源耗尽测试）
- protocol freeze（协议冻结）
- Independent Security Review（独立安全审查）
- 正式 Release（正式版本）

> 当前仍为 Experimental（实验性）。未经过实际运行测试或独立安全审查的功能不能标记为稳定功能。