# APT GFW/真实网络测试计划

> 本文只定义合法的网络可靠性与协议鲁棒性测试，不提供规避网络审查或隐藏流量特征的操作步骤。

## 目标

验证 APT 在不同真实网络条件下是否稳定、是否正确处理 TLS/HTTP2/QUIC 等传输变化，以及丢包、延迟、重排、连接重置、MTU 变化等异常。

## 分层

### L0：可重复实验室测试

使用 Linux `tc netem` 模拟 RTT、loss、jitter、reorder、bandwidth 和 MTU 变化。

### L1：真实公网

使用不同云厂商/不同 AS 的 VPS，测试 IPv4、IPv6、高 RTT、长连接、100/1000 concurrent streams（并发流）和 100MB/1GB transfer（大数据传输）。

### L2：网络异常

记录 TCP RST、TLS alert、connection timeout、burst loss、path MTU change、NAT rebinding（若传输层支持）。

### L3：协议可观测性

检查非法 frame、超大 payload、错误版本、错误 WINDOW_UPDATE、未认证连接和 replay（重放）测试。

## PASS 标准

任何测试都必须保存原始日志、版本号、测试参数和结果摘要；不能因为“能连接”就标记 PASS。

## GFW 相关说明

不能用“是否被 GFW 封锁”作为本地自动化单元测试。真实网络测试只能记录客观现象，例如 TCP reset、TLS handshake failure、DNS failure、connectivity loss 和恢复时间；不得将结果解释为规避审查能力的保证。
