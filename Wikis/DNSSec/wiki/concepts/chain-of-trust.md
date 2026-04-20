---
title: "Chain of Trust"
type: concept
tags: ["dnssec", "trust", "pki", "chain-of-trust", "ds-record", "ksk-rollover", "kaminsky", "trust-anchor", "validation", "emergency-rollover", "stand-by-keys", "ksk", "zsk", "pedagogy"]
created: 2026-04-13
updated: 2026-04-14
source_count: 5
confidence: high
status: active
related: ["[[dnssec-internals]]", "[[dnssec-vs-dns-resolution]]", "[[dnssec-in-bind]]", "[[dnssec-debugging]]", "[[dns-cache-poisoning]]"]
open_questions:
  - "How was the root KSK trust anchor rolled in 2018, and what went wrong for some resolvers?"
  - "How does the CDS/CDNSKEY automation (RFC 7344/8078) interact with the DS update TTL requirement?"
contradictions: []
---

# Chain of Trust

> **Domain**: DNSSEC / PKI
> **Also known as**: DNSSEC trust chain, delegation of trust

---

## Definition

The DNSSEC chain of trust is the cryptographic hierarchy that lets a validating resolver verify that a DNS response is authentic, by building a verified path from the **root zone's trust anchor** (hardcoded in the resolver) all the way down to the final answer.

It is the central organizing principle of DNSSEC. Every other DNSSEC mechanism — KSK/ZSK split, DS records, RRSIGs, NSEC — exists to establish, maintain, or verify this chain.

---

## Why This Chain Had to Be Built: The First-Principles Motivation

DNS was designed in the 1980s without authentication. Any server could respond to a query claiming to be authoritative, and resolvers would accept the answer. This became critically exploitable in 2008.

**The [[dns-cache-poisoning]] problem**: a recursive resolver caches answers to serve thousands of clients. If an attacker can inject a forged answer into that cache — by guessing the 16-bit transaction ID of an outbound query — every client that resolver serves gets redirected silently to an attacker-controlled IP.

[[Dan-Kaminsky]] demonstrated in 2008 that this was not just theoretically possible but *practically easy*: by flooding the resolver with queries for uncached subdomains, an attacker could generate thousands of transaction ID guessing attempts per second. The emergency fix (source port randomization) expanded the guessable space from 16 bits to ~32 bits — harder, but not cryptographically sound.

**The structural answer**: if every DNS response carries a cryptographic signature from the zone that originated it, an attacker who forges a response cannot produce a valid signature without the zone's private key. The forged answer is rejected regardless of whether the transaction ID matches. This is what the chain of trust implements.

> The chain of trust is DNSSEC's answer to the question: "How do we know this response really came from the zone that owns this name?"

Sources: [[kaminsky-exploit]], [[how-does-dnssec-work-cloudflare]]

---

## The Trust Anchor (Root)

- The root zone's KSK public key is **hardcoded** into all validating resolvers — this is the trust anchor
- It cannot be overridden by any DNS response — only by updating the resolver software/configuration
- Updated via RFC 5011 (automated trust anchor rollover) or manual configuration
- The **Root Signing Ceremony** is the physical, multi-party, publicly audited event that generates and controls the root private KSK. Several individuals worldwide each hold a key share — no single person can sign the root alone.

**BIND trust anchor configuration** (from [[dnssec-guide-bind9]]):

| Config value | Behavior |
|---|---|
| `dnssec-validation auto` | BIND ships with root trust anchor; updates it automatically via RFC 5011 (`managed-keys`). **Default in BIND 9.18+.** |
| `dnssec-validation yes` | Requires operator to configure trust anchors manually in `named.conf` |
| `dnssec-validation no` | Validation disabled entirely; not recommended |

`auto` is correct for almost all deployments. The root trust anchor is maintained by ISC and distributed with BIND.

Sources: [[how-does-dnssec-work-cloudflare]], [[dnssec-guide-bind9]]

---

## How Trust Delegates Downward

Trust flows exclusively **top-down** — from root → TLD → zone. A child zone cannot self-certify; it must be endorsed by its parent.

```
Root KSK (hardcoded trust anchor in resolver)
  └─ signs → Root ZSK (via RRSIG over root DNSKEY RRset)
       └─ signs → .com DS record (hash of .com KSK, RRSIG'd by root ZSK)
            └─ .com KSK (verified via DS)
                 └─ signs → .com ZSK
                      └─ signs → example.com DS (hash of example.com KSK)
                           └─ example.com KSK (verified via DS)
                                └─ signs → example.com ZSK
                                     └─ signs → A/AAAA/MX/etc. RRsets
```

Each `└─ signs →` is a verifiable cryptographic signature. A resolver walks this chain and verifies every link. If any link fails — the chain breaks and the resolver returns SERVFAIL.

