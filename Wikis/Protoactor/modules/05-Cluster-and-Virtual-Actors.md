---
tags: [protoactor, golang, cluster, grains, virtual-actors, distributed]
module: 5
---

# Module 5 — Cluster & Virtual Actors

> [!info] Goal
> Run actors across multiple machines using Proto.Cluster and virtual actors (grains).

**← Prev** [[04-Behaviors-and-Context]] | **Next →** [[06-Gossip-Protocol]]

---

## Three Deployment Modes

```
┌─────────────────────────────────────────────────────┐
│  LOCAL (single process)                             │
│                                                     │
│  ┌────────┐    msg    ┌────────┐                    │
│  │Actor A │──────────►│Actor B │                    │
│  └────────┘           └────────┘                    │
│                                                     │
│  Both actors in same memory space                   │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│  REMOTE (two processes, different hosts)            │
│                                                     │
│  [Host 1]               [Host 2]                    │
│  ┌────────┐   gRPC/PB   ┌────────┐                  │
│  │Actor A │────────────►│Actor B │                  │
│  └────────┘             └────────┘                  │
│                                                     │
│  Addressed by PID. You manage placement.            │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│  CLUSTER (virtual actors / grains)                  │
│                                                     │
│  Client → "user/123"  (cluster identity)            │
│                │                                    │
│                ▼  cluster resolves location         │
│  [Node 1]   [Node 2]       [Node 3]                 │
│  ┌───────┐  ┌──────────┐  ┌───────┐                 │
│  │       │  │ user/123 │  │       │                 │
│  │       │  │  grain   │  │       │                 │
│  └───────┘  └──────────┘  └───────┘                 │
│                                                     │
│  Node 2 dies → user/123 respawns on Node 1 or 3     │
│  Client never notices. Location is transparent.     │
└─────────────────────────────────────────────────────┘
```

---

## Virtual Actors (Grains)

Proto.Cluster uses the **Virtual Actor Model** pioneered by Microsoft Orleans.

> [!tip] Key Idea
> You never manually spawn or locate a grain. You just address it by *kind + identity* and the cluster does the rest — spawning it on demand, moving it on failure.

### Grain vs Regular Actor

| Aspect | Regular Actor | Virtual Actor (Grain) |
|---|---|---|
| Referenced by | PID | `kind/identity` e.g. `user/123` |
| Spawned by | You, explicitly | Cluster, on first request |
| Location | You decide | Cluster decides |
| Failover | Manual re-spawn | Automatic, transparent |
| Lifecycle | You manage | Cluster manages |
| Multiple instances | Possible | Exactly one active at a time |

---

## Cluster Identity

A grain is uniquely addressed by two parts:

```
"user/123"
  │    │
  │    └── identity: unique ID within the kind (maps to domain entity)
  └── kind: the grain type (maps to a registered ClusterKind)

More examples:
  "vehicle/JB007"
  "airport/ATL"
  "temperatureSensor/ff7e1600-ecdd-4150-80d3-17dd1074bcc0"
```

---

## Grain Lifecycle

```
1. Cluster starts — grain "user/123" does not exist anywhere

2. Client sends request to cluster identity "user/123"
        │
        ▼
3. Cluster spawns "user/123" as an actor on Node 2
   → actor receives *actor.Started
   → request is delivered

4. Node 2 goes down
        │
        ▼
5. Cluster detects failure via gossip / heartbeat
   → "user/123" is respawned on Node 1 (new PID)
   → next request routes to Node 1 automatically

6. Grain idle for too long → set ReceiveTimeout → poison self
```

---

## Setting Up a Cluster in Go

```go
// main.go
package main

import (
    "github.com/asynkron/protoactor-go/actor"
    "github.com/asynkron/protoactor-go/cluster"
    "github.com/asynkron/protoactor-go/cluster/clusterproviders/consul"
    "github.com/asynkron/protoactor-go/remote"
)

func main() {
    system := actor.NewActorSystem()

    // 1. Configure remote transport
    remoteConfig := remote.Configure("127.0.0.1", 8080)

    // 2. Register grain kinds
    userGrainKind := cluster.NewKind(
        "user",
        actor.PropsFromProducer(func() actor.Actor { return &UserGrain{} }),
    )

    // 3. Configure and start the cluster
    clusterConfig := cluster.Configure(
        "my-cluster",
        consul.New(),        // cluster provider (tracks members)
        remoteConfig,
        cluster.WithKinds(userGrainKind),
    )

    c := cluster.New(system, clusterConfig)
    c.StartMember()          // joins the cluster as a member
    // or c.StartClient()    // joins as client-only (no grain hosting)
}
```

---

## Sending to a Grain

```go
// No need to know where "user/123" lives
// The cluster resolves it automatically

ctx := c.GetClusterIdentity("user", "123")
response, err := ctx.RequestFuture(&GetProfile{}, 5*time.Second).Result()
```

---

## Proto.Cluster Architecture

```
┌──────────────────────────────────────────────────┐
│               Proto.Cluster                      │
│                                                  │
│  ┌──────────────┐    ┌───────────────────────┐   │
│  │ Cluster      │    │  Identity Lookup       │   │
│  │ Provider     │    │  (where is grain X?)   │   │
│  │ (Consul/K8s) │    └───────────┬───────────┘   │
│  └──────┬───────┘                │               │
│         │ member list            │ PID            │
│         ▼                        ▼               │
│  ┌──────────────┐    ┌───────────────────────┐   │
│  │   Gossip     │    │  Actor Cache           │   │
│  │  (state sync)│    │  (recently used PIDs)  │   │
│  └──────────────┘    └───────────────────────┘   │
│                                                  │
│  Underlying: Proto.Remote (gRPC + Protobuf)      │
└──────────────────────────────────────────────────┘
```

---

## Shutting Down a Grain

Grains don't garbage collect automatically. Two patterns:

```go
// Pattern 1: Event-driven stop
func (u *UserGrain) Receive(ctx actor.Context) {
    switch ctx.Message().(type) {
    case *UserDisconnected:
        ctx.Stop(ctx.Self()) // terminates this grain
    }
}

// Pattern 2: Idle timeout
func (u *UserGrain) Receive(ctx actor.Context) {
    switch ctx.Message().(type) {
    case *actor.Started:
        ctx.SetReceiveTimeout(10 * time.Minute) // idle threshold

    case *actor.ReceiveTimeout:
        ctx.Stop(ctx.Self())
    }
}
```

---

## Mini Project

> [!todo] Build this
> 1. Start two cluster members on ports 8080 and 8081
> 2. Register a `"counter"` grain kind
> 3. From a client node, send `Increment` and `GetCount` to `"counter/global"`
> 4. Kill one member, verify the grain migrates and count is reset (stateless restart)
> 5. Bonus: persist count using Proto.Actor persistence to survive restarts

---

## References

- [Proto.Actor Docs — Cluster](https://proto.actor/docs/ProtoActor/cluster)
- [Getting Started With Grains (Go)](https://proto.actor/docs/ProtoActor/cluster/getting-started-go)
- [Virtual Actors (Go)](https://proto.actor/docs/ProtoActor/cluster/virtual-actors-go)
- [Go Cluster Example](https://github.com/asynkron/protoactor-go/tree/dev/examples/cluster-basic)

---

**← Prev** [[04-Behaviors-and-Context]] | **Next →** [[06-Gossip-Protocol]]
