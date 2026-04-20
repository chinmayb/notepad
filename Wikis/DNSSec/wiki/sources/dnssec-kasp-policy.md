---
title: "DNSSEC Key and Signing Policy (ISC KB)"
type: source
tags: ["dnssec", "bind", "isc", "kasp", "dnssec-policy", "key-rollover", "key-states", "inline-signing", "migration"]
created: 2026-04-14
updated: 2026-04-14
source_count: 1
confidence: high
status: active
related: ["[[dnssec-in-bind]]", "[[dnssec-internals]]", "[[dnssec-debugging]]", "[[Internet-Systems-Consortium]]"]
open_questions:
  - "What is the exact behavior of `parental-agents` when the parent zone's DS state cannot be determined automatically?"
  - "How does the inline-signing upgrade path interact with zones already using dynamic updates + DNSSEC in BIND 9.18?"
contradictions: []
---

# Source: DNSSEC Key and Signing Policy (ISC KB)

> **Raw file**: `raw/dnssec0guide/DNSSEC Key and Signing Policy.md`
> **Source URL**: https://kb.isc.org/docs/dnssec-key-and-signing-policy
> **Author**: [[Internet-Systems-Consortium]]
> **Type**: ISC knowledge base article — operational guide for `dnssec-policy` configuration, migration, and key rollover
> **Coverage axes**: BIND ops (Axis 5 deep-dive), key rollover, migration from legacy `auto-dnssec`

---

## Summary

Focused ISC KB article on `dnssec-policy` (KASP). Covers: enabling DNSSEC on new zones, migrating existing zones from `auto-dnssec`, the key state machine, KSK/CSK rollover procedure with manual DS signaling, NSEC3 configuration, and the inline-signing vs dynamic-update file behavior difference across BIND versions. Complements the broad BIND guide ([[dnssec-guide-bind9]]) with migration-focused detail and the key state machine not covered there.

---

## Key Claims

### `dnssec-policy` Replaces `auto-dnssec` Permanently

1. **`auto-dnssec` is removed**, not deprecated — removed since BIND 9.19.16 development release. The wiki's prior language ("deprecated") should be read as "removed in new versions." Any production deployment on BIND 9.20+ cannot use `auto-dnssec`.

2. **`dnssec-policy default` built-in parameters** — explicit values now known:
   - `dnskey-ttl PT1H` (1 hour)
   - `signatures-refresh P5D` (5 days)
   - `signatures-validity P14D` (14 days)
   - `max-zone-ttl P1D` (1 day)
   - `parent-ds-ttl P1D` (1 day)
   - `purge-keys P90D` (90 days after key removal)
   - Key: single CSK, algorithm ECDSAP256SHA256, `lifetime unlimited`

### BIND Version Gate for Inline Signing

3. **BIND 9.18 requires explicit `inline-signing yes`** — In 9.18, inline-signing is not automatic; it must be added to zone config or globally. In **BIND 9.20+**, `inline-signing` defaults to `yes` — `dnssec-policy` alone is sufficient with no extra config.

4. **9.20 upgrade pitfall** — If zones were signed with dynamic updates (no inline-signing) on 9.18, upgrading to 9.20 without adding `inline-signing no;` causes BIND to sign the already-signed zone file again, producing double-signed `...signed` files. Three mitigations: (a) revert to insecure first, (b) add `inline-signing no;` per zone before upgrade, (c) add `inline-signing no;` to the `dnssec-policy` block.

### KASP Key State Machine

5. **Four KASP key states** (distinct from key lifecycle timestamps Created/Publish/Activate/Inactive/Delete):

   | State | Meaning |
   |-------|---------|
   | `hidden` | Key not yet published; no resolver knows about it |
   | `rumoured` | Key published but propagation not complete; some resolvers may have it |
   | `omnipresent` | Key fully propagated; all resolvers either have it cached or will fetch it fresh |
   | `unretentive` | Key being retired; transitioning out of use |

   These states are tracked per resource record type (DNSKEY state, DS state, zone RRSIG state, KSK RRSIG state) in `.state` files alongside the key files.

