---
tags: [protoactor, golang, gossip, cluster, distributed-systems]
module: 6
---

# Module 6 — Gossip Protocol Deep Dive

> [!info] Goal
> Understand how cluster members share state without any central coordinator.

**← Prev** [[05-Cluster-and-Virtual-Actors]] | **Next →** [[07-State-Sharing-Across-100-Nodes]]

---

## What is a Gossip Protocol?

> [!quote]
> A gossip protocol spreads information the same way a rumor spreads — each node periodically tells a few random neighbors what it knows, who in turn tell others, until *everyone knows*. Also called an **epidemic protocol**.

No leader. No central registry. Pure peer-to-peer convergence.

---

## How Information Spreads

```
Initial state: Only Node A knows something new

Round 1:
  A ──► B  (A gossips to B and C)
  A ──► C

Round 2:
  B ──► D
  B ──► E
  C ──► F
  C ──► G

Round 3:
  D, E, F, G each gossip to 2 more...

Result: All N nodes converge in O(log N) rounds
```

### Convergence at Scale

```
Fan-out = 3, Interval = 300ms

  10 nodes:   log₃(10)  ≈ 2 rounds  → ~600ms
  100 nodes:  log₃(100) ≈ 4 rounds  → ~1.2s
  1000 nodes: log₃(1000)≈ 6 rounds  → ~1.8s

Wall-clock growth is logarithmic — scales well.
```

---

## GossipState Structure

Each cluster member maintains a `GossipState` — a nested dictionary:

```
GossipState (full cluster view on one node)
│
├── member-1 (GossipMemberState)
│   ├── "topology"    → {alive: [m1, m2, m3]}     seq: 42
│   ├── "heartbeat"   → {actorCounts: {...}}        seq: 105
│   └── "myAppState"  → {configVersion: "v3"}      seq: 7
│
├── member-2 (GossipMemberState)
│   ├── "topology"    → {alive: [m1, m2]}           seq: 38
│   └── "heartbeat"   → {...}                        seq: 99
│
└── member-3
    └── ...
```

> [!note] State ownership
> State is **per member** — you know what *each other node knows*. To compute a cluster-wide value, **merge entries from all members**. This gives eventual consistency.

---

## Delta Sync via Committed Offsets

To avoid resending the full state every gossip tick, Proto.Actor tracks **committed offsets** — the highest sequence number each node has confirmed sending to each peer.

```
Scenario: Node A gossips to Node B

Node A's committed offsets for B:
  "member1.topology"  → last confirmed at seq 38
  "member1.heartbeat" → last confirmed at seq 99

Node A's current state:
  "member1.topology"  → current seq 42
  "member1.heartbeat" → current seq 105

Delta sent to B:
  topology  seq 39, 40, 41, 42   (only the new stuff)
  heartbeat seq 100, 101, 102, 103, 104, 105

No full-state retransmission → bandwidth stays manageable.
```

---

## Gossip Fan-out

```
ClusterConfig.GossipFanout = 3   (default)

Every 300ms, each node:
  1. Picks 3 random cluster members
  2. Sends them a delta of its state
  3. Receives their deltas in return
  4. Merges incoming state into local GossipState

Message volume (100 nodes, fan-out 3, interval 300ms):
  100 × (1000/300) × 3 ≈ 1000 gossip msgs/sec cluster-wide
  Each msg = delta only → manageable
```

---

## Built-in Gossip Keys

Proto.Cluster uses gossip internally for:

| Key | Purpose |
|---|---|
| `topology` | Which members are alive, cluster generation |
| `heartbeat` | Per-member actor counts and health stats |
| `banned-members` | Members that have been blocklisted |

---

## User-Defined State (Go)

You can gossip your own application state — any Protobuf message:

```go
// Define your state as a Protobuf message (state.proto):
// message MyFeatureFlags {
//   bool enable_new_ui   = 1;
//   int32 max_batch_size = 2;
// }

// Set state — automatically gossiped to all cluster members
cluster.Gossip.SetKey("featureFlags", &MyFeatureFlags{
    EnableNewUi:  true,
    MaxBatchSize: 500,
})

// Read state from all known members (eventual consistency)
states, err := cluster.Gossip.GetState("featureFlags")
for memberId, anyValue := range states {
    var flags MyFeatureFlags
    if err := anyValue.UnmarshalTo(&flags); err == nil {
        fmt.Printf("Member %s: enableNewUI=%v\n", memberId, flags.EnableNewUi)
    }
}
```

---

## Subscribing to Gossip Changes

```go
// React whenever a gossip key is updated on any member
system.EventStream.Subscribe(func(evt interface{}) {
    update, ok := evt.(*cluster.GossipUpdate)
    if !ok || update.Key != "featureFlags" {
        return
    }

    var flags MyFeatureFlags
    if err := update.Value.UnmarshalTo(&flags); err == nil {
        fmt.Printf("featureFlags changed on %s: %+v\n", update.MemberID, flags)
    }
})
```

---

## Tuning Gossip

```go
clusterConfig := cluster.Configure(...).
    WithGossipInterval(500 * time.Millisecond).      // default: 300ms
    WithGossipRequestTimeout(2 * time.Second).       // default: 1500ms
    WithGossipFanout(5)                              // default: 3
```

> [!warning] Trade-offs
> - **Shorter interval** → faster convergence, more CPU/network overhead
> - **Longer interval** → less overhead, slower convergence
> - **Higher fan-out** → fewer rounds to converge, more messages per round

---

## Gossip vs Central Coordinator

```
Central Coordinator (e.g. Zookeeper):          Gossip:
──────────────────────────────────             ──────────────────────────────
All nodes read/write one place                 All nodes talk to random peers
Strong consistency                             Eventual consistency
Single point of failure                        No SPOF — any node can fail
Bottleneck at scale                            Scales with log(N) convergence
Complex failover                               Self-healing by design
```

---

## Mini Project

> [!todo] Build this
> 1. Start a 3-node cluster
> 2. On Node 1, call `Gossip.SetKey("config", &AppConfig{FeatureX: true})`
> 3. On Node 2 and 3, subscribe to `GossipUpdate` for `"config"`
> 4. Measure: how many milliseconds until all nodes see the update?
> 5. Try increasing `GossipInterval` to 2s — how does convergence time change?

---

## References

- [Proto.Actor Docs — Cluster Gossip](https://proto.actor/docs/ProtoActor/cluster/gossip)
- [Wikipedia — Gossip Protocol](https://en.wikipedia.org/wiki/Gossip_protocol)

---

**← Prev** [[05-Cluster-and-Virtual-Actors]] | **Next →** [[07-State-Sharing-Across-100-Nodes]]
