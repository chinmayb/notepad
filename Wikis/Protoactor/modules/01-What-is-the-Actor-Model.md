---
tags: [protoactor, golang, actors, fundamentals]
module: 1
---

# Module 1 — What is the Actor Model?

> [!info] Goal
> Understand the mental model before touching any code.

**Next →** [[02-Core-Mechanics-in-Go]]

---

## The Big Idea

An actor is a **self-contained unit of computation**. It owns everything it needs and shares nothing with the outside world. The only way to interact with an actor is by **sending it a message**.

> [!tip] Analogy
> Think of actors like employees in a company. Each one has their own inbox (mailbox), handles one task at a time, and can delegate to sub-employees (children). They never share a desk (state) — they only communicate by sending memos (messages).

---

## What's Inside an Actor

```
┌─────────────────────────────────────┐
│               ACTOR                 │
│                                     │
│  PID ──────────────────────────────►│  ← only external handle
│                                     │
│  ┌─────────┐   ┌─────────────────┐  │
│  │ Mailbox │──►│    Behavior     │  │
│  │ (FIFO)  │   │ (message logic) │  │
│  └─────────┘   └────────┬────────┘  │
│                         │           │
│                    ┌────▼────┐      │
│                    │  State  │      │
│                    │(private)│      │
│                    └─────────┘      │
│                                     │
│  Children: [child1, child2, ...]    │
│  Supervisor Strategy: Restart/Stop  │
└─────────────────────────────────────┘
```

| Component               | What it does                                                          |
| ----------------------- | --------------------------------------------------------------------- |
| **PID**                 | Your only external handle — pass it around freely, cannot look inside |
| **Mailbox**             | FIFO queue of incoming messages — one processed at a time             |
| **State**               | Private variables — no locks needed, only this actor touches them     |
| **Behavior**            | The `Receive` function — can change dynamically at runtime            |
| **Children**            | Sub-actors spawned by this actor — forms a supervision tree           |
| **Supervisor Strategy** | Policy when a child crashes: Restart, Stop, or Escalate               |

---

## How Actors Communicate

```
┌──────────┐   send msg   ┌──────────┐
│ Actor A  │─────────────►│ Actor B  │
│          │              │ Mailbox  │
│          │◄─────────────│          │
└──────────┘   response   └──────────┘

Rules:
  - Messages are immutable
  - Delivery is asynchronous (fire and forget by default)
  - One message processed at a time — no concurrent access to state
  - Actor A does NOT wait while B processes (unless using request/future)
```

---

## Why Actors?

> [!success] Benefits
> - **No locks** — each actor is single-threaded in its own world
> - **No shared memory** — eliminates entire classes of race conditions
> - **Location transparent** — send to a PID; you don't care if it's local or remote
> - **Fault isolation** — a crash in one actor doesn't cascade unless you allow it

---

## Actor vs Traditional Thread Model

```
Traditional (shared memory):              Actor Model:
──────────────────────────────            ──────────────────────────────
Thread 1 ──┐                              Actor A  ──msg──►  Actor B
Thread 2 ──┼──► shared state             Actor C  ──msg──►  Actor B
Thread 3 ──┘    (needs locks!)
                                          No shared state. No locks.
Problem: deadlocks, race conditions       Messages are the only interface.
```

---

## The Actor Tree

Actors form a **hierarchy**. Every actor has a parent that supervises it.

```
                    [Actor System]
                          │
              ┌───────────┴───────────┐
              │                       │
         [Service A]             [Service B]
              │                       │
       ┌──────┴──────┐           ┌────┴────┐
  [Worker 1]   [Worker 2]   [Worker 3]  [Worker 4]
```

> [!note]
> This tree structure is what makes fault tolerance possible — a supervisor handles failures of its direct children.

---

## References

- [Proto.Actor Docs — Actors](https://proto.actor/docs/ProtoActor/actors)
- [Proto.Actor Docs — PID](https://proto.actor/docs/ProtoActor/pid)
- [Proto.Actor Docs — Mailboxes](https://proto.actor/docs/ProtoActor/mailboxes)

---

**Next →** [[02-Core-Mechanics-in-Go]]
