---
title: "DNSSEC vs Normal DNS Resolution"
type: concept
tags: ["dnssec", "dns", "resolution", "workflow", "dig", "dnskey", "rrsig", "ds", "do-bit", "opt-record", "nsec", "chain-of-trust", "trust-anchor"]
created: 2026-04-13
updated: 2026-04-14
source_count: 2
confidence: high
status: active
related: ["[[dnssec-internals]]", "[[chain-of-trust]]", "[[dnssec-debugging]]", "[[dnssec-in-bind]]"]
open_questions:
  - "Does a validating resolver re-validate cached responses, or only validate on initial fetch?"
  - "How does DNSSEC interact with DNS over HTTPS (DoH) / DNS over TLS (DoT)?"
  - "What is the performance overhead of DNSSEC validation in practice?"
contradictions: []
---

# DNSSEC vs Normal DNS Resolution

> **Domain**: DNS / DNSSEC
> **Also known as**: DNSSEC resolution workflow, validating resolver behavior

---

## Definition

When DNSSEC validation is enabled, a recursive resolver performs additional cryptographic verification steps at each delegation boundary in the resolution chain. This page captures *exactly what changes* compared to plain DNS resolution — walking through a concrete `dig`-based example at every step.

---

## What Stays the Same

- The iterative resolution process: resolver → root → TLD → authoritative
- The record types being requested (A, AAAA, MX, etc.)
- UDP/TCP transport, port 53
- Caching of responses (TTL behavior)
- NXDOMAIN for non-existent domains (but now authenticated)

---

## What Changes with DNSSEC

At each delegation step, the resolver additionally:
1. Sets the **DO (DNSSEC OK)** bit in its EDNS0 OPT record — signals that it wants DNSSEC records returned
2. Requests the `DNSKEY` records for the zone
3. Verifies the `RRSIG` over the `DNSKEY` RRset using the KSK
4. Verifies the `DS` record in the parent zone matches the child's KSK (hashed)
5. Verifies the `RRSIG` over the actual answer using the ZSK

### The DO Bit and the OPT Record

The **DO (DNSSEC OK)** bit travels inside an EDNS0 **OPT pseudo-record**. This OPT record is encoded as an "additional" record for backwards compatibility — meaning:

> **ADDITIONAL count in `dig` output = actual additional records + 1**

The +1 is always the OPT pseudo-record carrying the DO bit and the EDNS0 buffer size. This is not an error — it is by design. `dig` shows it in the `; OPT PSEUDOSECTION:` block.

```
;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags: do; udp: 1232
```

`flags: do` = the DO bit is set. `udp: 1232` = EDNS0 payload size announced by the querier.

Without the DO bit (plain `dig` without `+dnssec`), authoritative servers suppress RRSIG and DNSKEY records entirely — backwards compatible with pre-DNSSEC resolvers.

### On Validation Failure

A validating resolver returns **SERVFAIL** — it does NOT return the unvalidated answer. This is a hard stop: the client gets no answer at all.

> This is why DNSSEC misconfiguration causes total outages, not partial degradation. See [[dnssec-debugging]] for the triage methodology.

---

## Full DNSSEC Resolution Walk: `eecs.berkeley.edu`

Source: [[cs161-dnssec]]. This is the complete resolver workflow, step by step, with the actual records returned at each stage.

### Trust Anchor Foundation

Before any query, the resolver has one hardcoded fact: **the root KSK public key**. This is the trust anchor. It cannot be obtained from DNS — it is baked into the resolver software. Everything else is derived from it.

> "If we don't trust anybody, then DNSSEC will never work (we'll never trust any records we get), so we must first choose a trust anchor." — CS161

---

### Step 1 — Fetch and Verify the Root Zone's Keys

**Query**: `dig DNSKEY .` (root zone DNSKEY)

**Returns**:
- `DNSKEY` flag **256** (ZSK) — the root's zone signing key
- `DNSKEY` flag **257** (KSK) — the root's key signing key
- `RRSIG` over the DNSKEY RRset (signed by the root's private KSK)

**Verification**:
- Use the **hardcoded root KSK** (trust anchor) to verify the RRSIG over the DNSKEY RRset
- This confirms: the root ZSK (flag 256) is authentic

**Result**: Root ZSK is trusted.

---

### Step 2 — Root Delegates to `.edu`

**Query**: Root nameservers for `eecs.berkeley.edu A`

**Returns**:
- NS records pointing to `.edu` nameservers
- `DS` record = hash of `.edu`'s KSK
- `RRSIG` over the DS record, signed by the **root ZSK**

**Verification**:
- Use the trusted root ZSK (verified in Step 1) to verify the RRSIG over the DS record
- The DS is now trusted → `.edu`'s KSK is authentic (anyone whose key hashes to that DS is `.edu`'s legitimate KSK)

**Result**: `.edu`'s KSK is trusted via root's endorsement.

---

### Step 3 — Fetch and Verify `.edu` Zone's Keys

**Query**: `dig DNSKEY edu.` (`.edu` nameservers)

**Returns**:
- `DNSKEY` flag **256** (ZSK) — `.edu`'s zone signing key
- `DNSKEY` flag **257** (KSK) — `.edu`'s key signing key
- `RRSIG` over the DNSKEY RRset, signed by `.edu`'s private KSK

