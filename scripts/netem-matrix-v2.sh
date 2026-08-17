#!/usr/bin/env bash
set -euo pipefail

# Reproducible APT transport robustness matrix.
# Run on an isolated Linux test host with CAP_NET_ADMIN.
# This only simulates deterministic network impairments; it does not
# bypass filtering or alter routing policy.

IFACE="${IFACE:-eth0}"
CASE_SLEEP="${CASE_SLEEP:-1}"
CASES=(
  "baseline|0ms|0%|0ms|0%|0"
  "rtt80|80ms|0%|0ms|0%|0"
  "rtt150-loss1|150ms|1%|0ms|0%|0"
  "rtt250-jitter|250ms|1%|10ms|0%|0"
  "reorder3|80ms|0%|0ms|3%|0"
  "loss3|80ms|3%|0ms|0%|0"
  "bandwidth20|80ms|0%|0ms|0%|20mbit"
)

cleanup() { tc qdisc del dev "$IFACE" root 2>/dev/null || true; }
trap cleanup EXIT

for spec in "${CASES[@]}"; do
  IFS='|' read -r name delay loss jitter reorder bandwidth <<<"$spec"
  echo "=== $name ==="
  cleanup

  if [[ "$bandwidth" != "0" ]]; then
    tc qdisc add dev "$IFACE" root handle 1: tbf rate "$bandwidth" burst 32kbit latency 400ms
    netem=(netem delay "$delay")
    [[ "$jitter" != "0ms" ]] && netem+=(10ms)
    [[ "$loss" != "0%" ]] && netem+=(loss "$loss")
    [[ "$reorder" != "0%" ]] && netem+=(reorder "$reorder" 50%)
    tc qdisc add dev "$IFACE" parent 1:1 handle 10: "${netem[@]}"
  else
    netem=(root netem delay "$delay")
    [[ "$jitter" != "0ms" ]] && netem+=(10ms)
    [[ "$loss" != "0%" ]] && netem+=(loss "$loss")
    [[ "$reorder" != "0%" ]] && netem+=(reorder "$reorder" 50%)
    tc qdisc add dev "$IFACE" "${netem[@]:1}"
  fi

  echo "RUN: $name"
  echo "Collect APT E2E logs, throughput, SHA256, stream recovery and connection metrics here."
  sleep "$CASE_SLEEP"
done
