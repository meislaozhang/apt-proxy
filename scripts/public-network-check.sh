#!/usr/bin/env bash
set -euo pipefail

# Public-network reliability probe. This measures connectivity and TLS/HTTP
# behavior; it does not attempt to bypass filtering or evade inspection.
TARGET=${1:?usage: $0 host:port [server_name]}
SERVER_NAME=${2:-${TARGET%%:*}}
HOST=${TARGET%:*}
PORT=${TARGET##*:}
TIMEOUT=${TIMEOUT:-10}

printf 'target=%s\nserver_name=%s\n' "$TARGET" "$SERVER_NAME"
printf '\n== DNS ==\n'
getent ahosts "$HOST" || true
printf '\n== TCP ==\n'
if timeout "$TIMEOUT" bash -c "</dev/tcp/$HOST/$PORT" 2>/dev/null; then
  echo 'tcp_connect=PASS'
else
  echo 'tcp_connect=FAIL'
fi
printf '\n== TLS ==\n'
if command -v openssl >/dev/null 2>&1; then
  if timeout "$TIMEOUT" openssl s_client -connect "$HOST:$PORT" -servername "$SERVER_NAME" -brief </dev/null 2>&1 | grep -q 'CONNECTION'; then
    echo 'tls_handshake=PASS'
  else
    echo 'tls_handshake=FAIL'
  fi
else
  echo 'tls_handshake=SKIP openssl-not-installed'
fi
printf '\n== Timing ==\n'
if command -v curl >/dev/null 2>&1; then
  curl -sk --connect-timeout "$TIMEOUT" -o /dev/null -w 'connect=%{time_connect} tls=%{time_appconnect} total=%{time_total}\n' "https://$HOST:$PORT/" || true
else
  echo 'curl=SKIP'
fi
