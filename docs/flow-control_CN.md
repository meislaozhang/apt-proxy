# APT 流量控制（Flow Control）设计

APT 使用显式 `WINDOW_UPDATE`（窗口更新）帧逐步增加对端可发送的字节额度。

## 设计目标

- 防止一个慢接收端耗尽服务端内存。
- 每条 Stream（数据流）拥有独立窗口。
- 后续可扩展到 Session（会话）级总窗口。
- 窗口增量为非零 32-bit 无符号整数。

## 帧语义

`WINDOW_UPDATE` 的 Payload（负载）为 4 字节大端整数，表示新增可发送字节数。

接收方处理后增加对应 Stream 的发送额度；发送方在额度耗尽时必须产生 Backpressure（反压），而不是无限缓存 DATA（数据）帧。

## 后续实现

APT/0.2 将把该帧接入 Stream 状态机，并加入：

1. 初始 Stream Window。
2. 消费数据后的窗口回收。
3. Session Window。
4. 最大窗口限制与溢出检查。
5. 超时与关闭策略。
6. 并发压力测试与 race detector（数据竞争检测）。
