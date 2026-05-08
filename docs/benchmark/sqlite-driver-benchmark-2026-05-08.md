# SQLite Driver Benchmark: mattn/go-sqlite3 vs modernc.org/sqlite

**Date**: 2026-05-08
**Tester**: Qwen Code
**Purpose**: Evaluate SQLite driver performance for Laju Go production deployment

---

## Server Specifications

| Resource | Specification |
|----------|---------------|
| Provider | Vultr |
| CPU | Intel Xeon Processor (Cascadelake) @ 2.0GHz |
| Cores | 3 cores, 2 threads/core = **6 vCPUs** |
| RAM | 15 GB |
| Storage | SSD |
| OS | Ubuntu 24.04 (Kernel 6.8.0-111-generic) |
| Go Version | 1.24.3 |

---

## Part 1: Microbenchmark (go test -bench)

**Methodology**: `go test -bench=. -benchmem -count=5 -timeout=10m`
**Database**: In-memory & file-based SQLite with WAL mode
**Date**: 2026-05-08 09:30 UTC

### Raw Results (Average of 5 runs)

| Operation | mattn (ns/op) | modernc (ns/op) | mattn Faster | mattn Allocs | modernc Allocs |
|-----------|--------------:|----------------:|:------------:|:------------:|:--------------:|
| Insert | 9,300 | 47,200 | **5.1x** | 17 | 12 |
| Bulk Insert (100 rows) | 721,000 | 1,224,000 | **1.7x** | 1,719 | 1,207 |
| Select | 11,400 | 14,100 | **1.2x** | 38 | 36 |
| Select All (1000 rows) | 1,969,000 | 2,050,000 | **1.04x** | 9,669 | 15,663 |
| Update | 5,040 | 7,630 | **1.5x** | 11 | 6 |
| Delete | 22,340 | 95,910 | **4.3x** | 39 | 30 |
| Prepared Insert | 5,210 | 48,400 | **9.3x** | 14 | 12 |
| Insert (file-based) | 33,970 | 43,320 | **1.3x** | 17 | 12 |
| Concurrent Read | 14,680 | 17,210 | **1.2x** | 34 | 32 |
| Concurrent Write | 19,380 | 63,600 | **3.3x** | 17 | 12 |
| Concurrent ReadWrite | 20,390 | 28,400 | **1.4x** | 31 | 28 |
| HTTP GET (file DB) | 12,800 | 18,420 | **1.4x** | 71 | 69 |
| HTTP GET + Middleware | 12,270 | 20,670 | **1.7x** | 74 | 72 |

### Microbenchmark Analysis

- **Biggest gap**: Prepared statements (9.3x) and single inserts (5.1x)
- **Smallest gap**: Large SELECT queries (1.04x) — reading many rows narrows the difference
- **Memory**: modernc uses fewer allocations in most operations (better GC pressure)
- **Write operations**: mattn significantly faster (3-9x) due to CGO direct kernel access

---

## Part 2: HTTP Load Test (wrk)

**Methodology**: `wrk` with Fiber v2 HTTP server, SQLite file-based DB with WAL mode
**Test Duration**: 30 seconds per test
**Database**: 50,000 pre-populated users
**Endpoint**: `GET /users/:id` (prepared SELECT statement)

### Test 1: 100 Concurrent Connections

```
wrk -t6 -c100 -d30s --latency http://localhost:PORT/users/1
```

| Metric | mattn | modernc | Difference |
|--------|------:|--------:|:----------:|
| **Requests/sec** | **99,469** | 53,653 | 1.85x |
| Avg Latency | 1.43ms | 2.01ms | 1.41x |
| p50 Latency | 762µs | 1.69ms | 2.22x |
| p75 Latency | 1.91ms | 2.85ms | 1.49x |
| p90 Latency | 3.70ms | 4.22ms | 1.14x |
| p99 Latency | 8.29ms | 7.60ms | ~same |
| Max Latency | 26.79ms | 52.03ms | 1.94x |
| Total Requests | 2,987,069 | 1,610,662 | 1.85x |
| Transfer/sec | 15.56 MB | 8.39 MB | 1.85x |

### Test 2: 500 Concurrent Connections (Full Saturation)

```
wrk -t6 -c500 -d30s --latency http://localhost:PORT/users/1
```

| Metric | mattn | modernc | Difference |
|--------|------:|--------:|:----------:|
| **Requests/sec** | **101,555** | 53,009 | **1.92x** |
| Avg Latency | 5.50ms | 10.28ms | 1.87x |
| p50 Latency | 4.60ms | 8.63ms | 1.88x |
| p75 Latency | 7.51ms | 13.89ms | 1.85x |
| p90 Latency | 11.26ms | 20.48ms | 1.82x |
| p99 Latency | 23.26ms | 37.13ms | 1.60x |
| Max Latency | 60.86ms | 116.10ms | 1.91x |
| Total Requests | 3,054,709 | 1,592,966 | 1.92x |
| Transfer/sec | 15.88 MB | 8.29 MB | 1.92x |

