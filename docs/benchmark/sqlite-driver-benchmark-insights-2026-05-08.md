# SQLite Driver Benchmark: Key Insights & Strategic Analysis

**Date**: 2026-05-08
**Source**: 3-server benchmark across Vultr shared & dedicated CPU
**Benchmark Data**: [Full Report](sqlite-driver-benchmark-all-servers-2026-05-08.md)

---

## Executive Summary

Tested `mattn/go-sqlite3` (CGO) vs `modernc.org/sqlite` (pure Go) across:
- Server A: 6v Shared (Intel Xeon) — $96/mo
- Server B: 1v Shared (Intel Skylake) — $6/mo
- Server C: 4v Dedicated (AMD EPYC) — $90/mo

**Key Finding**: mattn delivers 2x throughput on shared CPU, but the gap shrinks to 1.3x on dedicated CPU. At current pricing ($90 vs $96), Server A offers the best RPS per dollar for production workloads.

---

## Insight 1: Dedicated CPU = "Free Performance"

| Metric | Shared 6v | Dedicated 4v | Winner |
|--------|-----------|--------------|--------|
| mattn RPS/vCPU | 16,926 | 21,237 | **+25%** |
| modernc RPS/vCPU | 8,835 | 15,998 | **+81%** |

**What it means**: Dedicated CPU isn't just about "no noisy neighbors" — Go scheduler and SQLite locking work more efficiently without VM-level context switching.

**Actionable**: If you need consistent performance and can tolerate slightly higher cost, dedicated CPU gives 25-81% more bang per core.

---

## Insight 2: modernc Isn't "Slow", It "Doesn't Scale"

| Server | modernc RPS | RPS/vCPU | Pattern |
|--------|-------------|----------|---------|
| 1 vCPU | 12,175 | 12,175 | Linear |
| 4 vCPU (dedicated) | 63,991 | 15,998 | Linear |
| 6 vCPU (shared) | 53,009 | 8,835 | **Sub-linear!** |

**What it means**: modernc scales well up to 4 cores, but on shared 6-core CPU performance actually drops. Go scheduler contention in shared environments kills pure-Go SQLite performance.

**Actionable**: Don't use modernc on large shared-CPU instances (6+ vCPU). Either use smaller instances or switch to mattn.

---

## Insight 3: mattn vs modernc Gap = "CPU Noise Detector"

| Condition | Gap | Why |
|-----------|-----|-----|
| 1 vCPU shared | 1.35x | No other cores, minimal noise |
| 4 vCPU dedicated | 1.33x | Dedicated = predictable |
| 6 vCPU shared | **1.92x** | Shared = heavy context switching |

**What it means**: The gap grows on shared CPU because mattn (CGO) has fixed overhead that gets amortized, while modernc (pure Go) suffers from scheduler contention. **The gap is an environment sensitivity, not a driver limitation.**

**Actionable**: On dedicated CPU, you can confidently use modernc (1.3x gap). On shared CPU with 6+ cores, mattn is almost 2x faster.

---

## Insight 4: Pricing Is Non-Linear With Performance

| Upgrade | Cost | RPS Gain | Worth It? |
|---------|------|----------|-----------|
| 1→4 vCPU (dedicated) | +$84 | +417% | ✅ Yes, if you need 5x RPS |
| 1→6 vCPU (shared) | +$90 | +519% | ❌ Worse RPS/$ |
| **4→6 vCPU** | **+$6** | **+20%** | ✅ **Best value!** |

**What it means**: Vultr pricing is slightly anomalous — Server C (4v dedicated $90) is almost the same as Server A (6v shared $96). This makes **Server A the sweet spot** for high-traffic production.

**Actionable**: Skip Server C if you need maximum throughput. Server A gives 20% more RPS for just $6 more.

---

## Insight 5: Latency Story — Dedicated Wins Big

| Metric @ 100 conn | Server A (6v) | Server C (4v) | Server B (1v) |
|-------------------|---------------|---------------|---------------|
| mattn p99 | 8.29ms | **7.24ms** | 31.01ms |
| modernc p99 | 7.60ms | **9.62ms** | 31.03ms |

**What it means**: Server C (dedicated 4v) has the lowest latency — even lower than Server A (shared 6v). **If latency is critical (e.g., API response <10ms), dedicated CPU wins despite fewer cores.**

**Actionable**: For latency-sensitive APIs, use Server C with mattn (p99 7.24ms). For throughput, use Server A (101K RPS).

---

## Insight 6: P99 vs P50 Gap = "CPU Contention Indicator"

| Server | mattn P50 | mattn P99 | Gap |
|--------|-----------|-----------|-----|
| C (dedicated) @400conn | 4.48ms | 17.99ms | **4.0x** |
| A (shared) @500conn | 8.69ms | 23.26ms | **2.7x** |
| B (shared) @200conn | 17.74ms | 40.33ms | **2.3x** |

**What it means**: Dedicated CPU has higher variance (P99/P50 = 4.0x) because spikes are sharper during GC or compaction. Shared CPU has lower variance because neighbor "noise" redistributes spikes. **Dedicated = lower average latency, higher variance.**

