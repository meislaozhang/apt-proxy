# Final Acceptance Checklist

## TCP Core

- [ ] 8/100/1000 stream E2E
- [ ] 1000-stream large integrity gate
- [ ] 1/16/64/100 MiB integrity gate
- [ ] optional 1 GiB integrity gate
- [ ] stream/session flow control
- [ ] zero-window backpressure
- [ ] HALF_CLOSE
- [ ] RESET
- [ ] race detector

## Remaining final gates

- [ ] netem latency/loss/jitter/reordering
- [ ] long-duration stability
- [ ] resource leak detection
- [ ] reconnect/recovery
- [ ] final CI evidence

100% is reported only after every required gate has concrete passing evidence.
