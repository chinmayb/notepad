---
title: "Overview"
type: overview
tags: ["dnssec", "dns"]
created: 2026-04-10
updated: 2026-04-20
source_count: 8
confidence: medium
status: active
related: ["[[dnssec-learning-roadmap]]", "[[dnssec-internals]]", "[[chain-of-trust]]", "[[dnssec-vs-dns-resolution]]", "[[dnssec-in-bind]]", "[[dnssec-debugging]]", "[[dns-cache-poisoning]]"]
open_questions:
  - "What BIND version is in use in the target environment?"
  - "Is the goal authoritative zone signing, validating resolver, or both?"
contradictions: []
---

# Overview

> **Domain**: DNSSEC — DNS Security Extensions
> **Goal**: Operational mastery of DNSSEC from first principles through BIND configuration and production debugging

---

## Current Focus

This knowledge base is being built to achieve deep, practical understanding of DNSSEC — the cryptographic security layer for DNS. The owner's six learning axes are captured in [[01-learn-dnssec]] and the reading plan is in [[dnssec-learning-roadmap]].

**11 raw sources** are staged in `raw/` covering theory, RFC specifications, BIND operations, and debugging tools. Ingestion is in progress (8 of 11 ingested). The BIND operational layer is now fully documented (3 BIND sources + RFC 6781). The conceptual/pedagogical layer is complete (CS161 ingested; all concept pages are now active). Remaining: RFC 9364 reference, and 2 debugging tool sources.

---

## Evolving Thesis

*Mid-stage — thesis solidifying as BIND ops sources are ingested.*

**Current thesis**: DNSSEC is a chain-of-trust protocol whose operational complexity (especially key rollovers and BIND configuration) vastly exceeds its conceptual complexity. Most DNSSEC failures in production trace to either (a) broken parent-child DS delegation (security lameness) or (b) signature expiry — both preventable with proper automation via BIND's `dnssec-policy`. Two additional silent killers are EDNS0 filtering and clock skew, both common in enterprise environments and both producing SERVFAIL with no obvious cause.

---

## Key Themes

### 1. Authentication Without Confidentiality
DNSSEC signs DNS responses so forgeries are detectable, but does not encrypt them. The threat model is specifically **data integrity and origin authentication** — not privacy.

### 2. The Chain of Trust is Everything
Every DNSSEC concept ultimately serves one purpose: verifying the chain from the root trust anchor to the final answer. A single broken link causes SERVFAIL. See [[chain-of-trust]] for the full 12-step validation walkthrough.

### 3. BIND's `dnssec-policy` is the Modern Path
Manual `dnssec-keygen` + `dnssec-signzone` still works but is error-prone. BIND 9.16+ `dnssec-policy` automates the key lifecycle. Any new deployment should use it. **ISC explicitly flags all pre-KASP documentation as outdated** — confirmed by [[dnssec-and-bind9-isc]] and [[dnssec-guide-bind9]].

### 4. The `default` Policy Uses a CSK, Not KSK/ZSK
`dnssec-policy default` creates a **Combined Signing Key** — one key pair serving both roles. This simplifies management but means every rollover requires parent DS coordination. For zones needing frequent ZSK rotation, use a custom policy with explicit `ksk`/`zsk` lines. Documented in [[dnssec-internals]] and [[chain-of-trust]].

### 5. Debugging Requires Specialized Tools
`dig` alone is insufficient for DNSSEC debugging. `delv` provides chain validation traces; [[DNSViz]] visualizes the full chain. Understanding SERVFAIL anatomy is critical for production. The full triage methodology is in [[dnssec-debugging]].

### 6. Two Hidden DNSSEC Killers: EDNS0 Filtering and Clock Skew
DNSSEC can be perfectly configured and still fail silently in production due to (a) network devices stripping EDNS0 packets, or (b) NTP clock sync issues making valid RRSIG records appear expired. Both produce SERVFAIL with no obvious cause. Documented in [[dnssec-in-bind]] and [[dnssec-debugging]].

---

| Concept | Status | Description |
|---------|--------|-------------|
| [[chain-of-trust]] | **Active (high)** | The cryptographic hierarchy from root to answer; 12-step validation; DS record mechanics; KSK rollover; CSK |
| [[dnssec-internals]] | **Active (high)** | RRsets, DNSKEY/RRSIG/DS/NSEC/NSEC3, KSK/ZSK/CSK, key lifecycle, EDNS0/DO bit |
| [[dns-cache-poisoning]] | **Active (high)** | The Kaminsky attack; why DNSSEC exists |
| [[dnssec-in-bind]] | **Active (high)** | BIND 9 config, `dnssec-policy`, signing, validation, rollover, `rndc` commands |
| [[dnssec-debugging]] | **Active (high)** | `ad` flag, security lameness, clock skew, NTA, triage methodology |
| [[dnssec-vs-dns-resolution]] | **Active (high)** | Full 6-step `dig` walkthrough (eecs.berkeley.edu); DNSKEY flags; OPT pseudosection; offline signing motivation |

---

## Source Inventory

| Source | Type | Ingested | Primary Axes |
|--------|------|----------|--------------|
| [[01-learn-dnssec]] | Goals anchor | ✅ | All |
| [[how-does-dnssec-work-cloudflare]] | Tutorial | ✅ | Internals, chain of trust |
| [[kaminsky-exploit]] | Retrospective | ✅ | First principles |
| [[dnssec-and-bind9-isc]] | ISC portal | ✅ | BIND ops prereqs, tooling |
| [[dnssec-guide-bind9]] | ISC BIND guide | ✅ | BIND ops, signing, validation, rollover, debugging |
| [[dnssec-kasp-policy]] | ISC KB — KASP ops | ✅ | `dnssec-policy` migration, rollover signals, KASP states |
| [[cs161-dnssec]] | Academic — CS161 UC Berkeley | ✅ | Query walkthrough, DNSKEY flags, OPT pseudosection, offline signing |
| [[rfc-6781]] | IETF RFC (2012) | ✅ | Rollover procedures, timing terminology, emergency rollover |
| RFC 9364 | Standard | ❌ | Reference |
| dig and delv | Tooling | ❌ | Debugging |
| delv manpage | Tooling | ❌ | Debugging |

---

## Open Questions

- What BIND version is the target environment (9.16 LTS vs 9.18 LTS vs 9.20)?
- DNSViz documentation — should it be added as a raw source?

---

## Contradictions Under Review

*None yet.*

---

*Updated on: 2026-04-14 | Sources ingested: 8 / 11*
