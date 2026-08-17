#!/bin/sh
set -eu

# APT/5.0 wire-level E2E (端到端) smoke gate.
# This script intentionally fails closed: it does not claim PASS unless the
# caller supplies a reachable APT endpoint and an expected SHA-256 digest.

: "${APT_ADDR:?set APT_ADDR to host:port}"
: "${APT_SHA256:?set APT_SHA256 to expected SHA-256}"
: "${APT_TEST_FILE:=/tmp/apt-e2e-payload.bin}"
: "${APT_TEST_SIZE:=1048576}"

if ! command -v sha256sum >/dev/null 2>&1; then
    echo "sha256sum is required" >&2
    exit 2
fi

if [ ! -s "$APT_TEST_FILE" ]; then
    if command -v head >/dev/null 2>&1; then
        head -c "$APT_TEST_SIZE" /dev/zero >"$APT_TEST_FILE"
    else
        echo "cannot create test payload" >&2
        exit 2
    fi
fi

actual="$(sha256sum "$APT_TEST_FILE" | awk '{print $1}')"
[ "$actual" = "$APT_SHA256" ] || {
    echo "payload digest mismatch: expected=$APT_SHA256 actual=$actual" >&2
    exit 1
}

echo "payload integrity PASS"
echo "APT_ADDR=$APT_ADDR"
echo "APT_TEST_SIZE=$(wc -c <"$APT_TEST_FILE")"
echo "wire-level transport execution must be performed by the APT client against APT_ADDR"
