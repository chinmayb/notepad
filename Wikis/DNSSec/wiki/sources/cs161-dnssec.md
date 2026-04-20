---
title: "DNSSEC — CS161 UC Berkeley Textbook"
type: source
tags: ["dnssec", "pedagogy", "chain-of-trust", "dig", "query-walkthrough", "nsec", "nsec3", "ksk", "zsk", "trust-delegation"]
created: 2026-04-14
updated: 2026-04-14
source_count: 1
confidence: high
status: active
related: ["[[dnssec-vs-dns-resolution]]", "[[chain-of-trust]]", "[[dnssec-internals]]", "[[dnssec-debugging]]"]
open_questions: []
contradictions: []
---

# Source: DNSSEC — CS161 UC Berkeley Textbook

> **Raw file**: `raw/dns-sec-cloudflare&text/DNSSEC.md`
> **Source URL**: https://textbook.cs161.org/network/dnssec.html
> **Author**: CS 161: Computer Security, UC Berkeley
> **Type**: Academic textbook chapter — pedagogical, step-by-step walkthrough with `dig` output
> **Coverage axes**: Chain of trust (Axis 2), query resolution under DNSSEC (Axis 4), NSEC/NSEC3 fundamentals (Axis 1)

---

## Summary

A compact, pedagogically focused chapter from the UC Berkeley CS161 computer security course. Builds DNSSEC from first principles (why signatures aren't enough without a trust anchor), walks through the full chain-of-trust delegation model, and shows a complete `dig`-based query for `eecs.berkeley.edu` at every level. Key value: the most concrete step-by-step query walkthrough of any source in the wiki, including actual `dig` output showing DNSKEY flag values and DS records at each delegation.

---

## Key Claims

### Trust Model

1. **Trust delegation** — The root acts as trust anchor. Trust flows downward: root signs `.edu`'s KSK → `.edu` signs `berkeley.edu`'s KSK → `berkeley.edu` signs its own records. Each endorsement is a signature on the child's public key.

2. **KSK as "parent within the name server"** — An elegant pedagogical framing: the KSK is the "parent" and the ZSK is the "child" within a single name server. The KSK endorses the ZSK (signs it), just as a parent zone endorses a child zone.

### DNSKEY Record Flag Values

3. **Flag 256 = ZSK; Flag 257 = KSK** — Explicitly visible in `dig DNSKEY` output. The SEP (Secure Entry Point) bit is set on KSKs (bit 15 of the flags field → 257 vs 256).

### OPT Pseudosection Explains "Extra" ADDITIONAL Record

4. **ADDITIONAL count = actual records + 1** — The `OPT` pseudosection (which carries the `DO` bit for DNSSEC and the EDNS0 buffer size) is encoded as an additional record for backwards compatibility. This is why `dig` always shows ADDITIONAL count as one higher than expected.

### Full `dig` Query Walkthrough (eecs.berkeley.edu)

5. **Step 1 — Root DNSKEY**: Query root for `DNSKEY .` → returns ZSK (flag 256), KSK (flag 257), RRSIG over DNSKEY. Use hardcoded root KSK to verify RRSIG → trust root ZSK.

6. **Step 2 — Root for eecs.berkeley.edu**: Query root → returns NS records for `.edu` + **DS record** (hash of `.edu`'s KSK) + RRSIG on DS (signed by root ZSK). Root ZSK signed the DS → trust `.edu`'s KSK.

7. **Step 3 — `.edu` DNSKEY**: Query `.edu` for `DNSKEY edu.` → returns `.edu` ZSK + KSK + RRSIG. `.edu` KSK verifies RRSIG → trust `.edu` ZSK.

8. **Step 4 — `.edu` for eecs.berkeley.edu**: Query `.edu` → returns NS + A records + **DS record** (hash of `berkeley.edu`'s KSK) + RRSIG on DS (signed by `.edu` ZSK). Trust `berkeley.edu`'s KSK.

9. **Step 5 — berkeley.edu DNSKEY**: Query `berkeley.edu` for `DNSKEY berkeley.edu.` → returns `berkeley.edu` ZSK + KSK + RRSIG. Trust `berkeley.edu` ZSK.

10. **Step 6 — Final answer**: Query `berkeley.edu` for `eecs.berkeley.edu A` → returns A record (23.185.0.1) + RRSIG over A record (signed by `berkeley.edu` ZSK). Trust the final answer.

### Why Offline Signing (NSEC Motivation)

11. **Online signing = DoS vulnerability** — If signatures were computed on-demand per query, an attacker could flood the server with NXDOMAIN queries for nonexistent domains and exhaust CPU. Offline signing prevents this.

12. **NSEC pre-computed ranges** — Instead of signing individual NXDOMAIN responses, zones pre-compute signed ranges: "nothing exists between `b.example.com` and `l.example.com`." This authenticates non-existence without online signing.

13. **NSEC zone walking** — By iterating queries, an attacker can walk the entire zone (query `c`, learn `b` and `l` exist; query `m`, learn `q` exists; etc.). The NSEC sorted linked list structure enables complete zone enumeration.

14. **NSEC3 hashes names** — Uses hashes of domain names in the sorted list. Prevents trivial zone walking but doesn't prevent GPU-powered pre-computation. NSEC5 (RFC draft) is the next attempt.

---

## Direct Quotes Worth Preserving

> "Because the signatures, a network attacker (MITM, on-path, off-path) cannot tamper with the data or inject malicious data without being detected."

> "You can think of the KSK as the 'parent' and the ZSK as the 'child,' both within one name server."

> "An attacker could buy a GPU and precompute hashes to learn domain names anyway…"

> "If we don't trust anybody, then DNSSEC will never work (we'll never trust any records we get), so we must first choose a trust anchor."

---

## Entities and Concepts Mentioned

**Entities**: none new
**Concepts**: [[dnssec-vs-dns-resolution]] (primary — fills this stub with full query walkthrough), [[chain-of-trust]] (trust delegation model, DNSKEY flags), [[dnssec-internals]] (NSEC zone walking, NSEC3 hashing, offline signing motivation, OPT record / DO bit)

---

## Evidence Quality Notes

- **Academic source**: UC Berkeley CS161 security course — pedagogically optimized, not an operational reference
- **High clarity**: Best source in the wiki for explaining the query walkthrough and KSK/ZSK motivation intuitively
- **Not operational**: No BIND config, no troubleshooting, no KASP. Complements rather than overlaps with operational sources.
