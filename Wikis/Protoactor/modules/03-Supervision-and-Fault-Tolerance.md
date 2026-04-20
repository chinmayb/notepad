---
tags: [protoactor, golang, supervision, fault-tolerance]
module: 3
---

# Module 3 — Supervision & Fault Tolerance

> [!info] Goal
> Understand the "let it crash" philosophy and how Proto.Actor handles failures.

**← Prev** [[02-Core-Mechanics-in-Go]] | **Next →** [[04-Behaviors-and-Context]]

---

## The Core Philosophy

> [!quote] Let it crash
> Don't defensively code every possible failure. Instead, design your supervision tree so failures are *expected*, *isolated*, and *handled at the right level*.

In Proto.Actor, every actor has a parent. When an actor fails, its **parent (supervisor) decides what to do** — not the actor itself.

---

## Supervisor Tree

```
                    [Actor System Root]
                           │
               ┌───────────┴────────────┐
               │                        │
          [ServiceA]               [ServiceB]
          supervisor                supervisor
               │                        │
       ┌───────┴───────┐           ┌────┴────┐
  [Worker 1]     [Worker 2]   [Worker 3]  [Worker 4]
```

- Each node supervises its **direct children only**
- Crashes don't automatically propagate upward unless escalated
- You decide granularity — fine-grained or coarse

---

## What Happens When an Actor Crashes

```
[Worker 2] panics
      │
      ▼
[ServiceA] receives failure notification
      │
      ├─► Restart  ──► Worker 2 gets fresh state, receives *Started again
      │
      ├─► Stop     ──► Worker 2 is terminated, mailbox drained to dead letters
      │
      └─► Escalate ──► ServiceA's parent (Root) now handles it
```

> [!warning] State on restart
> When an actor **restarts**, its state is **reset from scratch** — like a fresh spawn. If you need to survive a restart, use [[07-State-Sharing-Across-100-Nodes|Persistence]].

---

## Supervision Strategies

### OneForOne (default)
Only the failing child is affected.

```
[Parent]
  ├── [Child A] ← crashes → restarted
  ├── [Child B] ← unaffected
  └── [Child C] ← unaffected
```

### AllForOne
If one child fails, all children are affected by the same directive.

```
[Parent]
  ├── [Child A] ← crashes → all restarted
  ├── [Child B] ←            restarted
  └── [Child C] ←            restarted
```

Use `AllForOne` when children are tightly coupled and a partial failure leaves the system in an inconsistent state.

---

## Defining a Supervisor Strategy in Go

```go
// supervisor.go
package main

import (
    "fmt"
    "github.com/asynkron/protoactor-go/actor"
)

// Worker that can panic
type RiskyWorker struct{}

func (r *RiskyWorker) Receive(ctx actor.Context) {
    switch msg := ctx.Message().(type) {
    case *DoWork:
        if msg.ShouldFail {
            panic("something went wrong!")
        }
        fmt.Println("work done")
    }
}

// Supervisor that defines restart policy
type Supervisor struct{}

func (s *Supervisor) Receive(ctx actor.Context) {
    switch ctx.Message().(type) {
    case *actor.Started:
        workerProps := actor.PropsFromProducer(func() actor.Actor {
            return &RiskyWorker{}
        })
        ctx.Spawn(workerProps)
    }
}

// Attach a supervision strategy to props
func supervisorProps() *actor.Props {
    return actor.
        PropsFromProducer(func() actor.Actor { return &Supervisor{} }).
        WithSupervisor(actor.NewOneForOneStrategy(
            3,                    // max retries
            1000*time.Millisecond, // within this window
            func(reason interface{}) actor.Directive {
                fmt.Println("supervising failure:", reason)
                return actor.RestartDirective
            },
        ))
}
```

---

## Supervision Decision Flow

```
Actor fails
    │
    ▼
Supervisor's decider function called with (reason)
    │
    ├─► actor.RestartDirective   → restart the child
    ├─► actor.StopDirective      → stop the child permanently
    ├─► actor.EscalateDirective  → escalate to my own parent
    └─► actor.ResumeDirective    → ignore and continue (use with care)
```

---

## Restart vs Stop vs Escalate — When to Use Each

| Directive | Use when |
|---|---|
| **Restart** | Transient failures — bad message, temporary resource unavailability |
| **Stop** | Fatal errors — actor cannot meaningfully recover |
| **Escalate** | Parent doesn't know how to handle it — let a higher supervisor decide |
| **Resume** | Non-critical errors where state is still valid |

---

## Dead Letters

When a stopped actor's mailbox drains, all remaining messages go to **Dead Letters**:

```
Actor stopped
    │
    ▼
Remaining mailbox messages ──► Dead Letter Mailbox
                                        │
                                        ▼
                               EventStream as DeadLetterEvent

// Subscribe to dead letters for observability
system.EventStream.Subscribe(func(evt interface{}) {
    if dl, ok := evt.(*actor.DeadLetterEvent); ok {
        fmt.Printf("dead letter: %v sent to %v\n", dl.Message, dl.PID)
    }
})
```

> [!warning]
> Dead letters are best-effort. Don't rely on them for guaranteed delivery. Use request/response patterns or acknowledgment protocols if you need delivery guarantees.

---

## Mini Project

> [!todo] Build this
> 1. Create a `DatabaseActor` that randomly panics 30% of the time
> 2. Wrap it in a supervisor that restarts it up to 5 times within 2 seconds
> 3. Subscribe to `DeadLetterEvent` and log any dropped messages
> 4. Observe: what happens after the 5th failure in 2 seconds?

---

## References

- [Proto.Actor Docs — Supervision](https://proto.actor/docs/ProtoActor/supervision)
- [Proto.Actor Docs — Actor Lifecycle](https://proto.actor/docs/ProtoActor/life-cycle)
- [Proto.Actor Docs — EventStream](https://proto.actor/docs/ProtoActor/eventstream)

---

**← Prev** [[02-Core-Mechanics-in-Go]] | **Next →** [[04-Behaviors-and-Context]]
