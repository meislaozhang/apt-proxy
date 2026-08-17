# 1000-stream Large Integrity Gate

This gate combines multiplexing concurrency with data-integrity verification.

- 1000 concurrent streams
- 64 KiB payload per stream
- 64 MiB total payload
- independent SHA-256 verification per stream
- 60 second per-stream read deadline
- all stream failures are collected so one failure does not hide sibling failures

The gate is an acceptance test for the TCP multiplexing/data-integrity layer. It does not by itself claim netem resilience, long-duration stability, leak-freedom, or release readiness.
