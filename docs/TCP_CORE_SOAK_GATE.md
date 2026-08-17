# TCP Core Soak Gate

This gate is the final repository-local stability check before release acceptance.

## Scope

- No external host or IP target.
- No Cloudflare Tunnel.
- Ten repeated functional test passes.
- Three repeated `-race` test passes.
- Final `go vet` and both binary builds.

## Acceptance

The gate passes only when every iteration exits successfully. A timeout or race failure fails the workflow immediately.

This is intentionally bounded for CI. It is evidence of repeatability, not a claim of an unbounded production soak test or OS-level leak proof.
