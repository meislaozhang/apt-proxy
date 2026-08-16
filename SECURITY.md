# Security Policy

APT is experimental software. Please do not assume it provides anonymity, censorship resistance, or production-grade security.

For vulnerabilities, use a private GitHub security report when available. Do not disclose exploitable details publicly before a fix is available.

## Current security boundary

APT/0.1 delegates confidentiality and server authentication to TLS 1.3. The application token is only an additional credential. The current MVP is intentionally not advertised as replay-resistant, independently audited, or production hardened.