---

## The 12-Step Validation Walkthrough

This is how a resolver validates a query for `www.example.com A`, step by step. Each step is a DNS query + verification. Source: [[dnssec-guide-bind9]].

```
 1. Resolver queries root servers for example.com A
 2. Root returns referral to .com + DS record for .com (signed by root ZSK)
 3. Resolver verifies root DS with trust anchor → root chain confirmed
 4. Resolver queries .com servers for example.com A
 5. .com returns referral to example.com + DS record for example.com (signed by .com ZSK)
 6. Resolver fetches .com DNSKEY, verifies DS match → .com chain confirmed
 7. Resolver queries example.com servers for www.example.com A
 8. example.com returns A record + RRSIG (signed by example.com ZSK)
 9. example.com also returns DNSKEY RRset + RRSIG over DNSKEY (signed by example.com KSK)
10. Resolver hashes example.com KSK → compares against .com DS record → KSK confirmed
11. Resolver uses example.com KSK to verify RRSIG over DNSKEY RRset → ZSK confirmed
12. Resolver uses example.com ZSK to verify RRSIG over A RRset → answer confirmed
    → Sets AD flag in response. DNSSEC validation complete.
```

**The AD flag** in `dig` output (`flags: qr rd ra ad`) is the resolver's assertion that every link in this chain verified successfully. No `ad` flag = answer not validated.

**Caching optimizes this**: resolvers cache DNSKEY records at their TTL. Steps 9–11 are only performed on first query per TTL window. This is also why **rollover complexity exists**: cached DNSKEY records must expire before a new key is trusted.

---

## The DS Record: The Concrete Cross-Zone Link

The DS (Delegation Signer) record is how a parent zone endorses a child zone:

```
Zone operator: hash( child zone's KSK public key ) → submit to parent registrar as DS record
Parent zone: publishes DS record, signs it with parent's ZSK
Resolver: fetches child KSK, hashes it, compares to parent's DS record
           → if they match: child KSK is authentic → chain continues
           → if they don't match: SERVFAIL
```

**Key property**: the DS record lives in the **parent** zone, not the child zone. The child cannot forge or modify it — the parent controls it. This is what makes the top-down endorsement meaningful.

Source: [[how-does-dnssec-work-cloudflare]]

---

## KSK Rollover: The Operationally Expensive Operation

Changing a zone's KSK is the most dangerous operation in DNSSEC. The DS record in the parent must be updated, but there's a timing constraint:

**Why rollover is complex**: recursive resolvers cache DNSKEY records at their TTL. If you swap a KSK instantly, any resolver with the old DNSKEY cached will fail to verify new signatures → SERVFAIL. All rollover procedures insert TTL-expiry wait periods to safely drain caches.

**Three rollover methods** (from [[dnssec-guide-bind9]]):

| Method | Parent interactions | DNSKEY RRset size during transition | Notes |
|--------|--------------------|------------------------------------|-------|
| **Double-KSK** | 1 (add+remove DS in sequence) | Larger (two KSKs present) | Simplest; default for most deployments |
| **Double-DS** | 2 (add new DS → swap KSK → remove old DS) | Normal | Two TTL waits; smaller RRset |
| **Double-RRset** | 2 | Larger | Fastest rollover window; most complex |

**`dnssec-policy` automates all of this**: BIND 9 manages all timing, key generation, and transitions. The only remaining manual task is uploading the DS record to the parent registrar (automatable with CDS/CDNSKEY + `parental-agents`).

**Manual Double-DS (for reference)**:
1. Generate new KSK
2. **Add** new DS record to parent zone (both old and new DS exist simultaneously)
3. Wait for old DS TTL to expire everywhere
4. Switch zone to use new KSK for signing
5. **Remove** old DS record from parent

Sources: [[how-does-dnssec-work-cloudflare]], [[dnssec-guide-bind9]]

---

## ZSK Rollover: Low-Cost by Design

ZSK rollover requires no parent coordination — the new ZSK public key is added to the zone's own DNSKEY RRset (signed by the KSK), and after the TTL expires the old ZSK is removed. Because the KSK (which is in the parent's DS record) doesn't change, the parent is uninvolved.

This is precisely why the two-key architecture exists. When using `dnssec-policy`, ZSK rollover is fully automatic.

---

## CSK: The Single-Key Variant

The BIND 9 `default` policy creates a **Combined Signing Key (CSK)** — a single key pair that serves both KSK and ZSK roles. This simplifies management at the cost of rollover flexibility: since the CSK hash is in the parent DS record, every rollover requires parent coordination (no low-cost ZSK-only rollover is possible).

Use CSK for: low-traffic zones, internal zones, zones where simplicity > rollover agility.
Use KSK/ZSK split for: high-traffic public zones, zones needing frequent ZSK rotation.

