---
title: "Learning Goals: DNSSEC"
type: source
tags: ["dnssec", "learning-goals", "dns"]
created: 2026-04-13
updated: 2026-04-13
source_count: 1
confidence: high
status: active
related: ["[[dnssec-learning-roadmap]]", "[[dnssec-internals]]", "[[chain-of-trust]]", "[[dnssec-vs-dns-resolution]]", "[[dnssec-in-bind]]", "[[dnssec-debugging]]"]
source_url: ""
source_author: "cbharadwaj"
source_date: 2026-04-13
source_type: personal_note
open_questions:
  - "What specific BIND version is the target environment?"
  - "Is this for operating a signed authoritative zone, validating resolver, or both?"
  - "What does 'debugging in future' mean in practice — production incidents, lab, or both?"
contradictions: []
entities_mentioned: []
concepts_mentioned:
  - "[[dnssec-internals]]"
  - "[[dnssec-vs-dns-resolution]]"
  - "[[chain-of-trust]]"
  - "[[dnssec-in-bind]]"
  - "[[dnssec-debugging]]"
---

# Learning Goals: DNSSEC

> **Type**: Personal learning intent document
> **Author**: cbharadwaj · **Date**: 2026-04-13
> **Role in wiki**: Curriculum anchor — defines the 6 learning axes that all DNSSEC source pages map back to

---

## Summary

This document captures the owner's stated learning goals for DNSSEC. It is the **north star** for this knowledge base: every source ingested should be evaluated against how well it advances one or more of these six axes. The wiki's synthesis pages, concept pages, and source summaries are all organized around these goals.

This is not a content source — it contains no facts to extract. Its value is as a **structured intent** that shapes how all other sources are processed and connected.

---

## The Six Learning Axes

These are the six dimensions the owner wants to achieve mastery across:

### Axis 1 — DNSSEC Internals
*The cryptographic machinery underneath.*
- Record types: DNSKEY, RRSIG, DS, NSEC, NSEC3
- KSK vs ZSK — the two-key architecture and why it exists
- How signatures are created, stored, and verified
- EDNS0 and the DO (DNSSEC OK) bit
- → See [[dnssec-internals]]

### Axis 2 — Workflow: DNSSEC vs Normal DNS Resolution
*What actually changes, step by step.*
- How a validating resolver walks the chain vs a plain resolver
- Where RRSIG records appear in the resolution path
- Caching behavior differences
- What happens at each delegation boundary
- → See [[dnssec-vs-dns-resolution]]

### Axis 3 — First Principles
*Why DNSSEC exists. What it protects against. What it doesn't.*
- The Kaminsky cache poisoning attack (2008) — the canonical motivation
- What "authentication without confidentiality" means
- The threat model: on-path vs off-path attackers, cache poisoning vs zone tampering
- What DNSSEC does NOT protect against (DDoS, resolver compromise, registrar hijack)
- → See [[chain-of-trust]]

### Axis 4 — Practical Use
*Deploying and operating signed zones.*
- Signing a zone: `dnssec-keygen`, `dnssec-signzone`, or `dnssec-policy`
- DS record submission to parent / registrar
- Key rollover procedures (ZSK rollover vs KSK rollover — very different!)
- Signature expiry and re-signing schedules
- → See [[dnssec-in-bind]]

### Axis 5 — DNSSEC in BIND
*BIND-specific configuration and operation.*
- `named.conf` directives: `dnssec-validation`, `dnssec-policy`, `auto-dnssec`
- Inline signing vs manual signing workflow
- Key management with `dnssec-keymgr` and the newer `dnssec-policy`
- Zone file behavior when BIND manages signing automatically
- → See [[dnssec-in-bind]]

### Axis 6 — Debugging and Tools
*Diagnosing problems when DNSSEC breaks.*
- `dig +dnssec` — the baseline query tool
- `delv` — DNSSEC-aware validation tool (BIND 9.10+)
- `dnsviz` — visual chain-of-trust debugger
- `dnssec-verify` — offline zone verification
- Common failure modes: signature expiry, missing DS, broken delegation, algorithm mismatch
- SERVFAIL anatomy under DNSSEC validation
- → See [[dnssec-debugging]]

---

## Key Claims

1. The owner's primary purpose is operational: learning DNSSEC well enough to configure it in BIND and debug production issues — not just theoretical understanding. `confidence: high`
2. BIND is the specific implementation target — content about other DNS servers (PowerDNS, NSD, Unbound) is secondary. `confidence: high`
3. Debugging is explicitly forward-looking ("issues that may arise") — implies a production or near-production environment context. `confidence: medium`

---

## Coverage Map

How well the current raw source collection covers each axis (as of 2026-04-13):

| Axis | Coverage | Best Current Source | Gap |
|------|----------|---------------------|-----|
| Internals | ✅ Good | CS161 Berkeley (`DNSSEC.md`) | EDNS0/DO bit depth |
| Workflow delta | ✅ Good | CS161 Berkeley (`DNSSEC.md`) | Covered |
| First principles | ⚠️ Partial | Cloudflare intro | Kaminsky attack story missing |
| Practical use | ⚠️ Partial | `DNSSEC Guide — BIND 9` | Key rollover detail thin |
| DNSSEC in BIND | ⚠️ Partial | `DNSSEC and BIND 9.md` | `dnssec-policy` depth needed |
| Debugging & tools | ⚠️ Partial | `dig and delv.md` | `dnsviz`, SERVFAIL triage missing |

---

## Concepts Covered

- [[dnssec-internals]]
- [[dnssec-vs-dns-resolution]]
- [[chain-of-trust]]
- [[dnssec-in-bind]]
- [[dnssec-debugging]]

---

*Source file: `raw/01-learn-dnssec.md` | This page is the curriculum anchor for the DNSSEC knowledge base.*
