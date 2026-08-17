#!/usr/bin/env bash
set -euo pipefail

# Reproducible APT transport robustness matrix.
# Run on the isolated test host with CAP_NET_ADMIN.
# TARGET_HOST is the APT E2E endpoint; default is the supplied test host.
# APT_E2E_CMD must be provided. It is executed once per impairment case.
# The command must return non-zero on any data-integrity/recovery failure.

IFACE="${IFACE:-eth0}"
TARGET_HOST="${TARGET_HOST:-192.236.230.171}"
CASE_SLEEP="${CASE_SLEEP:-1}"
LOG_DIR="${LOG_DIR:-./artifacts/netem}"
APT_E2E_CMD="${APT_E2E_CMD:-}"

if [[ -z "$APT_E2E_CMD" ]]; then
  echo "ERROR: APT_E2E_CMD is required; refusing to report a synthetic PASS." >&2
  exit 2
fi

mkdir -p "$LOG_DIR"
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

run_case() {
  local name="$1" delay="$2" loss="$3" jitter="$4" reorder="$5" bandwidth="$6"
  local log="$LOG_DIR/${name}.log"

  echo "=== $name ==="
  cleanup

  if [[ "$bandwidth" != "0" ]]; then
    tc qdisc add dev "$IFACE" root handle 1: tbf rate "$bandwidth" burst 32kbit latency 400ms
    local netem=(netem delay "$delay")
    [[ "$jitter" != "0ms" ]] && netem+=(10ms)
    [[ "$loss" != "0%" ]] && netem+=(loss "$loss")
    [[ "$reorder" != "0%" ]] && netem+=(reorder "$reorder" 50%)
    tc qdisc add dev "$IFACE" parent 1:1 handle 10: "${netem[@]}"
  else
    local netem=(netem delay "$delay")
    [[ "$jitter" != "0ms" ]] && netem+=(10ms)
    [[ "$loss" != "0%" ]] && netem+=(loss "$loss")
    [[ "$reorder" != "0%" ]] && netem+=(reorder "$reorder" 50%)
    tc qdisc add dev "$IFACE" root "${netem[@]}"
  fi

  {
    echo "case=$name"
    echo "target=$TARGET_HOST"
    echo "delay=$delay loss=$loss jitter=$jitter reorder=$reorder bandwidth=$bandwidth"
    echo "started=$(date -Is)"
    echo "command=$APT_E2E_CMD"
  } | tee "$log"

  if TARGET_HOST="$TARGET_HOST" APT_NETEM_CASE="$name" bash -lc "$APT_E2E_CMD" >>"$log" 2>&1; then
    echo "result=PASS" | tee -a "$log"
  else
    echo "result=FAIL" | tee -a "$log"
    return 1
  fi

  echo "finished=$(date -Is)" | tee -a "$log"
  sleep "$CASE_SLEEP"
}

failures=0
for spec in "${CASES[@]}"; do
  IFS='|' read -r name delay loss jitter reorder bandwidth <<<"$spec"
  if ! run_case "$name" "$delay" "$loss" "$jitter" "$reorder" "$bandwidth"; then
    failures=$((failures + 1))
  fi
done

cleanup
if (( failures > 0 )); then
  echo "NETEM_MATRIX_FAIL failures=$failures" >&2
  exit 1
fi

echo "NETEM_MATRIX_PASS cases=${#CASES[@]} target=$TARGET_HOST"
