# APT — Adaptive Proxy Transport

APT（Adaptive Proxy Transport，自适应代理传输）是一个实验性开源代理传输协议。

> **状态：Experimental（实验性）**。当前仓库用于协议和实现研究，不代表已经通过独立安全审计，也不应被描述为匿名、隐身或抗封锁协议。

## 当前 MVP（最小可行版本）

当前版本实现：

- TLS 1.3（传输层安全）
- APT/0.1 Frame（协议帧）
- Token Authentication（令牌认证）
- Stream ID（数据流标识）
- TCP CONNECT（TCP 目标连接）
- 本地 SOCKS5（SOCKS5 代理接口）

### 尚未作为稳定功能发布

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

这些能力只有在实现、互操作测试、资源限制、模糊测试和安全审查完成后才会进入稳定发布范围。

## 中文使用说明

完整的服务端和客户端中文部署、启动、TLS 证书、SOCKS5 和 systemd 示例见：

**[`docs/使用说明_CN.md`](docs/使用说明_CN.md)**

## 快速开始

### 编译

```bash
go test ./...
go vet ./...
go build -o apt-server ./cmd/apt-server
go build -o apt-client ./cmd/apt-client
```

### 本地测试证书

```bash
openssl req -x509 -newkey rsa:2048 -nodes -days 7 \
  -keyout server.key -out server.crt \
  -subj '/CN=localhost'
```

### 服务端

```bash
./apt-server \
  -listen 127.0.0.1:8443 \
  -cert server.crt \
  -key server.key \
  -token 'dev-token-change-me'
```

### 客户端

```bash
./apt-client \
  -listen 127.0.0.1:1080 \
  -server 127.0.0.1:8443 \
  -server-name localhost \
  -token 'dev-token-change-me' \
  -insecure-skip-verify=true
```

> `-insecure-skip-verify=true` 仅用于本地自签名证书测试。正式部署应使用受信任 CA（Certificate Authority，证书颁发机构）证书并设置为 `false`。

### SOCKS5 测试

```bash
curl --proxy socks5h://127.0.0.1:1080 https://example.com/
```

## 服务端/客户端架构

```text
应用程序
   │
   ▼
SOCKS5 :1080
   │
   ▼
APT Client
   │
   │ TLS 1.3
   ▼
APT Server :8443
   │
   ▼
Target TCP Service
```

APT 将逻辑代理会话与底层传输层分离，后续目标是让同一协议会话能够承载 TCP、UDP 和 IP，并适配 TCP/TLS、HTTP/2、QUIC/HTTP/3 等标准传输。

## 开发测试

```bash
go test ./...
go test -race ./...
go vet ./...
go build ./cmd/apt-server
go build ./cmd/apt-client
```

## Release（版本发布）原则

APT 在以下项目完成之前保持 Experimental（实验性）状态：

1. 协议状态机和并发模型稳定
2. 多路复用和流量控制经过测试
3. UDP/QUIC/HTTP3 互操作测试
4. 资源耗尽和异常输入防护
5. Fuzz Testing（模糊测试）
6. 跨平台构建与集成测试
7. Benchmark（性能基准测试）
8. 独立安全审查

版本历史见 [`CHANGELOG.md`](CHANGELOG.md)。

## Security（安全）

APT 当前不是生产级安全边界。TLS 负责传输加密和服务器身份验证；应用层 Token 是额外认证凭据，并不替代 TLS。

请参阅 [`SECURITY.md`](SECURITY.md)。

## License

Apache-2.0.
