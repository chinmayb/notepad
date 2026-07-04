---
title: "DNSSEC Internals"
type: concept
tags: ["dnssec", "cryptography", "dns", "record-types", "rrset", "ksk", "zsk", "csk", "edns0", "key-lifecycle", "nsec3", "kasp-states", "signature-timing", "dnskey-flags", "opt-record", "offline-signing"]
created: 2026-04-13
updated: 2026-05-03
source_count: 6
confidence: high
status: active
related: ["[[chain-of-trust]]", "[[dnssec-vs-dns-resolution]]", "[[dnssec-in-bind]]", "[[dnssec-debugging]]"]
open_questions:
  - "What is the practical difference between RSA/SHA-256 and ECDSA P-256 for zone signing today?"
contradictions:
  - "RFC 6781 (2012) §5.3.2 recommended 100 NSEC3 iterations. RFC 9276 (2022) supersedes this: use 0 iterations. Modern GPUs make higher iterations ineffective while adding server CPU cost. Current guidance follows RFC 9276."
---

# DNSSEC Internals

> **Domain**: DNS Security / Cryptography
> **Also known as**: DNS Security Extensions

---

## Definition

DNSSEC is a suite of IETF extensions to DNS that adds **origin authentication** and **data integrity** to DNS responses using public-key cryptography. It does **not** add confidentiality — DNS responses remain visible to anyone on the wire.

Key contrast: **HTTPS encrypts traffic** (confidentiality). **DNSSEC signs responses** (authentication only). They solve different problems. A zone can be DNSSEC-signed and its contents still be readable by any observer.

---

## RRsets: The Atomic Unit of Signing

Before any signing can happen, all records of the **same type** at the **same label** are grouped into a **Resource Record Set (RRset)**.

Example: three `AAAA` records at `api.example.com` → one AAAA RRset. A single RRSIG covers the entire RRset, not individual records.

**Consequence**: a resolver that receives any record from an RRset must validate the whole RRset together. You cannot validate a single record in isolation.

Source: [[how-does-dnssec-work-cloudflare]]

---

## Key Record Types

| Record | Who Stores It | Purpose |
|--------|--------------|---------|
| `DNSKEY` | Zone's own nameserver | Public key (KSK or ZSK) for the zone |
| `RRSIG` | Zone's own nameserver | Cryptographic signature over one RRset |
| `DS` | **Parent** zone's nameserver | Hash of the child zone's KSK — the inter-zone trust link |
| `NSEC` | Zone's own nameserver | Authenticated proof that a name does NOT exist (sorted linked list) |
| `NSEC3` | Zone's own nameserver | Same as NSEC but uses hashed names to prevent zone walking |
| `CDNSKEY` / `CDS` | Child zone's nameserver | Child signals desired DS update to parent (automation, RFC 7344 / 8078) |

Source: [[how-does-dnssec-work-cloudflare]]

---

## Key Types: KSK, ZSK, and CSK

### The Validator Doesn't Care — The Operator Does

A foundational point that the spec makes explicit but is often missed: **the DNSSEC validation protocol does not distinguish between key types**. From [[rfc-6781]] §3.1:

> "The DNSSEC validation protocol does not distinguish between different types of DNSKEYs. The motivations to differentiate between keys are purely operational; validators will not make a distinction."

A resolver presented with the DNSKEY RRset sees a bag of public keys. It will try each one against an RRSIG until verification succeeds or fails. The KSK/ZSK distinction is a deployment convention layered *on top of* a uniform protocol — not a protocol feature. Once the DNSKEY RRset is signed by the KSK, **any key in that RRset can be used as a ZSK** as far as the protocol is concerned (RFC 6781 §3.1).

What separates them in practice is a single bit in the DNSKEY flags field — the **SEP (Secure Entry Point) bit** (see DNSKEY Flag Values section below) — and an operational policy about *which* key signs *what*.

### What Each Key Actually Signs

| Key type | Signs | Rotation frequency | Rollover cost | When to use |
|----------|-------|--------------------|---------------|-------------|
| **ZSK** (Zone Signing Key) | Every non-DNSKEY RRset in the zone (A, AAAA, MX, NSEC/NSEC3, etc.) | Frequently (months) | Low — no parent coordination | High-traffic zones needing frequent rotation |
| **KSK** (Key Signing Key) | **Only the apex DNSKEY RRset** | Infrequently (yearly+) | **High** — requires parent DS update | Standard two-key deployments |
| **CSK** (Combined Signing Key) | Everything (serves both roles) | Infrequently | **High** — always requires parent DS update | Simple/internal zones; BIND `dnssec-policy default` |

