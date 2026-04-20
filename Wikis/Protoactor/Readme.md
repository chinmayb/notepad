---
tags:
  - protoactor
  - golang
  - index
study: gossip
---

# Proto.Actor in Go — Learning Space

> Reference: https://proto.actor/docs/ProtoActor/actors
> Topics: Actor Model · Proto.Actor in Go · Gossip Protocol · Distributed State Sharing

---

## Learning Path

```
01 → 02 → 03 → 04 → 05 → 06 → 07
```

| Module | Topic | Key Concept |
|---|---|---|
| [[modules/01-What-is-the-Actor-Model\|01]] | What is the Actor Model? | PID, Mailbox, State, Behavior, Supervision |
| [[modules/02-Core-Mechanics-in-Go\|02]] | Core Mechanics in Go | Spawn, Send, Receive, Lifecycle |
| [[modules/03-Supervision-and-Fault-Tolerance\|03]] | Supervision & Fault Tolerance | Restart, Stop, Escalate, Dead Letters |
| [[modules/04-Behaviors-and-Context\|04]] | Behaviors & Context | Become, EventStream, ReceiveTimeout, Stash |
| [[modules/05-Cluster-and-Virtual-Actors\|05]] | Cluster & Virtual Actors | Grains, Cluster Identity, Transparent Failover |
| [[modules/06-Gossip-Protocol\|06]] | Gossip Protocol | Epidemic spread, GossipState, Delta Sync |
| [[modules/07-State-Sharing-Across-100-Nodes\|07]] | State Sharing Across 100 Nodes | Gossip vs Grains, Architecture, Trade-offs |

---

## Quick Reference

### The Actor in One Picture

```
┌─────────────────────────────────────┐
│               ACTOR                 │
│  PID ──────────────────────────────►│  ← only external handle
│  ┌─────────┐   ┌─────────────────┐  │
│  │ Mailbox │──►│    Behavior     │  │
│  │ (FIFO)  │   │ (message logic) │  │
│  └─────────┘   └────────┬────────┘  │
│                    ┌────▼────┐      │
│                    │  State  │      │
│                    │(private)│      │
│                    └─────────┘      │
│  Supervisor Strategy: Restart/Stop  │
└─────────────────────────────────────┘
```

### Cluster in One Picture

```
Client → "user/123" (cluster identity)
              │
              ▼ cluster resolves
   [Node 1]  [Node 2: user/123]  [Node 3]
              │
         grain actor
         owns state
              │
   Node 2 dies → auto-respawn on Node 1
   Client never notices
```

### Gossip in One Picture

```
Round 1: A ──► B, C
Round 2: B ──► D, E  │  C ──► F, G
Round 3: all spread further...
Result: O(log N) rounds to converge
```

---

## Resources

- [Proto.Actor Docs — Actors](https://proto.actor/docs/ProtoActor/actors)
- [Proto.Actor Docs — Cluster](https://proto.actor/docs/ProtoActor/cluster)
- [Proto.Actor Docs — Gossip](https://proto.actor/docs/ProtoActor/cluster/gossip)
- [Go Cluster Example](https://github.com/asynkron/protoactor-go/tree/dev/examples/cluster-basic)
- [Go Virtual Actors Getting Started](https://proto.actor/docs/ProtoActor/cluster/getting-started-go)