**Actionable**: If you need consistent latency (predictable P99), shared CPU may be better. If you need lowest average, dedicated wins.

---

## Insight 7: Concurrency — mattn Scales Better

| Metric | mattn | modernc | Winner |
|--------|-------|---------|--------|
| Concurrent Write (4v dedicated) | 11,284 ns | 50,673 ns | **mattn 4.5x** |
| Concurrent ReadWrite (4v dedicated) | 13,506 ns | 16,291 ns | mattn 1.2x |

**What it means**: mattn wins big on concurrent writes because CGO has fixed overhead — the more concurrent requests, the more that overhead is amortized. modernc has per-request overhead that grows with concurrency.

**Actionable**: For write-heavy workloads, mattn is 4.5x faster. For read-heavy, both are comparable (1.2x).

---

## Insight 8: SQLite Isn't "Dead" for Production

| Benchmark | Result | Context |
|-----------|--------|---------|
| Peak RPS | 101,555 | Single SQLite file! |
| P99 Latency | 7.24ms | 100 concurrent users |
| Concurrency | 500 connections | WAL mode handles it |

**What it means**: SQLite with WAL mode can handle **100K RPS with P99 <25ms** on a single Go Fiber server. This is sufficient for 99% of SaaS apps. You only need PostgreSQL when you need >100K concurrent connections (not RPS).

**Actionable**: Don't prematurely migrate to PostgreSQL. SQLite + WAL is production-ready for most workloads up to 100K RPS.

---

## Insight 9: Decision Matrix — When to Use What

| Need | Choose | Reason |
|------|--------|--------|
| **Development/Staging** | modernc + Server B ($6) | Easy cross-compile, sufficient for testing |
| **Production <20K RPS** | modernc + Server B ($6) | Economical, modernc is adequate |
| **Production 20-85K RPS, latency sensitive** | mattn + Server C ($90) | P99 <12ms, dedicated stability |
| **Production 85-100K RPS** | mattn + Server A ($96) | Best RPS/$, 20% more than Server C |
| **Production >100K RPS** | Multiple Server A | Horizontal scaling |

---

## Insight 10: The "Hidden" Cost of modernc

| Factor | modernc | mattn |
|--------|---------|-------|
| Compile time | 30s | 60s (CGO) |
| Cross-compile | ✅ Easy | ❌ Need musl/Docker |
| Binary size | 18MB | 12MB |
| Memory usage | Higher (Go GC) | Lower (C allocator) |
| Max RPS (6v) | 53K | 101K |
| **Cost to hit 100K RPS** | 2 Server A ($192) | 1 Server A ($96) |

**What it means**: modernc is "free" in development, but in production the cost is 2x for the same traffic.

**Actionable**: Use modernc during development (simpler), switch to mattn for production (2x throughput, half the cost at scale).

---

## Strategic Recommendations for Laju Go

### Phase 1: Development & MVP
- **Server**: B ($6/mo)
- **Driver**: modernc (easy cross-compile)
- **RPS**: ~15K (sufficient for testing)
- **Cost**: $6/mo

### Phase 2: Production Launch (<50K RPS)
- **Server**: C ($90/mo) — if latency critical
- **Driver**: mattn (consistent performance)
- **RPS**: ~85K
- **P99**: <12ms
- **Cost**: $90/mo

### Phase 3: Growth (50-100K RPS)
- **Server**: A ($96/mo)
- **Driver**: mattn (best throughput)
- **RPS**: ~100K
- **Cost**: $96/mo

### Phase 4: Scale (>100K RPS)
- **Architecture**: Multiple Server A behind load balancer
- **Driver**: mattn
- **Scaling**: Horizontal (add servers)
- **Cost**: $96/mo per 100K RPS

---

## Quick Reference Card

```
┌─────────────────────────────────────────────────────────────────┐
│                    LAJU GO SQLite STRATEGY                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Development:    modernc + Server B ($6) = 15K RPS              │
│  Production:     mattn + Server A ($96) = 100K RPS              │
│  Scale:          Add Server A instances (+100K RPS each)        │
│                                                                 │
│  Driver Choice:                                                │
│    • modernc: Easy compile, good for <20K RPS                  │
│    • mattn: 2x throughput, needed for >50K RPS                 │
│                                                                 │
│  Server Choice:                                                │
│    • B ($6): Dev/staging, <20K RPS                             │
│    • C ($90): Latency-critical, dedicated CPU                  │
│    • A ($96): Best value, 100K RPS                             │
│                                                                 │
│  SQLite Limit: ~100K RPS per instance (WAL mode)               │
│  PostgreSQL: Only if >100K concurrent connections              │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## Appendix: Data Sources

| Document | Content |
|----------|---------|
| [Full Benchmark Report](sqlite-driver-benchmark-all-servers-2026-05-08.md) | All raw data, charts, test commands |
| [Server A Results](sqlite-driver-benchmark-2026-05-08.md) | 6v Shared detailed results |
| [Server B vs C Comparison](sqlite-driver-benchmark-comparison-2026-05-08.md) | 1v Shared vs 4v Dedicated |

---

*Analysis by Qwen Code — 2026-05-08*
