---
title: "Synthesis: DNSSEC Learning Roadmap"
type: synthesis
tags: ["dnssec", "learning-roadmap", "curriculum"]
created: 2026-04-13
updated: 2026-04-20
source_count: 8
confidence: medium
status: active
related: ["[[01-learn-dnssec]]", "[[dnssec-internals]]", "[[chain-of-trust]]", "[[dnssec-vs-dns-resolution]]", "[[dnssec-in-bind]]", "[[dnssec-debugging]]"]
open_questions:
  - "What BIND version is the target (BIND 9.16 LTS vs 9.18 LTS vs 9.20 current)?"
  - "Is the goal authoritative server signing, validating resolver config, or both?"
  - "Are there planned sources for the Kaminsky attack and zone walking / NSEC5?"
contradictions: []
---

# Synthesis: DNSSEC Learning Roadmap

> **Purpose**: Curriculum map — what to read in what order, what each source covers, what gaps remain
> **Sources synthesized**: 8 sources (as of 2026-04-20)
> **Derived from**: [[01-learn-dnssec]] (goals) + gap analysis across all raw/ sources

---

## Bottom Line

The current source collection covers **DNSSEC theory and BIND operations well** (8 of 11 sources ingested). The conceptual foundation (chain of trust, record types, signing mechanics) and the full BIND operational stack (`dnssec-policy`, key lifecycle, rollover, validation) are now documented. Debugging remains the weakest axis — the `dig`/`delv` tool sources and DNSViz documentation have not yet been ingested (Phase 3 remaining).

**Suggested reading/ingestion order below** — optimized for building understanding layer by layer, not source-by-source.

---

## Phase 1 — Foundation: Why and What (Axes 1–3)

*Build the mental model before touching any config.*

### Step 1.1 — The threat: why DNS needed securing
**Ingest next**: `raw/dnssec0guide/DNS Vulnerabilities Lessons from the Kaminsky Exploit...md`
- Establishes the Kaminsky 2008 attack as the canonical motivation
- Without this, DNSSEC feels like complexity for its own sake
- Covers: cache poisoning, birthday attacks, port randomization as mitigation, why that wasn't enough

### Step 1.2 — The mechanism: theory-first walkthrough
**Already reviewed**: `raw/dns-sec-cloudflare&text/DNSSEC.md` (CS161 Berkeley)
- Best single source for the conceptual model
- Chain of trust, KSK/ZSK, NSEC/NSEC3, full `dig` walkthrough with real output
- Read this before the RFC

### Step 1.3 — The mechanism: operator's perspective
**Already reviewed**: `raw/dns-sec-cloudflare&text/How Does DNSSEC Work?.md` (Cloudflare)
- Complements CS161 — focuses on RRsets, DS record mechanics, registrar DS submission
- Good for understanding why key rollovers are operationally painful

### Step 1.4 — The standard: what the IETF says
**Already reviewed**: `raw/dns-sec-rfc/RFC 9364 DNS Security Extensions (DNSSEC).md`
- Not for learning — for reference. Lists all normative RFCs and their scope.
- Use as a map to know which RFC to consult for specific questions
- RFC 6781 (operational practices) is the most useful follow-on from here

---

## Phase 2 — Operations: How to Deploy (Axes 4–5)

*Once the model is clear, learn how to actually run it.*

### Step 2.1 — BIND-specific overview
**Ingest next**: `raw/dnssec0guide/DNSSEC and BIND 9.md`
- BIND-specific orientation: how BIND implements signing, key management, validation
- Establishes BIND 9 terminology that the full guide uses

### Step 2.2 — Full BIND guide
**Ingest next**: `raw/dnssec0guide/DNSSEC Guide — BIND 9 9.21.22-dev documentation.md`
- The authoritative operational reference for BIND DNSSEC
- Covers: inline signing, `dnssec-policy`, key generation, zone signing, rollovers
- This is the most important practical source in the collection

