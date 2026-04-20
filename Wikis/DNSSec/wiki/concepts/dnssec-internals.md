---
title: "DNSSEC Internals"
type: concept
tags: ["dnssec", "cryptography", "dns", "record-types", "rrset", "ksk", "zsk", "csk", "edns0", "key-lifecycle", "nsec3", "kasp-states", "signature-timing", "dnskey-flags", "opt-record", "offline-signing"]
created: 2026-04-13
updated: 2026-04-14
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

Every DNSSEC zone has at minimum two key pairs; a simplified single-key variant also exists:

| Key type | Signs | Rotation frequency | Rollover cost | When to use |
|----------|-------|--------------------|---------------|-------------|
| **ZSK** (Zone Signing Key) | All RRsets in the zone | Frequently (months) | Low — no parent coordination | High-traffic zones needing frequent rotation |
| **KSK** (Key Signing Key) | Only the DNSKEY RRset | Infrequently (yearly+) | **High** — requires parent DS update | Standard two-key deployments |
| **CSK** (Combined Signing Key) | Everything (serves both roles) | Infrequently | **High** — always requires parent DS update | Simple/internal zones; BIND `dnssec-policy default` |

**Why two keys?**
The KSK's hash is the DS record stored in the **parent zone**. Updating the DS record requires coordinating with the parent (registrar, TLD operator) — a slow, error-prone process. By separating the KSK from the ZSK:
- ZSK can be rotated frequently for security without touching the parent
- KSK remains stable, minimizing expensive parent coordination
- A compromised ZSK can be replaced immediately without a parent DS change

The KSK signs the DNSKEY RRset (which contains both the KSK public key and the ZSK public key). This creates an internal chain: KSK → endorses → ZSK → signs → all other records.

**BIND `dnssec-policy default`** uses a **CSK** with algorithm ECDSAP256SHA256. There is no separate ZSK — the single key handles everything. This is simpler but every key rollover requires parent DS interaction.

Sources: [[how-does-dnssec-work-cloudflare]], [[dnssec-guide-bind9]]

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
- [[rfc-6781]] — formal signature timing terminology (validity period, Re-Sign Period, Refresh Period, key effectivity period); NSEC3 2012 guidance (superseded by RFC 9276)
- [[cs161-dnssec]] — DNSKEY flag values (256=ZSK, 257=KSK, SEP bit); OPT pseudo-record / ADDITIONAL count +1; offline signing DoS motivation
