# APT — Adaptive Proxy Transport

APT（Adaptive Proxy Transport，自适应代理传输）是一个基于 Go 实现的实验性代理传输协议与 TCP 代理实现。

> **当前里程碑：APT 5.0 开发目标已完成。** Core、并发/多路复用、流生命周期、流量控制、Race Detector、重复稳定性以及最终仓库级 CI Gate 均已通过。
>
> **发布定位仍为 Experimental（实验性）**：尚未完成独立安全审计、完整 Fuzz Testing（模糊测试）覆盖以及长期生产环境认证。因此不要将 APT 描述为匿名、隐身、抗封锁或经过独立安全审计的安全协议。

## 目录

- [特性](#特性)
- [快速开始](#快速开始)
- [服务端](#服务端)
- [客户端](#客户端)
- [SOCKS5 使用](#socks5-使用)
- [TLS 证书](#tls-证书)
- [systemd 部署](#systemd-部署)
- [Docker / 容器](#docker--容器)
- [验证安装](#验证安装)
- [架构](#架构)
- [安全说明](#安全说明)
- [当前范围与后续计划](#当前范围与后续计划)
- [开发测试](#开发测试)
- [版本](#版本)
- [License](#license)

## 特性

当前 APT 5.0 Core 主要覆盖：

- TLS 1.3（传输层安全）
- APT Frame（APT 协议帧）
- Token Authentication（令牌认证）
- Stream ID（数据流标识）
- TCP CONNECT（TCP 目标连接）
- 本地 SOCKS5（SOCKS5 代理接口）
- TCP Multiplexing（TCP 多路复用）
- 多 Stream 生命周期管理
- Flow Control（流量控制）与 Backpressure（反压）
- WINDOW_UPDATE
- RESET / HALF_CLOSE 生命周期处理
- Wire Lifecycle（传输连接生命周期）处理
- 共享 Multiplexed Session（多路复用会话）

APT 5.0 的核心验证已经覆盖 Functional / Integration Tests、Race Detector、`go vet`、Server/Client Build、Wire Lifecycle 以及重复稳定性 Gate。

## 快速开始

### 1. 获取源码

```bash
git clone https://github.com/meislaozhang/apt-proxy.git
cd apt-proxy
```

### 2. 编译

需要 Go（Golang）环境：

```bash
go test ./...
go vet ./...
go build -o apt-server ./cmd/apt-server
go build -o apt-client ./cmd/apt-client
```

也可以使用：

```bash
make test
make build
```

### 3. 生成本地测试证书

仅用于本机实验：

```bash
openssl req -x509 -newkey rsa:2048 -nodes -days 7 \
  -keyout server.key \
  -out server.crt \
  -subj '/CN=localhost'
```

### 4. 启动服务端

```bash
./apt-server \
  -listen 127.0.0.1:8443 \
  -cert server.crt \
  -key server.key \
  -token 'dev-token-change-me'
```

### 5. 启动客户端

```bash
./apt-client \
  -listen 127.0.0.1:1080 \
  -server 127.0.0.1:8443 \
  -server-name localhost \
  -token 'dev-token-change-me' \
  -insecure-skip-verify=true
```

> `-insecure-skip-verify=true` **只用于本地自签名证书测试**。正式部署必须使用受信任 CA（Certificate Authority，证书颁发机构）证书，并设置为 `false`。

### 6. 测试 SOCKS5

```bash
curl --proxy socks5h://127.0.0.1:1080 https://example.com/
```

`socks5h` 表示由代理端处理域名解析；实际 DNS 行为取决于客户端和目标连接实现。

## 服务端

服务端负责：

1. 监听 TLS 端口；
2. 校验应用层 Token；
3. 接收 APT Frame；
4. 创建/管理 Multiplexed Stream；
5. 将 TCP CONNECT 请求连接到目标地址；
6. 在多个 Stream 之间复用同一底层 TLS 会话。

### 服务端参数

| 参数 | 默认值 | 说明 |
|---|---|---|
| `-listen` | `:8443` | TLS 监听地址 |
| `-cert` | `server.crt` | TLS 证书文件 |
| `-key` | `server.key` | TLS 私钥文件 |
| `-token` | `change-me` | 应用层认证 Token |

生产环境必须替换默认 Token，建议使用随机生成的高熵长字符串。

示例：

```bash
./apt-server \
  -listen 0.0.0.0:8443 \
  -cert /opt/apt/server.crt \
  -key /opt/apt/server.key \
  -token 'REPLACE_WITH_A_LONG_RANDOM_TOKEN'
```

## 客户端

客户端负责：

1. 监听本地 SOCKS5；
2. 与 APT Server 建立 TLS 1.3 连接；
3. 复用一个底层 APT Session 承载多个 TCP Stream；
4. 将 SOCKS5 CONNECT 请求转换为 APT TCP Stream；
5. 处理 Stream 数据和生命周期。

### 客户端参数

| 参数 | 默认值 | 说明 |
|---|---|---|
| `-listen` | `127.0.0.1:1080` | 本地 SOCKS5 监听地址 |
| `-server` | `127.0.0.1:8443` | APT Server 地址 |
| `-server-name` | `localhost` | TLS ServerName / SNI |
| `-token` | `change-me` | 与服务端一致的认证 Token |
| `-insecure-skip-verify` | `true` | 是否跳过 TLS 证书验证；仅开发测试使用 |

正式部署建议：

```bash
./apt-client \
  -listen 127.0.0.1:1080 \
  -server apt.example.com:8443 \
  -server-name apt.example.com \
  -token 'YOUR_RANDOM_LONG_TOKEN' \
  -insecure-skip-verify=false
```

## SOCKS5 使用

启动客户端后，本机代理地址为：

```text
127.0.0.1:1080
```

### curl

```bash
curl --proxy socks5h://127.0.0.1:1080 https://example.com/
```

### 浏览器

将浏览器 SOCKS5 代理设置为：

```text
Host: 127.0.0.1
Port: 1080
```

如果希望域名由代理端处理，应优先使用支持 SOCKS5 hostname forwarding 的配置。

## TLS 证书

APT 使用 TLS 1.3 保护客户端与服务端之间的传输。

### 测试环境

可以使用自签名证书，并在客户端临时使用：

```text
-insecure-skip-verify=true
```

### 正式环境

应使用受信任 CA 签发的证书：

```text
-insecure-skip-verify=false
-server-name 与证书 SAN 匹配
```

例如证书包含：

```text
DNS:apt.example.com
```

客户端就使用：

```text
-server-name apt.example.com
```

不要在正式环境长期使用 `-insecure-skip-verify=true`。

## systemd 部署

Linux VPS（Virtual Private Server）上可以使用 systemd 管理服务端。

创建：

```text
/etc/systemd/system/apt-server.service
```

内容示例：

```ini
[Unit]
Description=APT Server
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=/opt/apt/apt-server \
  -listen 0.0.0.0:8443 \
  -cert /opt/apt/server.crt \
  -key /opt/apt/server.key \
  -token YOUR_LONG_RANDOM_TOKEN
Restart=on-failure
RestartSec=2

[Install]
WantedBy=multi-user.target
```

然后：

```bash
systemctl daemon-reload
systemctl enable --now apt-server
systemctl status apt-server
```

查看日志：

```bash
journalctl -u apt-server -f
```

## Docker / 容器

容器中通常没有 systemd，因此不要依赖 `systemctl`。

直接让 APT Server 成为前台进程：

```bash
exec /opt/apt/apt-server \
  -listen 0.0.0.0:8443 \
  -cert /opt/apt/server.crt \
  -key /opt/apt/server.key \
  -token "$APT_TOKEN"
```

由 Docker、Compose 或其他容器编排系统负责 Restart、日志和生命周期管理。

## 验证安装

### 基础测试

```bash
go test ./...
```

### Race Detector

```bash
go test -race ./...
```

### 静态检查

```bash
go vet ./...
```

### 构建

```bash
go build ./cmd/apt-server
go build ./cmd/apt-client
```

### 推荐本地验证顺序

```bash
go test ./... && \
go test -race ./... && \
go vet ./... && \
go build ./cmd/apt-server && \
go build ./cmd/apt-client
```

## 架构

```text
                    Application
                         │
                         ▼
                  SOCKS5 :1080
                         │
                         ▼
                    APT Client
                         │
                  ┌──────┴──────┐
                  │ APT Session │
                  │ TLS 1.3     │
                  └──────┬──────┘
                         │
              Multiplexed TCP Streams
             ┌───────────┼───────────┐
             ▼           ▼           ▼
          Stream 1    Stream 2    Stream N
             │           │           │
             └───────────┼───────────┘
                         │
                         ▼
                    APT Server
                         │
                         ▼
                  Target TCP Services
```

APT 的核心模型是将**逻辑 Stream**与**底层 Transport Session**分离，使多个 TCP Stream 可以共享同一个 TLS 会话。

## 安全说明

APT 当前仍然属于 Experimental（实验性）协议。

TLS 1.3 提供传输加密和服务器身份验证；应用层 Token 提供额外认证，但不能替代 TLS。

当前版本不应被描述为：

- 匿名协议；
- 隐身协议；
- 抗封锁协议；
- 保证绕过网络审查的协议；
- 已通过独立密码学安全审计的协议。

生产环境还应自行配置：

- 防火墙；
- 最小监听面；
- 强随机 Token；
- 正规 CA 证书；
- 日志与监控；
- 系统权限隔离；
- 资源限制。

## 当前范围与后续计划

APT 5.0 已完成当前定义的 TCP Core 开发和仓库级自动化验收。

以下能力仍不应作为当前稳定功能宣称：

- UDP Datagram；
- QUIC / HTTP/3 Transport；
- CONNECT-UDP；
- CONNECT-IP；
- Public-key Authentication；
- Connection Migration；
- Adaptive Transport；
- TUN/IP Tunnel；
- 独立安全审计；
- 完整 Fuzz Testing 覆盖；
- 长期生产环境认证。

这些能力应在实现、互操作测试、资源限制、模糊测试和安全审查完成后再进入正式稳定范围。

## 开发测试

完整测试：

```bash
go test ./...
go test -race ./...
go vet ./...
go build ./cmd/apt-server
go build ./cmd/apt-client
```

仓库还包含 TCP Core 的稳定性和最终验收 CI Gate，用于重复测试 Functional、Race、Vet 和 Build。

## 版本

当前开发里程碑：

```text
APT 5.0 — TCP Core Complete
```

APT 5.0 的开发目标及当前仓库级验收已经完成。正式 Release Tag / GitHub Release 应与实际发布动作保持一致，不在 README 中虚构尚未创建的版本号。

版本历史见 [`CHANGELOG.md`](CHANGELOG.md)。

## License

Apache-2.0.
