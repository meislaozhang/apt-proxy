#!/usr/bin/env bash
set -euo pipefail

# Reproducible APT transport robustness matrix.
# Run on an isolated Linux test host with CAP_NET_ADMIN.
# This does NOT attempt to bypass network filtering; it only simulates
# deterministic delay/loss/reordering/bandwidth conditions.

IFACE="${IFACE:-eth0}"
CASES=(
  "baseline|0ms|0%|0%|0"
  "rtt80|80ms|0%|0%|0"
  "rtt150-loss1|150ms|1%|0%|0"
  "rtt250-jitter|250ms|1%|10ms|0"
  "reorder3|80ms|0%|0%|3%"
  "loss3|80ms|3%|0%|0"
  "bandwidth20|80ms|0%|0%|0|20mbit"
)

cleanup() { tc qdisc del dev "$IFACE" root 2>/dev/null || true; }
trap cleanup EXIT

for spec in "${CASES[@]}"; do
  IFS='|' read -r name delay loss _ jitter reorder bandwidth <<<"$spec"
  echo "=== $name ==="
  cleanup
  args=(root netem delay "$delay")
  [[ "$loss" != "0%" ]] && args+=(loss "$loss")
  [[ "$jitter" != "0ms" ]] && args+=(10ms)
  [[ "$reorder" != "0%" ]] && args+=(reorder "$reorder" 50%)
  if [[ "${bandwidth:-0}" != "0" ]]; then
    tc qdisc add dev "$IFACE" root handle 1: tbf rate "$bandwidth" burst 32kbit latency 400ms
    tc qdisc add dev "$IFACE" parent 1:1 handle 10: netem delay "$delay" loss "$loss"
  else
    tc qdisc add dev "$IFACE" "${args[@]:1}"
  fi
  echo "RUN: $name"
  echo "Collect APT E2E logs, throughput, SHA256 and connection/recovery metrics here."
  sleep "${CASE_SLEEP:-1}"
done
