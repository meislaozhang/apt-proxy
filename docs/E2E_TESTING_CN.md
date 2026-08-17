# APT 端到端（E2E）验证方法

端到端验证不是只测试 Frame（协议帧）编码，而是实际启动 TLS listener（TLS 监听器）、APT server、测试目标服务和 APT client，然后验证完整数据路径。

## 当前 TCP E2E 测试

`internal/server/e2e_test.go` 执行以下链路：

```text
Test client
   │
   │ TLS 1.3
   ▼
APT Server
   │
   │ Stream #1..#8
   ▼
TCP Echo Server
   │
   │ echo
   ▼
APT Server
   │
   │ DATA frames
   ▼
APT Client
```

测试同时打开 8 个 Stream（数据流），每个 Stream 发送独立 payload，并验证返回内容与发送内容完全一致。

## 通过条件

- TLS 1.3 能完成握手
- Token Authentication（令牌认证）成功
- OPEN/OPEN_OK 成功
- 8 个 Stream 可同时存在
- DATA 能从 client 到 origin
- DATA 能从 origin 返回 client
- 每个 Stream 的数据不会串流
- Close 后资源能够释放

## 流控 E2E 下一步

TCP 多路复用基础链路通过后，再加入：

1. 小窗口配置。
2. 持续发送大于窗口的数据。
3. 确认发送端产生 Backpressure（反压）。
4. 接收端消费数据。
5. 发送 WINDOW_UPDATE。
6. 确认发送继续。
7. 慢消费者与资源上限测试。

## CI 验证

GitHub Actions 的 `test` workflow 会执行项目测试、race detector（数据竞争检测）、vet 和构建。当前最新 capability 测试提交已经有一次 `success` 的 workflow run；E2E 测试加入后必须等待新的 workflow run 成功，才能将 E2E 标记为 PASS。
