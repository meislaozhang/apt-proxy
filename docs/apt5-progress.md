# APT/5.0 Progress Gate

## Current target: 75%

Required gates before claiming 75%:

- [x] Bidirectional stream/session flow-control implementation
- [x] Client remote half-close wake-up implementation
- [x] Server origin half-close implementation
- [ ] Half-close lifecycle regression test passing in CI
- [ ] Bidirectional half-close E2E
- [ ] RST and close-race E2E
- [ ] 100 concurrent streams E2E
- [ ] Large-transfer integrity E2E
- [ ] Real VPS/Cloudflare Tunnel E2E
- [ ] UDP/QUIC/HTTP3 milestones

A percentage is only advanced when the corresponding CI or E2E evidence exists; code-only implementation is not counted as a passing gate.