### Test 3: 12 Threads / 500 Connections (Over-subscription)

```
wrk -t12 -c500 -d15s --latency http://localhost:PORT/users/1
```

| Metric | mattn | modernc | Difference |
|--------|------:|--------:|:----------:|
| **Requests/sec** | **84,500** | 50,010 | 1.69x |
| Avg Latency | 6.85ms | 10.75ms | 1.57x |
| p99 Latency | 30.91ms | 38.92ms | 1.26x |

**Note**: RPS dropped for both drivers due to context switching overhead (12 threads > 6 vCPUs).

---

## Part 3: CPU Utilization (500 Connections)

| Metric | mattn | modernc |
|--------|------:|--------:|
| User CPU | 47.8% | 53.5% |
| System CPU | 34.3% | 22.5% |
| **Total CPU** | **~82%** | **~76%** |
| Idle | 6.0% | 4.2% |
| **RPS per 1% CPU** | **1,238** | **697** |

**mattn delivers 1.78x more RPS per CPU percentage** — more efficient resource usage.

---

## Part 4: Scaling Projection

### Current Server (3 cores, 6 vCPUs)

| Driver | RPS | Latency (avg) | Latency (p99) |
|--------|----:|--------------:|--------------:|
| mattn | ~100K | 5.50ms | 23.26ms |
| modernc | ~53K | 10.28ms | 37.13ms |

### Projected with 6 cores (12 vCPUs)

| Driver | Est. RPS | Est. Latency (avg) | Est. Latency (p99) |
|--------|--------:|-------------------:|-------------------:|
| mattn | ~200K | ~5ms | ~22ms |
| modernc | ~106K | ~10ms | ~35ms |

### Scaling Formula

```
RPS = N_vCPU × base_RPS_per_vCPU

mattn:    RPS = N_vCPU × 16,926  (101,555 / 6 vCPU)
modernc:  RPS = N_vCPU × 8,835   (53,009 / 6 vCPU)

Rasio tetap: ~1.9x regardless of CPU count
```

---

## Part 5: Recommendations

### For Laju Go Project

| Use Case | Recommendation | Reason |
|----------|----------------|--------|
| **Auth/session (read-heavy)** | mattn | 1.85x faster reads, lower latency |
| **Dashboard/listing** | mattn | Better p99 latency (23ms vs 37ms) |
| **Form submissions** | mattn | 3-5x faster writes in microbenchmark |
| **File uploads** | mattn | 1.3x faster file-based inserts |
| **Early stage (<1000 users)** | Either works | modernc simpler cross-compile |
| **Production (>1000 users)** | **mattn** | 2x throughput headroom |

### Decision Matrix

| Factor | mattn | modernc | Winner |
|--------|-------|---------|:------:|
| Raw performance | ~100K RPS | ~53K RPS | mattn |
| Latency (p99) | 23ms | 37ms | mattn |
| CPU efficiency | 1,238 RPS/%CPU | 697 RPS/%CPU | mattn |
| Memory allocations | More | Less | modernc |
| Cross-compile | Needs CGO | Pure Go | modernc |
| Build complexity | Higher | Lower | modernc |

### Final Verdict

**Switch to `mattn/go-sqlite3` for production deployment.**

**Why:**
1. **2x throughput** under full load — significant for scaling
2. **40% better p99 latency** — better user experience
3. **78% more efficient** per CPU percentage — better cost efficiency
4. **Build complexity trade-off is acceptable** — can build on Linux server or use Docker

**When to stay with modernc:**
- Strict requirement for cross-compilation from macOS/Windows
- Very early prototype with <100 users
- Development environment only

---

## Appendix: Test Commands

### Microbenchmark
```bash
ssh root@141.164.44.250 "cd /opt/go-sqlite-benchmark-mattn-vs-modernc && go test -bench=. -benchmem -count=5 -timeout=10m ./..."
```

### HTTP Load Test (mattn)
```bash
ssh root@141.164.44.250 "/tmp/benchmark_mattn -driver mattn -port 3001 -db /tmp/mattn_bench.db &"
wrk -t6 -c500 -d30s --latency http://141.164.44.250:3001/users/1
```

### HTTP Load Test (modernc)
```bash
ssh root@141.164.44.250 "/tmp/benchmark_mattn -driver modernc -port 3002 -db /tmp/modernc_bench.db &"
wrk -t6 -c500 -d30s --latency http://141.164.44.250:3002/users/1
```

---

*Report generated by Qwen Code — 2026-05-08*
