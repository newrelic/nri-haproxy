---
description: Use when comparing nri-haproxy metric output against OTel haproxyreceiver output for the spike analysis, or when capturing a baseline nri-haproxy metric snapshot.
---

Produce a structured metric comparison between nri-haproxy and the OTel haproxyreceiver.
The demo stack must be running (`./env-setup-local-container/start.sh` from the repo root).

## Phase check

First, determine which phase we are in:
- **Phase 1 only (nri-haproxy):** OTel collector is not yet in docker-compose. Capture the nri-haproxy baseline and document it fully. Skip the OTel and comparison sections.
- **Phase 2 (both running):** Capture both, produce the full side-by-side comparison.

Check by running: `docker compose -f env-setup-local-container/docker-compose.yml ps | grep otel`
If no output → Phase 1. If otel-collector is listed → Phase 2.

---

## Step 1 — Capture nri-haproxy snapshot

```bash
./bin/nri-haproxy \
  --stats_url http://localhost:8404/stats \
  --ha_proxy_cluster_name demo-cluster \
  2>/dev/null > /tmp/nri-snapshot.json
```

Parse it into three groups:

```bash
# Frontends
jq '[.data[] | select(.entity.type=="ha-frontend") | {entity: .entity.name, metrics: .metrics[0]}]' /tmp/nri-snapshot.json

# Backends
jq '[.data[] | select(.entity.type=="ha-backend") | {entity: .entity.name, metrics: .metrics[0]}]' /tmp/nri-snapshot.json

# Servers
jq '[.data[] | select(.entity.type=="ha-server") | {entity: .entity.name, metrics: .metrics[0]}]' /tmp/nri-snapshot.json
```

---

## Step 2 — Build the nri-haproxy metric inventory

For each entity type, extract every metric key and classify it.
Use this command to get all unique metric keys per entity type:

```bash
# Frontend metric keys
jq -r '[.data[] | select(.entity.type=="ha-frontend") | .metrics[0] | keys[]] | unique[]' /tmp/nri-snapshot.json

# Backend metric keys
jq -r '[.data[] | select(.entity.type=="ha-backend") | .metrics[0] | keys[]] | unique[]' /tmp/nri-snapshot.json

# Server metric keys
jq -r '[.data[] | select(.entity.type=="ha-server") | .metrics[0] | keys[]] | unique[]' /tmp/nri-snapshot.json
```

For each metric, note:
- **Name** (e.g. `frontend.httpRequestsPerSecond`)
- **Type** — GAUGE (point-in-time value), RATE (per-second computed rate), ATTRIBUTE (string/label)
- **Unit** — requests/s, bytes/s, seconds, count, string
- **Granularity** — which entity types report it (frontend / backend / server)
- **Source field** — the HAProxy CSV field it maps from (find in `src/definition.go`)

---

## Step 3 — (Phase 2 only) Capture OTel snapshot

```bash
# OTel collector exposes Prometheus metrics at :8888 (adjust port if different)
curl -s http://localhost:8888/metrics > /tmp/otel-snapshot.txt

# Filter to only haproxy metrics
grep "^haproxy_" /tmp/otel-snapshot.txt > /tmp/otel-haproxy-metrics.txt
```

For each OTel metric, note:
- **Name** (e.g. `haproxy_frontend_connections_total`)
- **Type** — counter, gauge, histogram
- **Unit** — from metric name suffix or HELP line
- **Labels** — what dimensions/attributes are available (proxy, server, state, etc.)
- **Granularity** — does it have per-frontend / per-backend / per-server breakdown via labels

---

## Step 4 — (Phase 2 only) Side-by-side comparison

Produce a markdown table with one row per logical metric concept:

| Concept | nri-haproxy metric | nri type | OTel metric | OTel type | Match? | Notes |
|---|---|---|---|---|---|---|
| Request rate | `frontend.httpRequestsPerSecond` | RATE | `haproxy_frontend_http_requests_total` | counter | ✅ name differs, same data | |
| Current sessions | `frontend.currentSessions` | GAUGE | `haproxy_frontend_current_sessions` | gauge | ✅ | |
| Response time | `backend.averageResponseTimeInSeconds` | GAUGE | ❌ missing | — | ❌ gap | OTel has no equivalent |

Flag each row as one of:
- ✅ **Match** — both collectors have it (possibly different names/units)
- ⚠️ **Partial** — both have it but different granularity, unit, or labeling
- ❌ **nri only** — nri-haproxy has it, OTel does not
- ❌ **OTel only** — OTel has it, nri-haproxy does not

---

## Step 5 — Summary findings

After the table, write a short summary covering:

1. **Total metric counts** — how many metrics each collector exposes
2. **Coverage gaps** — what nri-haproxy collects that OTel misses (these are candidates for OTel contributions)
3. **OTel advantages** — anything OTel provides that nri-haproxy doesn't
4. **Naming/unit differences** — where the same concept is reported differently
5. **Granularity differences** — e.g. nri has per-server health checks, does OTel?
6. **Customer migration impact** — what would a customer lose/gain switching from nri-haproxy to OTel

This summary maps directly to the Confluence output structure used in NR-548322 (Redis spike).
