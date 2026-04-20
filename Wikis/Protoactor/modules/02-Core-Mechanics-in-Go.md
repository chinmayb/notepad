---
tags: [protoactor, golang, actors, hands-on]
module: 2
---

# Module 2 — Core Mechanics in Go

> [!info] Goal
> Write and run your first actor in Go.

**← Prev** [[01-What-is-the-Actor-Model]] | **Next →** [[03-Supervision-and-Fault-Tolerance]]

---

## Actor Lifecycle

Every actor goes through a predictable set of states:

```
                    ┌──────────┐
         spawn      │          │
    ───────────────►│ Starting │
                    │          │
                    └────┬─────┘
                         │ *actor.Started message delivered
                    ┌────▼─────┐
                    │          │
                    │  Running │◄──── processes messages from mailbox
                    │          │
                    └────┬─────┘
              crash │    │ stop requested
                    │    │
          ┌─────────▼┐  ┌▼──────────┐
          │Restarting│  │ Stopping  │
          └─────────┬┘  └─────┬─────┘
                    │         │
                    └────┬────┘
                         │
                    ┌────▼─────┐
                    │ Stopped  │
                    └──────────┘
```

| Lifecycle Message   | When it arrives                   |
| ------------------- | --------------------------------- |
| `*actor.Started`    | Right after actor is spawned      |
| `*actor.Stopping`   | Actor is about to stop            |
| `*actor.Stopped`    | Actor has fully stopped           |
| `*actor.Restarting` | Actor is restarting after a crash |

---

## Project Setup

```bash
mkdir counter-actor && cd counter-actor
go mod init counter-actor
go get github.com/asynkron/protoactor-go/actor
```

---

## Defining Messages

```go
// messages.go
package main

type Increment struct{}
type Decrement struct{}
type Reset     struct{}
type GetCount  struct{}
```

---

## Defining the Actor

```go
// counter.go
package main

import "github.com/asynkron/protoactor-go/actor"

type CounterActor struct {
    count int
}

func (c *CounterActor) Receive(ctx actor.Context) {
    switch ctx.Message().(type) {
    case *actor.Started:
        // called once when actor first spawns
        c.count = 0

    case *Increment:
        c.count++

    case *Decrement:
        c.count--

    case *Reset:
        c.count = 0

    case *GetCount:
        ctx.Respond(c.count)  // sends response back to sender
    }
}
```

---

## Spawning and Using the Actor

```go
// main.go
package main

import (
    "fmt"
    "time"
    "github.com/asynkron/protoactor-go/actor"
)

func main() {
    // 1. Create the actor system (one per application)
    system := actor.NewActorSystem()

    // 2. Define how to produce the actor
    props := actor.PropsFromProducer(func() actor.Actor {
        return &CounterActor{}
    })

    // 3. Spawn the actor — returns a PID
    pid := system.Root.Spawn(props)

    // 4. Fire-and-forget messages
    system.Root.Send(pid, &Increment{})
    system.Root.Send(pid, &Increment{})
    system.Root.Send(pid, &Decrement{})

    // 5. Request/response — waits for a reply
    future := system.Root.RequestFuture(pid, &GetCount{}, 5*time.Second)
    result, err := future.Result()
    if err == nil {
        fmt.Println("Count:", result) // Count: 1
    }

    // 6. Stop the actor
    system.Root.Stop(pid)
}
```

---

## Message Flow Diagram

```
main()                        CounterActor
  │                                │
  │── Send(&Increment{}) ─────────►│ count++ → 1
  │── Send(&Increment{}) ─────────►│ count++ → 2
  │── Send(&Decrement{}) ─────────►│ count-- → 1
  │                                │
  │── RequestFuture(&GetCount{}) ─►│
  │                                │── ctx.Respond(1)
  │◄────────────── 1 ──────────────│
  │
  │ prints: Count: 1
```

> [!note] One at a time
> Even though `main()` fires three messages rapidly, `CounterActor` processes them **one at a time** from its mailbox — no races, no locks.

---

## Tell vs Request vs RequestFuture

| Method                             | Blocking?         | Use when                                |
| ---------------------------------- | ----------------- | --------------------------------------- |
|  `Send(pid, msg)`                  | No                | Fire-and-forget, no reply needed        |
| `Request(pid, msg)`                | No                | Send with a return address (sender PID) |
| `RequestFuture(pid, msg, timeout)` | Yes (`.Result()`) | Need a response synchronously           |

---

## Spawning Child Actors

Actors can spawn children from inside `Receive`:

```go
func (a *ParentActor) Receive(ctx actor.Context) {
    switch ctx.Message().(type) {
    case *actor.Started:
        // spawn a child — parent automatically becomes supervisor
        childProps := actor.PropsFromProducer(func() actor.Actor {
            return &ChildActor{}
        })
        child := ctx.Spawn(childProps)
        ctx.Send(child, &DoWork{})
    }
}
```

---

## Mini Project

> [!todo] Build this
> Extend the counter actor to:
> 1. Accept a `SetCount{Value int}` message that jumps to a specific number
> 2. Print a log line on every state change using `*actor.Started` to set up a logger
> 3. Spawn a second "observer" actor that gets notified on every count change

---

## References

- [Proto.Actor Docs — Spawning Actors](https://proto.actor/docs/ProtoActor/spawn)
- [Proto.Actor Docs — Context](https://proto.actor/docs/ProtoActor/context)
- [Proto.Actor Docs — Communication](https://proto.actor/docs/ProtoActor/communication)
- [Go examples on GitHub](https://github.com/asynkron/protoactor-go/tree/dev/examples)

---

**← Prev** [[01-What-is-the-Actor-Model]] | **Next →** [[03-Supervision-and-Fault-Tolerance]]
