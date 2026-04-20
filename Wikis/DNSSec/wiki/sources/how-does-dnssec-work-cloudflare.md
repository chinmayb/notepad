---
title: "How Does DNSSEC Work? (Cloudflare)"
type: source
tags: ["dnssec", "cloudflare", "tutorial", "chain-of-trust", "rrset", "ksk", "zsk", "ds-record"]
created: 2026-04-13
updated: 2026-04-13
source_count: 1
confidence: high
status: active
related: ["[[dnssec-internals]]", "[[chain-of-trust]]", "[[dnssec-vs-dns-resolution]]", "[[dnssec-learning-roadmap]]"]
source_url: "https://www.cloudflare.com/en-gb/learning/dns/dnssec/how-dnssec-works/"
source_author: "Cloudflare"
source_date: 2026-04-13
source_type: article
open_questions:
  - "The CDS/CDNSKEY automation protocol (Internet Draft with CIRA) — was it ever standardized? Is it RFC 7344 / 8078?"
  - "Cloudflare mentions a 'smaller ZSK' being possible because of the KSK split — what are typical key sizes today?"
  - "What is Cloudflare's specific NSEC solution ('DNSSEC Done Right') referenced in the article?"
contradictions: []
entities_mentioned:
  - "[[Cloudflare]]"
concepts_mentioned:
  - "[[dnssec-internals]]"
  - "[[chain-of-trust]]"
  - "[[dnssec-vs-dns-resolution]]"
---

# How Does DNSSEC Work? (Cloudflare)

> **Author**: Cloudflare · **Type**: Tutorial article
> **Original**: [cloudflare.com/learning/dns/dnssec/how-dnssec-works/](https://www.cloudflare.com/en-gb/learning/dns/dnssec/how-dnssec-works/)

---

## Summary

Cloudflare's DNSSEC explainer is an operator-facing tutorial that builds up the full DNSSEC trust model layer by layer — from individual record signing through the complete chain of trust to the root. It is the most useful source in the current collection for understanding *why* the KSK/ZSK split exists and *exactly how* the DS record delegates trust between zones. The article ends with a practical note on the deployment friction of DS record management and Cloudflare's proposed automation.

The key framing insight: **DNSSEC is to DNS what a wax seal is to a letter** — it proves authenticity and detects tampering, but does not conceal the contents. This analogy is more precise than the common HTTPS comparison.

---

## Key Claims

1. **RRsets are the atomic unit of DNSSEC signing** — not individual records. All records of the same type at the same label are bundled into an RRset; the RRSIG covers the entire set. `confidence: high`

2. **The KSK signs only the DNSKEY RRset (including the ZSK public key)**; the ZSK signs all other RRsets. This two-level structure means ZSK rotation never requires touching the parent zone. `confidence: high`

3. **The DS record in the parent zone is a hash of the child zone's KSK** — not the ZSK. The resolver hashes the child's KSK and compares it to the parent's DS record to verify the chain link. `confidence: high`

4. **KSK rollover requires a TTL-aware two-phase process**: (1) add new DS record to parent, (2) wait for old DS TTL to expire, (3) remove old DS. Skipping step 2 causes a validation gap. `confidence: high`

5. **DNSSEC adds authentication, not confidentiality** — analogous to HTTPS in the "trust" dimension but fundamentally different: HTTPS encrypts traffic; DNSSEC only signs responses. `confidence: high`

6. **NSEC zone-walking is a real concern** — NSEC reveals the sorted list of all record names in a zone to anyone who queries sequentially. NSEC3 uses hashed names to mitigate this (but doesn't eliminate it). `confidence: high`

7. **The Root Signing Ceremony is the physical trust anchor** — a multi-party, publicly audited event that produces the RRSIG over the root DNSKEY RRset. This ceremony (not a DS record) is what bootstraps trust in the entire system. `confidence: high`

---

## Notable Quotes

> "Similar to HTTPS, DNSSEC adds a layer of security by enabling authenticated answers on top of an otherwise insecure protocol. Whereas HTTPS encrypts traffic so nobody on the wire can snoop on your Internet activities, DNSSEC merely signs responses so that forgeries are detectable."

> "Why do we use separate zone-signing keys and key-signing keys? [...] it's difficult to swap out an old or compromised KSK. Changing the ZSK, on the other hand, is much easier. This allows us to use a smaller ZSK without compromising the security of the server."

> "Note that any change in the KSK also requires a change in the parent zone's DS record. Changing the DS record is a multi-step process that can end up breaking the zone if it's performed incorrectly."

---

## Evidence Quality

- **Source type**: Operator documentation from a major DNS infrastructure company
- **Audience**: DNS operators and web developers
- **Bias**: Cloudflare has commercial interest in making DNSSEC appear simple to deploy via their platform. The article ends with a product pitch. Core technical content is unaffected by this.
- **Depth**: Conceptual/tutorial — correct and clear, but does not cover algorithm choices, TTL interactions, or BIND-specific configuration

---

## The Validation Flow (Intra-Zone)

From this source, the full intra-zone validation sequence a resolver performs:

```
1. Request target RRset (e.g., AAAA records)
   → Server returns: RRset + RRSIG (signed by private ZSK)

2. Request DNSKEY records
   → Server returns: public ZSK + public KSK + RRSIG over DNSKEY RRset (signed by private KSK)

3. Verify: use public KSK to validate RRSIG over DNSKEY RRset → confirms ZSK is trusted
4. Verify: use public ZSK to validate RRSIG over target RRset → confirms answer is authentic
```

The KSK itself is trusted because its hash matches the DS record in the parent zone (or, for root, because of the Root Signing Ceremony).

---

## DS Record Mechanics (Cross-Zone Trust)

The concrete trust-handoff between parent and child:

```
Zone operator action: hash(child KSK public key) → submit to parent as DS record
Resolver action: hash(child KSK received in answer) → compare to parent's DS record
If match: child KSK is trusted → child ZSK is trusted → child's RRsets are trusted
```

The DS record is itself signed (RRSIG) by the parent's ZSK — so validating the DS record follows the same chain.

---

## Entities Mentioned

- [[Cloudflare]] — author; also describes their automated DNSSEC deployment via registrar protocol
- CIRA (.ca registry) — co-author of the Cloudflare Internet Draft for automated DS management

---

## Concepts Covered

- [[dnssec-internals]] — RRsets, RRSIG, DNSKEY, DS, NSEC/NSEC3, CDNSKEY/CDS
- [[chain-of-trust]] — DS record mechanics, Root Signing Ceremony, how trust delegates top-down
- [[dnssec-vs-dns-resolution]] — what a resolver does differently under DNSSEC (validation steps)

---

*Source file: `raw/dns-sec-cloudflare&text/How Does DNSSEC Work?.md`*
