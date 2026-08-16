# APT 可行性验证记录

## 当前已验证的层级

1. Go 编译：服务端与客户端均可构建。
2. 协议帧：AUTH、OPEN、OPEN_OK、DATA、CLOSE、ERROR 已有基础实现。
3. TLS 1.3：APT Session 建立在 TLS 传输上。
4. SOCKS5：客户端可以把本地代理连接转发到 APT 服务端。
5. Multiplexing（多路复用）：Session 层已经具备多个 Stream ID 的分发模型。
6. WINDOW_UPDATE（窗口更新）：协议负载编码/解析已经实现。
7. SendWindow（发送窗口）：加入线程安全的窗口额度、更新、最大窗口和溢出检查。

## 当前验证边界

GitHub Connector 可以验证源码提交和仓库结构，但不能代替真实的两机网络集成测试。因此没有把未执行的网络测试标记为 PASS。

## 下一阶段自动验证

- `go test ./...`
- `go test -race ./...`
- `go vet ./...`
- 双端 TLS 集成测试
- 多 Stream 并发传输测试
- 窗口耗尽/恢复测试
- 慢接收端 Backpressure（反压）测试
- Stream FIN/RST/timeout 测试

通过这些测试后，再进入 UDP Datagram（UDP 数据报）实现。