The DNSKEY RRset contains *both* the KSK public key and the ZSK public key. The KSK signs that RRset, vouching for the ZSK. The ZSK then signs everything else. This creates the internal chain:

```
parent DS → (validates) KSK → (signs DNSKEY RRset containing) ZSK → (signs) all other RRsets
```

Source: [[how-does-dnssec-work-cloudflare]], [[rfc-6781]]

### Why Two Keys? — The Operational Argument

The KSK's hash is published as the **DS record in the parent zone**. Any KSK change must be reflected in the parent's DS RRset. That is the expensive operation. The cost comes from three things:

1. **Cross-organizational coordination.** The parent zone is operated by a different entity (a TLD registry, a registrar, an enterprise IT group) and updates flow through their interface — often a web console, sometimes a ticketing system.
2. **TTL-bounded waiting.** The parent must publish the new DS, then **wait for the original DS's TTL to expire** in caches before removing it. Per [[how-does-dnssec-work-cloudflare]]: "First, the parent needs to add the new DS record, then they need to wait until the TTL for the original DS record to expire before removing it."
3. **High blast radius on error.** A bungled DS update breaks the chain of trust → the entire zone goes Bogus → resolvers SERVFAIL.

Splitting the keys lets you decouple two different change cadences:

- **ZSK rotates frequently, locally, with zero parent involvement.** A compromised ZSK can be replaced *immediately* — the operator generates a new ZSK, signs it into the DNSKEY RRset with the existing KSK, and re-signs the zone. No external coordination.
- **KSK remains stable, infrequently rotated.** This minimizes how often you have to dance with the parent.

Per [[rfc-6781]] §3.1: *"If there has been an event that increases the risk that a ZSK is compromised, it can be simply replaced with a ZSK rollover. The new RRset is then re-signed with the KSK."*

### Why Two Keys? — The Storage / Risk-Tier Argument

The split also lets the two keys live in different security tiers. From [[rfc-6781]] §3.1:

> "A KSK can be stored off-line or with more limitations on access control than ZSKs, which need to be readily available for operational purposes such as the addition or deletion of zone data. A KSK stored on a smartcard that is kept in a safe, combined with a ZSK stored on a file-system accessible by operators for daily routine use, may provide better protection against key compromise without losing much operational flexibility."

In other words:
- **ZSK** — must be online to sign zone updates as records change (dynamic updates, daily edits). Lives on the signing host's filesystem or HSM. Higher exposure surface.
- **KSK** — only needs to come out of the safe when the *DNSKEY RRset itself* changes (i.e., during a ZSK rollover or KSK rollover). Can live in cold storage / offline HSM / smartcard.

This is the **Offline KSK** pattern (see Offline Signing section). It is not free: it adds an out-of-band ceremony every time a ZSK is rolled. But it removes the KSK's private material from any internet-connected machine.

Per the BIND guide ([[dnssec-guide-bind9]]): *"For operational reasons, it is possible to keep the KSK offline. Doing so minimizes the risk of the key being compromised through theft or loss."*

### Why NOT Two Keys? — When the Single-Type Scheme Wins

[[rfc-6781]] is candid that the KSK/ZSK split is not always justified. The cryptanalysis argument *alone* does not motivate the split:

> "Suppose one differentiates between a KSK and a ZSK, whereby the KSK effectivity period is X times the ZSK effectivity period. Then, in order for the resistance to cryptanalysis to be the same for the KSK and the ZSK, the KSK needs to be X times stronger than the ZSK. […] When translated to asymmetric keys, the size difference is still too insignificant to warrant a key-split; it only marginally affects the packet size and signing speed." (RFC 6781 §3.1)

The split's value is **operational**, not cryptographic. RFC 6781 explicitly names cases where the split is weakly justified:
- Keys are stored on an HSM (storage exposure is already low)
- The KSK is known not to be used as a trust anchor anywhere
- The cost of operational complexity outweighs the flexibility benefit

