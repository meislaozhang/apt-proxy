# TCP Core Soak Gate

This gate validates repeatability of the existing local TCP core test suite without any external host, IP-specific target, or Cloudflare Tunnel.

## Acceptance

- 10 consecutive `go test ./...` runs must pass.
- 3 consecutive `go test -race ./...` runs must pass.
- `go vet ./...` must pass.
- `go build ./cmd/apt-server` must pass.
- `go build ./cmd/apt-client` must pass.
- Each test invocation has a bounded timeout so a deadlock cannot hang the workflow indefinitely.

This is a repeatability gate, not a substitute for a separately measured multi-hour soak or OS-level resource-leak test.
