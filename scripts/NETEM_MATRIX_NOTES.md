# Netem robustness matrix

The matrix in `scripts/netem-matrix-v2.sh` covers baseline, fixed RTT, loss, jitter, reordering, and a 20 Mbit/s bandwidth cap.

Run on an isolated Linux test host with `CAP_NET_ADMIN`. The script changes the root qdisc and restores it on exit.

For each case, collect:

- APT E2E success/failure
- throughput and transfer duration
- SHA256 source/destination equality
- stream reset/recovery counts
- connection errors
- goroutine/FD/resource observations
