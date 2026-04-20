---
tags: [protoactor, golang, behaviors, context, eventstream]
module: 4
---

# Module 4 — Behaviors & Context

> [!info] Goal
> Make actors dynamic: swap behavior at runtime and use the actor context and EventStream.

**← Prev** [[03-Supervision-and-Fault-Tolerance]] | **Next →** [[05-Cluster-and-Virtual-Actors]]

---

## Behaviors — Runtime Behavior Switching

An actor's `Receive` function can be **swapped at runtime**. This is perfect for modeling state machines — the actor *becomes* something different in response to events.

### Behavior Stack

```
Initial Behavior (bottom of stack)
      │
      │  ctx.Become(newBehavior)     ← replaces top
      │  ctx.BecomeStacked(new)      ← pushes onto stack
      │  ctx.UnbecomeStacked()       ← pops back to previous
      ▼

Stack view with BecomeStacked:

  ┌─────────────────┐  ← current (top)
  │  Authenticated  │
  ├─────────────────┤
  │ Unauthenticated │  ← previous (bottom)
  └─────────────────┘

UnbecomeStacked() → pops Authenticated → back to Unauthenticated
```

---

## Auth State Machine Example

```go
// auth_actor.go
package main

import "github.com/asynkron/protoactor-go/actor"

type Login  struct{ Token string }
type Logout struct{}
type DoWork struct{ Payload string }

type AuthActor struct {
    behavior actor.Behavior
}

func NewAuthActor() *AuthActor {
    a := &AuthActor{}
    a.behavior.Become(a.unauthenticated)
    return a
}

func (a *AuthActor) Receive(ctx actor.Context) {
    a.behavior.Receive(ctx)
}

// State: not logged in
func (a *AuthActor) unauthenticated(ctx actor.Context) {
    switch ctx.Message().(type) {
    case *Login:
        fmt.Println("logged in, switching behavior")
        a.behavior.BecomeStacked(a.authenticated)
    default:
        fmt.Println("rejected — not authenticated")
    }
}

// State: logged in
func (a *AuthActor) authenticated(ctx actor.Context) {
    switch ctx.Message().(type) {
    case *Logout:
        fmt.Println("logged out")
        a.behavior.UnbecomeStacked()
    case *DoWork:
        fmt.Println("doing work:", ctx.Message().(*DoWork).Payload)
    }
}
```

### Flow Diagram

```
Message: &Login{}
    │  unauthenticated handler
    │  BecomeStacked(authenticated)
    ▼
Message: &DoWork{Payload: "process order"}
    │  authenticated handler
    │  prints: doing work: process order
    ▼
Message: &Logout{}
    │  authenticated handler
    │  UnbecomeStacked()
    ▼
Message: &DoWork{...}
    │  back to unauthenticated handler
    │  prints: rejected — not authenticated
```

---

## Context — Your Window into the Actor System

The `actor.Context` passed into every `Receive` call gives you access to everything:

```go
func (a *MyActor) Receive(ctx actor.Context) {
    ctx.Self()          // this actor's own PID
    ctx.Sender()        // PID of who sent the current message
    ctx.Parent()        // PID of this actor's supervisor
    ctx.Message()       // the current message (interface{})
    ctx.Respond(reply)  // send reply back to Sender
    ctx.Spawn(props)    // spawn a child actor
    ctx.Stop(pid)       // stop an actor
    ctx.Stash()         // stash current message for later
    ctx.Unwatch(pid)    // stop watching a PID for termination
    ctx.Watch(pid)      // get notified when pid terminates
}
```

### Watch / Unwatch Pattern

```go
// Watch another actor — receive Terminated when it stops
func (a *MonitorActor) Receive(ctx actor.Context) {
    switch msg := ctx.Message().(type) {
    case *actor.Started:
        ctx.Watch(a.watchedPID)

    case *actor.Terminated:
        fmt.Println("watched actor terminated:", msg.Who)
        // take action: restart, alert, etc.
    }
}
```

---

## EventStream — System-Wide Pub/Sub

The `EventStream` is a publish/subscribe bus built into every actor system. Actors (and non-actors) can publish and subscribe without knowing each other's PIDs.

```
Publisher                EventStream               Subscriber(s)
    │                        │                          │
    │── Publish(UserLoggedIn) ──►                       │
    │                        │──── deliver to all ─────►│
    │                        │     matching subs        │
```

```go
// Publish from anywhere
system.EventStream.Publish(&UserLoggedIn{UserId: "123"})

// Subscribe — returns a subscription handle
sub := system.EventStream.Subscribe(func(evt interface{}) {
    if e, ok := evt.(*UserLoggedIn); ok {
        fmt.Println("user logged in:", e.UserId)
    }
})

// Unsubscribe when done
system.EventStream.Unsubscribe(sub)
```

> [!tip] Type-based routing
> The EventStream dispatches by **Go type**. Your subscriber only receives the types it checks for. No message routing config needed.

---

## Receive Timeout

Make an actor react when it hasn't received a message for a while — useful for idle grain cleanup:

```go
func (a *SessionActor) Receive(ctx actor.Context) {
    switch ctx.Message().(type) {
    case *actor.Started:
        ctx.SetReceiveTimeout(30 * time.Second)

    case *actor.ReceiveTimeout:
        fmt.Println("idle too long, stopping")
        ctx.Stop(ctx.Self())

    case *Ping:
        // reset timeout on activity
        ctx.SetReceiveTimeout(30 * time.Second)
    }
}
```

---

## Stash — Defer Messages

If an actor isn't ready to handle certain messages yet, stash them. They'll be replayed after the next `Unstash()`:

```go
func (a *InitializingActor) Receive(ctx actor.Context) {
    switch ctx.Message().(type) {
    case *actor.Started:
        // not ready yet — stash incoming work
    case *InitDone:
        ctx.UnstashAll() // replay everything that arrived during init
    case *DoWork:
        ctx.Stash()      // queue it for later
    }
}
```

---

## Mini Project

> [!todo] Build this
> Create a **vending machine actor** using behavior switching:
> 1. `Idle` state: accepts `InsertCoin` → transitions to `HasCoin`
> 2. `HasCoin` state: accepts `SelectItem` → dispenses, returns to `Idle`; or `ReturnCoin` → back to `Idle`
> 3. Any message in wrong state → print "invalid operation"

---

## References

- [Proto.Actor Docs — Behaviors](https://proto.actor/docs/ProtoActor/behaviors)
- [Proto.Actor Docs — Context](https://proto.actor/docs/ProtoActor/context)
- [Proto.Actor Docs — EventStream](https://proto.actor/docs/ProtoActor/eventstream)
- [Proto.Actor Docs — Receive Timeout](https://proto.actor/docs/ProtoActor/receive-timeout)

---

**← Prev** [[03-Supervision-and-Fault-Tolerance]] | **Next →** [[05-Cluster-and-Virtual-Actors]]
