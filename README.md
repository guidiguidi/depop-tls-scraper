# Depop Real-Time Scraper & Monitor (Go)

A high-throughput, low-latency scraper and monitoring engine for Depop, built with Go.

## Key Highlights & Engineering Features

- **Cloudflare WAF Bypass without Headless Browsers:** Leverages `bogdanfinn/tls-client` and `fhttp` to spoof browser TLS Client Hello signatures (JA3/JA4 fingerprints, cipher suites, and HTTP/2 settings matching Chrome 120).
- **10x Resource Efficiency:** Operates at raw HTTP speeds with minimal RAM footprint (~20 MB vs ~500 MB per browser instance).
- **Concurrency & State Management:** Concurrent query polling with thread-safe in-memory cache (`alerter/cache.go`) for instant deduplication of seen listings.
- **Real-Time Notification Pipeline:** Formatted Telegram alerts with item pricing, images, and deep links.

## Architecture

```
Target Queries -> TLS Client (JA3/JA4 Spoofing) -> Depop API -> Deduplication Cache -> Telegram Notifier
```

## Configuration

Copy `.env.example` to `.env`:

```bash
cp .env.example .env
```

Parameters:
- `DEPOP_QUERIES`: Comma-separated search queries (e.g., `rick owens, balenciaga`).
- `DEPOP_INTERVAL`: Polling interval in seconds.
- `TG_TOKEN`: Telegram bot token.
- `TG_CHAT_ID`: Telegram chat ID for alerts.
- `SCRAPEOPS_API_KEY`: (Optional) Proxy / header aggregator key.

## Build & Run

```bash
go build -o depop-monitor .
./depop-monitor
```