### Step 2.3 — Key and Signing Policy deep dive
**Ingest next**: `raw/dnssec0guide/DNSSEC Key and Signing Policy.md`
- BIND 9.16+ `dnssec-policy` directive — the modern automated approach
- Replaces manual `dnssec-keygen` + `dnssec-signzone` workflow
- Critical for understanding how production BIND setups work today

### Step 2.4 — Operational practices (RFC)
**Ingest next**: `raw/dns-sec-rfc/RFC 6781 DNSSEC Operational Practices, Version 2.md`
- The IETF operational playbook — key rollover timing, TTL interactions, algorithm choices
- Use after the BIND guide — fills in the "why these settings" questions the BIND docs don't answer

---

## Phase 3 — Debugging: How to Fix It (Axis 6)

*The most perishable knowledge — best learned hands-on, but wiki provides the reference.*

### Step 3.1 — Tool fundamentals
**Ingest next**: `raw/delv/dig and delv.md`
- `dig +dnssec` vs `delv` — when to use which
- `delv` is the right tool for DNSSEC validation debugging

### Step 3.2 — `delv` reference
**Ingest next**: `raw/delv/Ubuntu Manpage — delv.md`
- Full flag reference
- Build a cheatsheet wiki page from this

### Step 3.3 — Visual debugging (gap — source not yet in raw/)
**Source needed**: dnsviz.net documentation or an explanatory article
- DNSViz renders the full chain of trust visually — invaluable for "where is the chain broken?"
- Add to raw/ before ingesting Phase 3

---

## Current Source Inventory

| Raw Source | Axis Coverage | Ingested | Priority |
|-----------|--------------|----------|----------|
| `01-learn-dnssec.md` | Goals anchor | ✅ 2026-04-13 | Done |
| `DNSSEC.md` (CS161) | Axes 1, 2, 3 | ✅ 2026-04-14 | Done |
| `How Does DNSSEC Work?.md` (Cloudflare) | Axes 1, 3, 4 | ✅ 2026-04-13 | Done |
| `DNS Vulnerabilities / Kaminsky` | Axis 3 | ✅ 2026-04-13 | Done |
| `DNSSEC and BIND 9.md` | Axes 4, 5 | ✅ 2026-04-13 | Done |
| `DNSSEC Guide — BIND 9.md` | Axes 4, 5 | ✅ 2026-04-13 | Done |
| `DNSSEC Key and Signing Policy.md` | Axis 5 | ✅ 2026-04-14 | Done |
| `RFC 6781` | Axis 4 | ✅ 2026-04-14 | Done |
| `RFC 9364` | Reference | ❌ | After Phase 3 |
| `dig and delv.md` | Axis 6 | ❌ | Phase 3 — next |
| `Ubuntu Manpage delv.md` | Axis 6 | ❌ | Phase 3 |

---

## Concept Pages to Build (in order)

As sources are ingested, these concept pages should be created and progressively filled:

1. [[chain-of-trust]] — root of all DNSSEC understanding
2. [[dnssec-internals]] — record types, signing mechanics
3. [[dnssec-vs-dns-resolution]] — the workflow delta
4. [[dnssec-in-bind]] — BIND-specific config and operation
5. [[dnssec-debugging]] — tools, failure modes, triage methodology
6. [[key-rollover]] — complex enough to deserve its own page
7. [[nsec-and-nsec3]] — denial of existence — nuanced enough to stand alone
8. [[kaminsky-attack]] — the historical motivation, entity page for the attack

---

## Known Gaps in Raw Collection

| Gap | Impact | Recommended Source |
|-----|--------|--------------------|
| DNSViz documentation / tutorial | High — visual debugging is the fastest way to diagnose broken chains | https://dnsviz.net or blog.apnic.net |
| SERVFAIL triage guide | High — what to check when a client gets SERVFAIL due to DNSSEC | RFC 8027 (Roadblock Avoidance) or ISC KB |
| Key rollover walkthrough (practical) | Medium — RFC 6781 covers timing but not the BIND-specific commands | ISC BIND Knowledge Base |
| DANE/TLSA | Low — shows downstream value of DNSSEC | RFC 6698 summary article |

---

*Derived from [[01-learn-dnssec]]. Update this page after each Phase is complete.*