In those cases the **Single-Type Signing Scheme** (RFC 6781 terminology) — which BIND calls a **CSK** — is reasonable. One key signs everything; every rollover involves the parent. Simpler operationally, fewer moving parts, but every rotation is "expensive."

**BIND `dnssec-policy default`** uses a CSK with ECDSAP256SHA256 and *never rotates it*. The conventional KSK/ZSK split must be enabled by writing a custom policy ([[dnssec-guide-bind9]] shows the typical "KSK yearly, ZSK every two months" pattern).

### Detection: Which Key Is Which?

Two practical recipes for identifying KSK vs. ZSK in a DNSKEY RRset:

1. **Flag value** — KSK is `257`, ZSK is `256`. The difference is the SEP bit (bit 15). See the DNSKEY Flag Values section below.
2. **Parity test** ([[rfc-6781]] §3.1): *"If the flag field is an odd number, it is a KSK; otherwise, it is a ZSK."* (Because the SEP bit is bit 0 in the flags field's low byte after byte-swap, setting it makes the integer odd.)

Sources: [[how-does-dnssec-work-cloudflare]], [[dnssec-guide-bind9]], [[rfc-6781]], [[cs161-dnssec]]

---

## The Intra-Zone Validation Sequence

When a DNSSEC-aware resolver receives a response:

```
Step 1: Receive target RRset + its RRSIG (signed by private ZSK)
Step 2: Fetch DNSKEY RRset → get public ZSK + public KSK + RRSIG over DNSKEY (signed by private KSK)
Step 3: Use public KSK → verify RRSIG over DNSKEY RRset → confirms ZSK is authentic
Step 4: Use public ZSK → verify RRSIG over target RRset → confirms answer is authentic
Step 5: Trust the KSK itself via the DS record in the parent zone (see [[chain-of-trust]])
```

The DNSKEY RRset (step 2) can be cached — resolvers don't re-fetch it on every query.

---

## EDNS0, the DO Bit, and the OPT Pseudo-Record

- DNSSEC records are only returned if the querier sets the **DO (DNSSEC OK)** bit in the EDNS0 OPT record
- Plain DNS resolvers without EDNS0 support won't receive DNSSEC records — backwards compatible by design
- `dig +dnssec` sets the DO bit; plain `dig` does not
- `delv` always sets the DO bit and performs full chain validation
- **Operational risk**: Network devices (firewalls, load balancers) that silently drop or strip EDNS0 packets will break DNSSEC — even if the zone is correctly signed. EDNS0 filtering is a common root cause of DNSSEC failures in enterprise environments.

### The OPT Pseudo-Record Explains the "+1 ADDITIONAL" Anomaly

The EDNS0 OPT record is encoded as an **additional record** for backwards compatibility. This means `dig` output always shows `ADDITIONAL` count as **one higher than the number of actual additional records**.

```
;; QUESTION SECTION:
;eecs.berkeley.edu.   IN   A

;; ANSWER SECTION:
eecs.berkeley.edu.   300   IN   A   23.185.0.1

;; ADDITIONAL SECTION:
;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags: do; udp: 1232
```

The `;; OPT PSEUDOSECTION:` is that extra count. `flags: do` confirms the DO bit is set; `udp: 1232` is the EDNS0 buffer size. Source: [[cs161-dnssec]].

Sources: [[how-does-dnssec-work-cloudflare]], [[dnssec-and-bind9-isc]], [[cs161-dnssec]]

---

## DNSKEY Flag Values

When examining `dig DNSKEY` output, the flags integer distinguishes key roles:

| Flag value | Key type | Meaning |
|-----------|---------|---------|
| **256** | ZSK | Zone Signing Key — signs all RRsets in the zone |
| **257** | KSK | Key Signing Key — signs only the DNSKEY RRset; SEP bit set |

The difference is the **SEP (Secure Entry Point) bit** — bit 15 of the DNSKEY flags field. KSKs have it set: `256 + 1 = 257`. This bit signals that this key's hash should be published as a DS record in the parent zone.

Example `dig DNSKEY` output:
```
; <<>> DiG 9.18 <<>> DNSKEY berkeley.edu.
;; ANSWER SECTION:
berkeley.edu.  3600  IN  DNSKEY  256 3 13 <base64-key>   ; ZSK
berkeley.edu.  3600  IN  DNSKEY  257 3 13 <base64-key>   ; KSK
berkeley.edu.  3600  IN  RRSIG   DNSKEY 13 2 3600 ...
```

The `3` is the protocol (always 3 = DNSSEC). The `13` is the algorithm (ECDSAP256SHA256). Source: [[cs161-dnssec]].

---

## Offline Signing

- Zone records are signed **ahead of time** and stored with the zone data
- **Why not online signing?** If signatures were computed on-demand per query, an attacker could flood the server with NXDOMAIN queries for nonexistent names, exhausting CPU with signing operations → DoS. Offline signing eliminates this attack surface entirely. Source: [[cs161-dnssec]].
- Consequence: signatures have an expiry (`RRSIG` Expiration field) — zones must be re-signed before expiry
- Typical re-signing interval: daily to weekly for production zones

---

## Key Lifecycle Metadata

BIND tracks each key through a lifecycle with explicit timestamps. These are set by `dnssec-policy` automatically or by `dnssec-settime` manually. Source: [[dnssec-guide-bind9]].

| Timestamp | Meaning |
|-----------|---------|
| **Created** | Key pair generated on disk |
| **Publish** | Public key added to DNSKEY RRset; resolvers start caching it |
| **Activate** | Key begins signing records |
| **Inactive** | Key stops signing; already-signed records with this key are still valid until RRSIG expiry |
| **Delete** | Public key removed from DNSKEY RRset entirely |

**Why the gap between Inactive and Delete?** After a ZSK stops signing (Inactive), resolvers may still have old RRSIGs in cache that were signed with that key. The public key must remain in the DNSKEY RRset so those cached signatures can still be validated. Only after the RRSIG TTL expires can the public key be safely removed.

The `dnssec-settime` utility reads and writes these timestamps on `.key` and `.private` files. `rndc signing -list <zone>` shows current key states for a loaded zone.

---

## NSEC vs NSEC3 Decision Guidance

Both provide authenticated non-existence. The choice is a security vs. privacy trade-off. Source: [[dnssec-guide-bind9]].

| Feature | NSEC | NSEC3 |
|---------|------|-------|
| Zone walking | Trivial (sorted linked list of plaintext names) | Hard (hashed names; GPU-crackable) |
| Query performance | Faster | Slower (hash computation per lookup) |
| Implementation complexity | Simple | More complex (hash + iterations + salt config) |
| BIND config | `nsec3param` absent (default) | Add `nsec3param 1 0 0 -;` to zone |

**BIND guidance on NSEC3 parameters**:
- `iterations`: use **0** (zero). Higher values were historically recommended to slow brute-force, but modern GPUs make them irrelevant for defense while adding server-side CPU cost. RFC 9276 formalizes this.
- `salt`: use **`-`** (empty / zero-length). Random salts provided minimal security benefit and complicated key rollovers.
- `opt-out flag`: use **0** (disabled) unless the zone has very many unsigned delegations.

> Recommended minimal NSEC3: `nsec3param 1 0 0 -;`

---

## NSEC and NSEC3: Authenticated Non-Existence

**The problem**: DNS returns an empty answer for non-existent names. There's nothing to sign. An attacker could inject fake NXDOMAIN responses.

**NSEC solution**: The zone pre-computes a sorted linked list of all its record names. For a query for a non-existent name, the server returns the NSEC record for the adjacent names ("nothing exists between `blog.example.com` and `www.example.com`"). This NSEC record is signed, so the non-existence is authenticated.

**NSEC vulnerability — zone walking**: By querying for names alphabetically, an attacker can enumerate every record name in the zone by following the NSEC chain.

**NSEC3 mitigation**: Same mechanism, but uses hashed names instead of plaintext names. Makes zone walking computationally harder (attacker must crack hashes), but not impossible with a GPU. Use zero iterations and empty salt per RFC 9276 guidance above.

Source: [[how-does-dnssec-work-cloudflare]]

---

## KASP Key States (BIND Runtime)

BIND's KASP tracks four **runtime propagation states** for each key (distinct from the key lifecycle timestamps like Publish/Activate/Delete). These track whether each record type (DNSKEY, DS, RRSIG) has propagated to resolvers. Stored in `.state` files alongside each key. Source: [[dnssec-kasp-policy]].

| State | Meaning |
|-------|---------|
| `hidden` | Not yet published; no resolver knows about this record |
| `rumoured` | Published, but propagation not complete; some resolvers may have it, others not |
| `omnipresent` | Fully propagated; all resolvers have it (either cached or will refresh from auth server) |
| `unretentive` | Being retired; transitioning out of active use |

These states are tracked **per record type** in the `.state` file:
- `DNSKEYState` — is the DNSKEY record propagated to resolvers?
- `DSState` — is the DS record propagated to parent zone resolvers?
- `ZRRSIGState` — are zone-signing RRSIGs propagated?
- `KRRSIGState` — are DNSKEY-signing (KSK) RRSIGs propagated?

**`GoalState`** is the target state BIND is working toward. During normal operation it is `omnipresent`; during key removal it becomes `hidden`.

**`DSState: rumoured` stuck** = DS is already published in the parent but BIND doesn't know. Fix: `rndc dnssec -checkds published <zone>`.

This is separate from the key lifecycle timestamps (Created/Publish/Activate/Inactive/Delete) which are scheduling milestones. KASP states are observed propagation facts.

---

## Signature Timing Terminology

Formal terms from RFC 6781 ([[rfc-6781]]). These map directly to BIND `dnssec-policy` parameters.

| Term | Definition |
|------|-----------|
| **Signature validity period** | Time between RRSIG `Signature Inception` and `Signature Expiration` timestamps |
| **Signature publication period** | Duration the RRSIG is actually published in the zone (may be shorter if refreshed early) |
| **Key effectivity period** | Duration from earliest RRSIG inception to latest RRSIG expiration for all RRSIGs made with a key. Spans multiple validity periods. |
| **Re-Sign Period** | How often the signer visits the zone to refresh signatures approaching expiry |
| **Refresh Period** | Window before RRSIG expiration within which signatures must be refreshed. Must be > Re-Sign Period. |

**Key relationship**: `Re-Sign Period < Refresh Period < Signature Validity Period`

The gap between Refresh Period and Validity Period = operator's response window when signing fails. If the signer breaks and the Refresh Period expires, signatures start expiring → SERVFAIL.

**BIND `dnssec-policy` mappings:**
- `signatures-validity` = Signature validity period (default: 14 days)
- `signatures-refresh` = Refresh Period (default: 5 days)
- BIND re-signs automatically at the Refresh Period boundary

**SOA expiry guidance** (RFC 6781): Set the SOA expire field to approximately 1/3 to 1/4 of the signature validity period. This ensures secondary servers fail to resolve (zone expires) before signatures expire, preventing secondaries from serving signature-expired (Bogus) data.

> ⚠️ **Contradiction — NSEC3 iterations**: RFC 6781 §5.3.2 recommended up to 100 NSEC3 iterations as "reasonable." **RFC 9276 (2022) supersedes this**: use 0 iterations. Modern GPUs make any iteration count ineffective as a defense while adding server CPU cost per query. **Current wiki guidance follows RFC 9276** (0 iterations, empty salt).

---

## Mentioned In

- [[how-does-dnssec-work-cloudflare]] — primary source (RRsets, KSK/ZSK split, DS mechanics, NSEC/NSEC3)
- [[dnssec-and-bind9-isc]] — EDNS0 as hard prerequisite; DO bit behavior
- [[dnssec-guide-bind9]] — CSK type, key lifecycle metadata (Created/Publish/Activate/Inactive/Delete), NSEC3 parameter guidance (RFC 9276: 0 iterations, empty salt)
- [[dnssec-kasp-policy]] — KASP key states (hidden/rumoured/omnipresent/unretentive), `.state` file format, CSK rollover strategy
- [[rfc-6781]] — formal signature timing terminology (validity period, Re-Sign Period, Refresh Period, key effectivity period); NSEC3 2012 guidance (superseded by RFC 9276); §3.1 "validators don't distinguish key types" — KSK/ZSK split is purely operational; storage-tier and parent-coordination motivations for the split; cryptanalysis does NOT motivate the split; "odd flag = KSK" parity rule; Single-Type Signing Scheme = CSK
- [[cs161-dnssec]] — DNSKEY flag values (256=ZSK, 257=KSK, SEP bit); OPT pseudo-record / ADDITIONAL count +1; offline signing DoS motivation