**Verification**:
- Use the trusted `.edu` KSK (verified via DS in Step 2) to verify the RRSIG over the DNSKEY RRset
- This confirms: `.edu`'s ZSK (flag 256) is authentic

**Result**: `.edu` ZSK is trusted.

---

### Step 4 — `.edu` Delegates to `berkeley.edu`

**Query**: `.edu` nameservers for `eecs.berkeley.edu A`

**Returns**:
- NS records + A records (glue) pointing to `berkeley.edu` nameservers
- `DS` record = hash of `berkeley.edu`'s KSK
- `RRSIG` over the DS record, signed by the **`.edu` ZSK**

**Verification**:
- Use the trusted `.edu` ZSK (verified in Step 3) to verify the RRSIG over the DS record
- The DS is now trusted → `berkeley.edu`'s KSK is authentic

**Result**: `berkeley.edu`'s KSK is trusted via `.edu`'s endorsement.

---

### Step 5 — Fetch and Verify `berkeley.edu` Zone's Keys

**Query**: `dig DNSKEY berkeley.edu.`

**Returns**:
- `DNSKEY` flag **256** (ZSK) — `berkeley.edu`'s zone signing key
- `DNSKEY` flag **257** (KSK) — `berkeley.edu`'s key signing key
- `RRSIG` over the DNSKEY RRset, signed by `berkeley.edu`'s private KSK

**Verification**:
- Use the trusted `berkeley.edu` KSK (verified via DS in Step 4) to verify the RRSIG over the DNSKEY RRset
- This confirms: `berkeley.edu`'s ZSK is authentic

**Result**: `berkeley.edu` ZSK is trusted.

---

### Step 6 — Final Answer: `eecs.berkeley.edu A`

**Query**: `dig eecs.berkeley.edu A` (to `berkeley.edu` nameservers)

**Returns**:
- `A` record: `23.185.0.1`
- `RRSIG` over the A RRset, signed by `berkeley.edu`'s private ZSK

**Verification**:
- Use the trusted `berkeley.edu` ZSK (verified in Step 5) to verify the RRSIG over the A RRset
- Answer is authentic

**Result**: Resolver sets the **AD flag** in the response. DNSSEC validation complete.

---

### The Chain Visualized

```
Hardcoded root KSK (trust anchor)
  └─ verifies → RRSIG(root DNSKEY RRset) → root ZSK trusted
       └─ verifies → RRSIG(DS for .edu) → .edu KSK trusted
            └─ verifies → RRSIG(.edu DNSKEY RRset) → .edu ZSK trusted
                 └─ verifies → RRSIG(DS for berkeley.edu) → berkeley.edu KSK trusted
                      └─ verifies → RRSIG(berkeley.edu DNSKEY RRset) → berkeley.edu ZSK trusted
                           └─ verifies → RRSIG(eecs.berkeley.edu A) → answer trusted ✓
```

Each arrow is a cryptographic signature check. A single broken link → SERVFAIL.

---

## Reading DNSKEY Flag Values

Two flag values appear consistently in `dig DNSKEY` output:

| Flag value | Key type | Role |
|-----------|---------|------|
| **256** | ZSK (Zone Signing Key) | Signs all RRsets in the zone |
| **257** | KSK (Key Signing Key) | Signs only the DNSKEY RRset; its hash is the DS record in the parent |

The difference is the **SEP bit** (Secure Entry Point, bit 15 of the flags field). KSKs have it set (256 + 1 = 257). This is visible in raw `dig` output as the flags integer and in the `dig` type column (`DNSKEY  256 3 13 <key>` vs `DNSKEY  257 3 13 <key>`).

Source: [[cs161-dnssec]], [[dnssec-internals]]

---

## Why Offline Signing: The NSEC/NXDOMAIN Problem

DNSSEC must authenticate not just positive answers but also **negative answers** — proof that a name does not exist. Two approaches are possible:

### Option A: Online Signing (rejected)
Sign each NXDOMAIN response on demand when queried.

**Problem**: An attacker can flood a nameserver with NXDOMAIN queries for nonexistent names, exhausting CPU with on-demand signing operations → DoS.

### Option B: Offline Signing (adopted)
Pre-compute signed proofs of non-existence at zone signing time. Serve them directly from the zone file.

**NSEC mechanism**: The zone maintains a **sorted linked list** of all its record names. For any non-existent name, the server returns the NSEC record for the adjacent names — "nothing exists between `b.example.com` and `l.example.com`" — with a valid RRSIG already attached.

**NSEC3**: Same mechanism, but uses hashed domain names in the sorted list to prevent trivial zone enumeration (zone walking). See [[dnssec-internals]] for NSEC vs NSEC3 decision guidance.

Sources: [[cs161-dnssec]], [[how-does-dnssec-work-cloudflare]]

---

## Caching Under DNSSEC

- RRSIG records are cached alongside the RRsets they sign
- TTL is the minimum of the RRset TTL and the RRSIG expiration — whichever is shorter
- Resolvers cache DNSKEY records at their TTL — this is why key rollover requires TTL-expiry wait periods before old keys can be removed
- RFC 9077 updates how NSEC/NSEC3 TTLs interact with caching

---

## Mentioned In

- [[cs161-dnssec]] — primary source: 6-step `dig` walkthrough, DNSKEY flag values 256/257, OPT pseudosection explanation, offline signing motivation, NSEC zone walking
- [[how-does-dnssec-work-cloudflare]] — DO bit, RRSIG/DS mechanics, caching behavior