Source: [[dnssec-guide-bind9]]

---

## KSK as "Parent", ZSK as "Child" — Within One Name Server

An elegant way to understand the KSK/ZSK relationship comes from CS161's framing:

> "You can think of the KSK as the 'parent' and the ZSK as the 'child,' both within one name server."

Just as a **parent zone endorses a child zone** by signing the child's KSK hash (the DS record), a **KSK endorses its own ZSK** by signing the DNSKEY RRset that contains the ZSK. The same trust delegation mechanism operates at two scales simultaneously:

| Scale | Who endorses | What is endorsed | How |
|-------|-------------|-----------------|-----|
| **Inter-zone** | Parent zone | Child zone's KSK | DS record (hash of child KSK), signed by parent ZSK |
| **Intra-zone** | KSK | ZSK | RRSIG over DNSKEY RRset (containing ZSK public key), signed by KSK |

This framing explains why the **DS record** is the concrete cross-zone link: it is the parent zone's "endorsement receipt" for the child zone's KSK — the same trust handoff, one level up.

Source: [[cs161-dnssec]]

---

## Insecure Delegation

If a child zone has DNSSEC records but no DS record in the parent:
- The child zone is **insecure** — resolvers treat it as unsigned
- Resolvers will not attempt validation — they accept responses without checking signatures
- This is not an error; it's the intended behavior for zones that haven't been enrolled in DNSSEC

If a child zone has a DS record in the parent but its own DNSKEY records don't match (e.g. KSK was rolled without updating the DS):
- This is a **broken chain** / **security lameness** — resolver returns SERVFAIL

---

## What Breaks the Chain

| Break Point | Resolver Behavior | Root Cause |
|-------------|------------------|-----------|
| No DS in parent | Insecure delegation — no SERVFAIL, no validation | Zone signed locally, DS never submitted |
| DS/KSK hash mismatch (security lameness) | SERVFAIL | KSK rolled without updating parent DS |
| RRSIG expired | SERVFAIL | Zone not re-signed on schedule |
| DNSKEY record missing | SERVFAIL | Zone file misconfiguration |
| Clock skew outside RRSIG validity window | SERVFAIL | NTP not running; server clock drifted |
| Algorithm not supported by resolver | SERVFAIL | Old resolver, new signing algorithm |
| NSEC/NSEC3 gap | SERVFAIL on NXDOMAIN queries | Zone signing incomplete |

---

## Emergency KSK Rollover

When a KSK is compromised, the zone administrator faces a dilemma: act fast (break the chain temporarily) or act carefully (maintain service continuity while rolling). Source: [[rfc-6781]].

**Option A — Keep chain of trust intact** (slower, but no SERVFAIL):
1. Introduce new KSK alongside compromised KSK
2. Lower DNSKEY TTL to drain caches faster
3. Sign DNSKEY RRset with short validity period (expires ~when new DS appears in parent)
4. Upload new DS to parent
5. Wait for new DS to propagate + old DS TTL to expire
6. Remove compromised KSK
- **Risk**: Attacker can continue using compromised key during the transition

**Option B — Break chain of trust / zone goes Bogus** (faster, but disrupts service):
- Replace KSK immediately → zone appears Bogus to validators until new DS in parent

**Option C — Zone goes Insecure** (most disruptive, but cleanest escape):
- Remove DS from parent entirely → zone goes insecure → re-key cleanly → re-submit DS
- Used when the losing DNS operator won't cooperate in a migration

---

## Stand-By Keys (Pre-Positioned Emergency Response)

Pre-publishing keys so that emergency rollover is fast. Source: [[rfc-6781]].

- **Stand-by ZSK**: Generate offline, publish its DNSKEY in the zone (not yet signing). Activation = add signing with it.
- **Stand-by KSK**: Generate offline, get its DS pre-published in parent zone. Activation = start signing with it.

When disaster strikes (key inaccessible or compromised): activate stand-by keys → sign zone → remove old keys after TTL expiry. Dramatically reduces recovery time because the new DS is already in the parent.

---

## Mentioned In

- [[how-does-dnssec-work-cloudflare]] — DS record mechanics, KSK rollover procedure, Root Signing Ceremony
- [[kaminsky-exploit]] — first-principles motivation: why authentication was needed
- [[dnssec-guide-bind9]] — 12-step validation walkthrough, trust anchor config (`dnssec-validation auto`), three rollover methods, CSK variant, caching as the root of rollover complexity
- [[rfc-6781]] — formal rollover procedures (Pre-Publish ZSK, Double-KSK), emergency rollover options (keep chain / go Bogus / go Insecure), stand-by keys
- [[cs161-dnssec]] — KSK/ZSK pedagogical framing ("KSK as parent, ZSK as child within one name server"), inter-zone vs intra-zone endorsement comparison
