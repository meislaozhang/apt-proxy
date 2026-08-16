# APT 中文使用说明

APT（Adaptive Proxy Transport，自适应代理传输）是一个实验性开源代理传输协议实现。

> 当前版本定位：实验性 MVP（最小可行版本）。它用于验证协议设计，不应被视为经过独立安全审计的生产级安全边界。

## 1. 当前版本能力

当前版本主要提供：

- TLS 1.3（传输层安全）
- APT/0.1 Frame（协议帧）
- Token Authentication（令牌认证）
- Stream ID（数据流标识）
- TCP CONNECT（TCP 目标连接）
- 本地 SOCKS5（SOCKS5 代理接口）

当前还**不应宣称**已经完成：UDP、QUIC、HTTP/3、CONNECT-UDP、CONNECT-IP、TUN、连接迁移和自适应传输。

## 2. 编译

要求 Go（Golang，Go 编程语言）环境。

```bash
go test ./...
go vet ./...
go build -o apt-server ./cmd/apt-server
go build -o apt-client ./cmd/apt-client
```

也可以：

```bash
make test
make build
```

## 3. 服务端配置

APT 服务端负责监听 TLS 端口并把 APT OPEN（打开流）请求转发到目标 TCP 地址。

### 3.1 生成测试证书

仅用于本机实验：

```bash
openssl req -x509 -newkey rsa:2048 -nodes -days 7 \
  -keyout server.key \
  -out server.crt \
  -subj '/CN=localhost'
```

### 3.2 启动服务端

```bash
./apt-server \
  -listen 0.0.0.0:8443 \
  -cert server.crt \
  -key server.key \
  -token '请替换成随机长令牌'
```

参数：

| 参数 | 说明 |
|---|---|
| `-listen` | 监听地址，例如 `0.0.0.0:8443` |
| `-cert` | TLS 证书文件 |
| `-key` | TLS 私钥文件 |
| `-token` | 应用层认证令牌 |

生产环境不要使用 `change-me`、`dev-token` 等示例令牌。

## 4. 客户端配置

客户端在本机监听 SOCKS5，然后通过 APT Server（APT 服务端）建立 TLS 1.3 连接。

```bash
./apt-client \
  -listen 127.0.0.1:1080 \
  -server 203.0.113.10:8443 \
  -server-name apt.example.com \
  -token '你的随机长令牌' \
  -insecure-skip-verify=false
```

参数：

| 参数 | 说明 |
|---|---|
| `-listen` | 本地 SOCKS5 监听地址 |
| `-server` | APT 服务端地址和端口 |
| `-server-name` | TLS Server Name（服务器名称），用于证书校验/SNI |
| `-token` | 与服务端一致的认证令牌 |
| `-insecure-skip-verify` | 是否跳过 TLS 证书验证；生产环境必须关闭 |

## 5. 使用 SOCKS5

客户端启动后：

```text
应用程序
   ↓
SOCKS5 127.0.0.1:1080
   ↓
APT Client
   ↓
TLS 1.3
   ↓
APT Server
   ↓
目标 TCP 服务
```

例如：

```bash
curl --proxy socks5h://127.0.0.1:1080 https://example.com/
```

`socks5h` 表示让代理端处理域名解析；具体 DNS 行为仍取决于客户端/目标连接实现。

## 6. TLS 证书

测试时可以使用自签名证书，但生产环境应使用正规 CA（Certificate Authority，证书颁发机构）签发的证书，并让客户端正常验证证书。

推荐：

```text
-client -insecure-skip-verify=false
-server-name 与证书 SAN（Subject Alternative Name，主题备用名称）匹配
```

不要把 `-insecure-skip-verify=true` 当作正式安全配置。

## 7. systemd 服务

在真正的 Linux VPS（Virtual Private Server，虚拟专用服务器）上可以使用 systemd（系统服务管理器）。

服务端示例：

```ini
[Unit]
Description=APT Server
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=/opt/apt/apt-server -listen 0.0.0.0:8443 -cert /opt/apt/server.crt -key /opt/apt/server.key -token YOUR_TOKEN
Restart=on-failure
RestartSec=2

[Install]
WantedBy=multi-user.target
```

保存为：

```text
/etc/systemd/system/apt-server.service
```

然后：

```bash
systemctl daemon-reload
systemctl enable --now apt-server
systemctl status apt-server
```

## 8. Docker / 容器环境

如果 PID 1（进程 ID 1）不是 systemd，不要直接依赖 `systemctl`。

例如容器中可以直接运行：

```bash
exec /opt/apt/apt-server -listen 0.0.0.0:8443 ...
```

或者由容器编排系统负责 Restart（自动重启）和开机启动。

## 9. 安全注意事项

APT 当前仍是实验性协议。

不要在当前版本中假设它已经具备：

- 完整抗重放机制
- 完整资源耗尽防护
- 完整流量控制
- 独立密码学审计
- 完整协议模糊测试覆盖
- 完整 QUIC/HTTP3 安全模型
- 完整匿名性或抗封锁能力

APT 不应被描述为“隐身协议”或“不会被封锁的协议”。

## 10. 开发测试

```bash
go test ./...
go test -race ./...
go vet ./...
go build ./cmd/apt-server
ngo build ./cmd/apt-client
```

完整功能完成并通过独立安全审查之前，版本应保持 Experimental（实验性）状态。
