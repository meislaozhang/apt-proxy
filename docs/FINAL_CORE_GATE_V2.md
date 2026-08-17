# Final TCP Core Gate v2

This gate is intentionally self-contained and does not depend on any external test host.

## Required checks

1. `go test -timeout 3m ./...`
2. `go test -timeout 6m -race ./...`
3. `go vet ./...`
4. `go build ./cmd/apt-server`
5. `go build ./cmd/apt-client`

## Acceptance boundary

The gate proves the repository builds cleanly, the complete test suite passes, and the race detector is clean within bounded execution time.

It does **not** claim that long-duration soak, operating-system resource leak testing, or arbitrary network impairment has been completed. Those remain separate evidence gates.

No external IP address and no Cloudflare Tunnel are part of this gate.
