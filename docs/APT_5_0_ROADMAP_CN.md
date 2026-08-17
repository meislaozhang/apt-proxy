# APT/5.0 快速推进路线

APT/5.0 是面向生产级代理传输的目标版本，而不是简单把版本号从 0.x 改成 5.0。协议必须经过兼容性、性能和安全验证后才能冻结。

## 目标架构

```text
                    APT/5.0
                       |
       +---------------+----------------+
       |               |                |
  TCP/TLS         QUIC/HTTP3       Future transports
       |               |
       +------- Session Layer ---------+
                       |
          +------------+-------------+
          |            |             |
        Stream      Datagram       IP Tunnel
          |            |             |
       Flow Ctrl    UDP flow       CONNECT-IP
          |
   Backpressure / FIN / RST
                       |
          Identity + Replay Protection
                       |
       Migration + Path Validation
                       |
            Adaptive Transport
                       |
       Fuzz / Benchmark / Interop
```

## 快速推进原则

1. 先稳定 Session/Stream 状态机，再扩展 UDP/IP。
2. Capability negotiation（能力协商）用于渐进式功能，不随意修改 wire version（线协议版本）。
3. 所有新帧必须有未知类型处理策略和长度上限。
4. Flow control（流量控制）必须同时覆盖 Stream 和 Session。
5. QUIC/HTTP3 优先复用成熟传输实现，不重新实现密码学或 QUIC。
6. 公钥身份、重放保护、连接迁移必须通过独立测试后才能稳定。
7. 任何“完成”必须同时区分源码实现、单元测试、集成测试和真实网络验证。

## 阶段

- 5.0-A：TCP Session + 完整 Flow Control
- 5.0-B：UDP Datagram + Keepalive
- 5.0-C：QUIC/HTTP3 Transport
- 5.0-D：CONNECT-UDP / CONNECT-IP
- 5.0-E：身份、密钥轮换、重放保护
- 5.0-F：Migration / Path Validation
- 5.0-G：Adaptive Transport
- 5.0-H：TUN + IPv4/IPv6
- 5.0-I：Fuzz / Benchmark / Interop / Security Review
- 5.0-Final：Protocol Freeze + Release

当前工作仍集中在 5.0-A，不能把尚未实现的阶段标记为完成。
