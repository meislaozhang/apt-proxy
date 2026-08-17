# Final Repository Stability Gate

This is the final repository-local stability gate before acceptance.

## Checks

- 3 consecutive functional `go test ./...` runs.
- 2 consecutive `go test -race ./...` runs.
- `go vet ./...`.
- Build `apt-server`.
- Build `apt-client`.
- Emit `FINAL_REPOSITORY_STABILITY_GATE_PASS` only after every check succeeds.

No external host, IP-specific target, or Cloudflare Tunnel is used.

This gate is deliberately fast and bounded. Long-duration soak and OS-level resource-leak certification remain separate evidence and are not falsely represented as complete by this workflow.