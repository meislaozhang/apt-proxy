# TCP Core Stability Gate v3

Purpose: close the remaining repository-local stability verification without external hosts, IP-specific targets, or Cloudflare Tunnel.

## Required checks

1. Five consecutive functional `go test ./...` runs.
2. Two consecutive `go test -race ./...` runs.
3. `go vet ./...`.
4. Build `apt-server` and `apt-client`.
5. Any failed iteration fails the workflow.
6. Each iteration has a bounded timeout.

This gate is a repeatability gate. It does not claim multi-hour soak testing or OS-level leak certification. Those remain separate acceptance evidence if required.