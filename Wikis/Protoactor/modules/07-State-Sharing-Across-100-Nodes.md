---
tags: [protoactor, golang, distributed-state, gossip, grains, architecture]
module: 7
---

# Module 7 — State Sharing Across 100s of Nodes

> [!info] Goal
> Evaluate when to use gossip vs. grain actors for state sharing at scale, and design a production-ready architecture.

**← Prev** [[06-Gossip-Protocol]]

---

## The Core Question

> Can Proto.Actor's gossip protocol be used for state sharing across 100s of nodes?

**Short answer**: Yes — but only for the *right kind* of state. Understanding *which pattern to use* is the key architectural decision.

---

## Two Patterns Side by Side

### Pattern A: Gossip (Metadata Propagation)

```
Use for: "What does everyone need to know?"
─────────────────────────────────────────────

         ┌────────┐
  ┌──────► Node 1 ◄──────┐
  │      └───┬────┘      │
gossip       │       gossip
  │       broadcast      │
  │          │           │
┌─┴──────┐   │   ┌───────┴─┐
│ Node 2 ◄───┘   │ Node 3  │
└────────┘       └─────────┘

Every node maintains its own copy.
Reads are local (no network hop).
Writes propagate in O(log N) rounds.
Consistency: eventual.
```

### Pattern B: Grain Actors (Authoritative Ownership)

```
Use for: "Who owns this piece of state?"
──────────────────────────────────────────

  Any Node     Cluster        Owning Node
     │            │                │
     │─"user/42"─►│                │
     │            │──────────────►grain "user/42"
     │            │               │  owns + mutates state
     │◄───────────┼───────────────│  responds with result
     │  response  │               │

One canonical copy. Strong sequential consistency per grain.
Reads/writes route to the grain (may cross network).
```

---

## Fit Analysis

| Use Case | Gossip | Grain Actors | Verdict |
|---|---|---|---|
| Cluster membership & health | ✅ Perfect | ❌ Overkill | **Gossip** |
| Feature flags / config | ✅ Good (eventual) | ✅ Works | **Gossip** (simpler) |
| Per-entity state (user, order) | ⚠️ Wrong tool | ✅ Perfect | **Grain** |
| Global counter | ⚠️ Races | ✅ One grain owns it | **Grain** |
| Leaderboard / rankings | ❌ Stale reads | ✅ Grain aggregates | **Grain** |
| Bulk data / large payloads | ❌ Not designed for | ✅ Stream to grain | **Grain** |
| Strong consistency required | ❌ Eventually consistent | ✅ Serialized by actor | **Grain** |
| Read-heavy, stale-OK metadata | ✅ Local reads, fast | ⚠️ Every read = network | **Gossip** |

---

## At 100 Nodes — What to Expect

### Gossip Overhead

```
100 nodes, fan-out = 3, interval = 300ms:

  Convergence:    log₃(100) ≈ 4-5 rounds × 300ms ≈ 1.2s
  Messages/sec:   100 × (1000/300) × 3 ≈ 1,000 msgs/sec cluster-wide
  Per-node load:  ~10 msgs/sec received + ~10 sent = 20 msgs/sec
  Payload:        deltas only (not full state) → typically small
```

> [!success] Verdict on gossip at 100 nodes
> Totally manageable for metadata (config, topology, health). Watch payload size — keep gossip values small (< a few KB). Large blobs will hurt.

### Grain Actor Overhead

```
100 nodes hosting 10,000 grains:

  Routing:   ~1 network hop per request (cluster identity → PID)
  Cache hit: frequently-used grains cached locally → 0 hops
  Failover:  grain respawn typically < 1s after member death
  Memory:    ~10KB per grain (rough estimate, depends on state size)
             10,000 grains × 10KB ≈ 100MB across the cluster
```

---

## What Can Go Wrong

> [!warning] Gossip pitfalls
> - **Large state values**: gossip is optimized for small metadata. Putting MB-sized blobs in gossip state will cause network saturation.
> - **High-churn state**: if your gossip key changes thousands of times per second, deltas pile up faster than they're consumed.
> - **Treating gossip as a database**: gossip is eventually consistent. If you read immediately after a write, you may see stale data.

> [!warning] Grain pitfalls
> - **Hot grains**: if one grain (e.g., `"counter/global"`) receives all traffic, it becomes a bottleneck. Shard it: `"counter/shard-0"` through `"counter/shard-N"`.
> - **State loss on restart**: unless you use persistence, grain state resets on crash.
> - **Memory pressure**: millions of grains across a small cluster → tune your idle timeouts.

---

## Recommended Architecture for 100-Node State Sharing

```
┌────────────────────────────────────────────────────────┐
│  LAYER 1: Gossip (cluster self-knowledge)              │
│                                                        │
│  ● Topology: who's alive                               │
│  ● Heartbeats: actor counts, health                    │
│  ● Config/feature flags: small, read-heavy             │
│  ● Rarely-changing metadata                            │
│                                                        │
│  Convergence: ~1s    Consistency: eventual             │
└──────────────────────────┬─────────────────────────────┘
                           │ informs grain placement
┌──────────────────────────▼─────────────────────────────┐
│  LAYER 2: Grain Actors (domain state ownership)        │
│                                                        │
│  "account/42"   → Node 7   (balance, transactions)     │
│  "session/99"   → Node 23  (session tokens, expiry)    │
│  "product/7"    → Node 55  (inventory, price)          │
│                                                        │
│  Reads/writes:  route to the owning grain              │
│  Consistency:   strong (sequential per grain)          │
│  Failover:      automatic, < 1s                        │
└──────────────────────────┬─────────────────────────────┘
                           │ for durability
┌──────────────────────────▼─────────────────────────────┐
│  LAYER 3: Persistence (optional, for crash recovery)   │
│                                                        │
│  Proto.Actor Persistence: replay event log on restart  │
│  External DB: write-behind from grain on state change  │
└────────────────────────────────────────────────────────┘
```

---

## The Proto.Actor Way (Golden Rule)

> [!quote] Don't share state — route to the grain that owns it.
>
> The gossip layer is for the cluster to understand **itself** (membership, health). Your application state lives in **grain actors** — each grain is the single authoritative owner of one entity's state. Route every read and write there.

---

## Decision Checklist

Use **gossip** when:
- [ ] The value is small (< a few KB)
- [ ] All nodes need the same info
- [ ] Eventual consistency is acceptable
- [ ] The value changes infrequently (config, flags)
- [ ] You want local reads with zero network hop

Use **grain actors** when:
- [ ] The state belongs to a specific entity (user, order, session)
- [ ] You need strong consistency (no stale reads)
- [ ] Multiple nodes might write to the same state
- [ ] You want automatic failover with transparent re-routing

---

## References

- [Proto.Actor Docs — Cluster Gossip](https://proto.actor/docs/ProtoActor/cluster/gossip)
- [Proto.Actor Docs — Cluster](https://proto.actor/docs/ProtoActor/cluster)
- [Proto.Actor Docs — Persistence](https://proto.actor/docs/ProtoActor/persistence)
- [Proto.Actor Docs — Sharding](https://proto.actor/docs/ProtoActor/sharding-and-partitioning)

---

**← Prev** [[06-Gossip-Protocol]]