6. **`.state` file** — BIND creates `Kzone.+alg+id.state` for each key. Contains: algorithm, length, lifetime, KSK/ZSK role, all timing timestamps, plus per-RR-type state (`DNSKEYState`, `DSState`, `ZRRSIGState`, `KRRSIGState`) and a `GoalState`.

### KSK/CSK Rollover — The Two Manual Signals

7. **KSK/CSK rollovers cannot be fully automated without `parental-agents`** — BIND pauses and waits for operator signals at two points:
   - After new DS is published in parent: `rndc dnssec -checkds -key <new-id> published <zone>`
   - After old DS is withdrawn from parent: `rndc dnssec -checkds -key <old-id> withdrawn <zone>`
   
   Without these signals, BIND holds the rollover in `rumoured` DS state indefinitely. `parental-agents { <ip>; }` automates this by letting BIND query the parent zone directly.

8. **DS readiness check** — Before signaling `published`, verify the new DS is actually in the parent by querying a parent zone server:
   ```
   dig +dnssec +norecurse +multi @<parent-ns-ip> <zone> DS
   ```
   The RRSIG on the DS RR will show the new key ID when the rollover is complete.

### ZSK vs KSK Rollover Strategies

9. **ZSK rollover = Pre-Publish strategy** (fully automatic): New DNSKEY is prepublished; after propagation, RRSIGs are gradually replaced. No operator action needed.

10. **KSK rollover = Double-KSK strategy** (requires operator DS signals): New DNSKEY prepublished next to existing KSK; after propagation, DS can be submitted to parent. Defined in RFC 7583.

11. **CSK rollover = Hybrid** (Pre-Publish + Double-KSK combined): New DNSKEY prepublished; zone RRSIGs gradually replaced; simultaneously DS can be submitted. Most complex but most common with `default` policy.

### Migration Prerequisites

12. **Private-key-format v1.3 required** — Existing keys must be at format v1.3 before `dnssec-policy` can use them. Older format keys have no metadata and are incompatible. Migrate with: `dnssec-settime -f <keyfile>`. Check: look for `Private-key-format: v1.3` in the `.private` file.

13. **Match existing keys in a custom policy first** — Switching a zone from `auto-dnssec` to `dnssec-policy default` triggers an immediate rollover if existing keys don't match (different algorithm or KSK/ZSK vs CSK). To migrate without rolling, create a custom `dnssec-policy` that exactly matches the current algorithm and key structure.

14. **`DSState: rumoured` after migration** — If migrating an existing zone whose DS is already published, BIND doesn't know that. Run `rndc dnssec -checkds published <zone>` to unstick it. Failing to do so leaves BIND waiting indefinitely.

---

## Direct Quotes Worth Preserving

> "Rollover correctness is guaranteed by DNSSEC records state machines. BIND 9 tries to follow the timings from the policy, but does not apply them if they would result in the zone becoming bogus."

> "`default` is a built-in policy that introduces a single ECDSAP256SHA key and accompanying signatures for the zone [...] Additionally it uses NSEC (not NSEC3 so far) for denial of existence."

> "Do not use extra iterations, salt, or opt-out unless their implications are fully understood. See RFC 9276."

---

## Entities and Concepts Mentioned

**Entities**: [[Internet-Systems-Consortium]]
**Concepts**: [[dnssec-in-bind]] (primary — `dnssec-policy` details, version gates, migration, KSK rollover signals), [[dnssec-internals]] (KASP key states, `.state` file format, CSK rollover strategy)

---

## Evidence Quality Notes

- **High authority**: Official ISC knowledge base article — written by BIND maintainers
- **Version-specific**: References BIND 9.16 through 9.20; inline-signing behavior explicitly version-gated
- **Complements**: [[dnssec-guide-bind9]] (broad guide) — this article is migration + rollover operations focused
