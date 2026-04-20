---
title: "DNSSEC in BIND"
type: concept
tags: ["dnssec", "bind", "named", "configuration", "zone-signing", "kasp", "edns0", "key-rollover", "inline-signing", "csk", "key-states", "parental-agents"]
created: 2026-04-13
updated: 2026-04-14
source_count: 3
confidence: high
status: active
related: ["[[dnssec-internals]]", "[[chain-of-trust]]", "[[dnssec-debugging]]", "[[Internet-Systems-Consortium]]", "[[DNSViz]]"]
open_questions:
  - "Is the DNSSEC Multi-Signer Model (RFC 8901) relevant to the target environment?"
  - "What is the exact behavior of `parental-agents` when the parent's DS state cannot be resolved automatically?"
contradictions: []
---

# DNSSEC in BIND

> **Domain**: BIND 9 / DNS Operations
> **Applies to**: BIND 9.18 LTS (current recommended; features require 9.16+)

---

## Definition

BIND 9's DNSSEC implementation covers two roles: **authoritative server** (signing zones) and **recursive resolver** (validating responses). Both are controlled via `named.conf`. The modern approach (`dnssec-policy`) automates virtually all ongoing key and signature management.

---

## Operational Prerequisites (Before You Sign Anything)

From [[dnssec-and-bind9-isc]] and [[dnssec-guide-bind9]]:

1. **EDNS0 must not be dropped** — DNSSEC requires EDNS0 for larger message sizes and the `DO` bit. Network devices (firewalls, load balancers) that silently strip EDNS0 packets will silently break DNSSEC. **Enterprise firewalls are a common culprit.**
2. **NTP / clock sync** — RRSIG records carry expiry timestamps. A significantly skewed clock makes valid signatures appear expired → SERVFAIL. **NTP sync is a hard operational dependency.**
3. **TCP port 53 must be open** — DNSSEC responses regularly exceed the 512-byte UDP limit. Firewalls blocking TCP/53 cause intermittent failures.
4. **All secondaries must support DNSSEC** — Secondaries that don't support DNSSEC will serve unsigned data.
5. **Parent zone must support DNSSEC** — Contact your registrar/parent zone operator before starting. The DS record upload is the final step to complete the chain.

---

## Enabling Validation (Resolver Side)

```nginx
# named.conf — basic validation
options {
    dnssec-validation auto;   /* default in BIND 9.18+ — usually no change needed */
};
```

| Setting | Behavior |
|---------|----------|
| `auto` | **Default in 9.18+.** BIND ships with root trust anchor, auto-updates via RFC 5011 (`managed-keys`). No manual key config needed. |
| `yes` | Validation enabled, but requires manually configured `trust-anchors`. |
| `no` | Validation disabled. **Never in production.** |

**How to verify validation is working:**
```bash
# The 'ad' flag in the response means validation passed
dig @resolver ftp.isc.org A +dnssec
# Look for: flags: qr rd ra ad

# To confirm SERVFAIL is a DNSSEC failure (not a network problem):
dig @resolver www.dnssec-failed.org A        # → SERVFAIL (expected with validation on)
dig @resolver www.dnssec-failed.org A +cd    # → NOERROR (bypasses validation)
# If the second succeeds and first fails → DNSSEC validation is working correctly
```

Source: [[dnssec-guide-bind9]]

---

## Signing a Zone: Modern Approach (`dnssec-policy`)

The entire signing setup is **3 config lines + `rndc reconfig`**:

```nginx
zone "example.com" in {
    type primary;
    file "example.com.zone";
    dnssec-policy default;    /* ← this line signs the zone */
};
```

`dnssec-policy` can be placed at `zone`, `view`, or `options` scope. BIND then:
1. Generates DNSSEC keys
2. Publishes DNSKEY records in the zone
3. Signs all RRsets (generating RRSIG records)
4. Generates NSEC/NSEC3 records
5. Re-signs automatically before signatures expire
6. Rolls keys according to the policy schedule

**The original zone file (`example.com.zone`) is never modified.** BIND generates a separate `example.com.zone.signed` (raw binary format) and `example.com.zone.signed.jnl` journal file.

> **Diagnostic tip**: If `*.signed.jnl` is missing, the zone was NOT signed — even if BIND logs look clean and the `*.signed` file exists.

Verify signing is complete:
```bash
rndc signing -list example.com
# Shows: "Done signing with key 10376/ECDSAP256SHA256"
```

Source: [[dnssec-guide-bind9]]

---

## The `default` Policy: What It Actually Does

The built-in `default` policy creates a **CSK (Combined Signing Key)** — a single key that acts as both KSK and ZSK. From [[dnssec-kasp-policy]], the exact parameters:

| Parameter | `default` value |
|-----------|----------------|
| Algorithm | ECDSAP256SHA256 (`ecdsa256`) |
| Key type | CSK (not separate KSK/ZSK) |
| Key lifetime | `unlimited` (no automatic rollover) |
| Signature validity | 14 days (`P14D`) |
| Signature refresh | 5 days (`P5D`) |
| DNSKEY TTL | 1 hour (`PT1H`) |
| Max zone TTL | 1 day (`P1D`) |
| Parent DS TTL | 1 day (`P1D`) |
| Publish/retire safety | 1 hour (`PT1H`) |
| Purge keys after delete | 90 days (`P90D`) |
| Denial of existence | NSEC (not NSEC3) |

The CSK approach simplifies configuration but means every key change requires parent DS interaction (since the CSK acts as KSK). For zones where ZSK-only automation is desired, use a custom policy with explicit `ksk` and `zsk` definitions.

---

## Inline Signing: Version-Sensitive Configuration

`dnssec-policy` requires either dynamic updates or inline-signing to be active. This is version-gated:

| BIND version | Inline signing behavior |
|---|---|
| 9.16 / 9.18 | Must explicitly add `inline-signing yes;` to zone config (or globally) |
| **9.20+** | **`inline-signing` defaults to `yes`** — `dnssec-policy` alone is sufficient |

**9.20 upgrade pitfall**: If a zone was signed with dynamic updates (no inline-signing) on 9.18, upgrading to 9.20 without adding `inline-signing no;` causes BIND to sign the already-signed zone file again, producing corrupted double-signed `.signed` files. Three mitigations:
1. Revert to unsigned first: use `dnssec-policy "insecure";`, then upgrade
2. Add `inline-signing no;` explicitly per zone before upgrading
3. Add `inline-signing no;` inside the `dnssec-policy` block

Source: [[dnssec-kasp-policy]]

---

## Custom Policy Example

```nginx
dnssec-policy standard {
    dnskey-ttl 600;
    keys {
        ksk lifetime 365d algorithm ecdsap256sha256;
        zsk lifetime 60d  algorithm ecdsap256sha256;
    };
    max-zone-ttl 600;
    parent-ds-ttl 600;
    parent-propagation-delay 2h;
    publish-safety 7d;
    retire-safety 7d;
    signatures-refresh 5d;
    signatures-validity 15d;
    zone-propagation-delay 2h;
};

zone "example.net" in {
    dnssec-policy standard;
    ...
};
```

Key parameters explained:
- `publish-safety` / `retire-safety` — Extra buffer time before a key starts/stops being used, in case rollover doesn't go to plan
- `parent-propagation-delay` — How long it takes for DS changes at the registrar to propagate
- `signatures-validity` — How long RRSIG records are valid; default is 14 days

---

## Uploading to the Parent Zone (Completing the Chain)

Until the parent zone publishes your DS record, your zone appears unsigned to the world:

```bash
# Generate DS record from your KSK
dnssec-dsfromkey Kexample.com.+013+10376.key
# → example.com. IN DS 10376 13 2 B92E22CAE0...

# Verify parent has published it
dig example.com. DS
# Look for DS record matching your key tag
```

If the parent supports CDS/CDNSKEY (RFC 7344/8078), BIND can signal the parent automatically without manual DS upload.

---

## Key Rollover: Why It's Complicated

The root cause of rollover complexity: **recursive resolvers cache DNSKEY records at their TTL**. If you swap keys instantly:
- Resolvers with the old key cached try to verify signatures made with the new key
- Verification fails → SERVFAIL

All rollover methods are designed to let caches expire organically before cutting over.

### ZSK Rollover Methods

| Method | How | Zone size impact |
|--------|-----|-----------------|
| Pre-Publication | Publish new ZSK → wait 1 TTL → start signing with new ZSK → wait 1 TTL → remove old ZSK | Minimal size increase |
| Double-Signature | Publish new ZSK + sign with both → wait 1 TTL → remove old ZSK + old RRSIGs | Zone doubles during rollover |

### KSK Rollover Methods (require parent DS change)

| Method | Steps | Parent interactions |
|--------|-------|---------------------|
| **Double-KSK** | Add new KSK → sign DNSKEY with both → update parent DS → wait TTL → remove old KSK | 1 |
| **Double-DS** | Add new DS to parent → wait TTL → swap KSK at child → wait TTL → remove old DS | 2 |
| **Double-RRset** | Add new KSK + new DS simultaneously → wait TTL → remove old KSK + old DS | 2 (fastest end-to-end) |

**With `dnssec-policy`, BIND automates all of this.** The only remaining manual task for ZSK rollovers is nothing — they are fully automatic. For KSK/CSK rollovers, the operator must signal DS changes to the parent (see below).

Source: [[dnssec-guide-bind9]]

---

## KSK/CSK Rollover: The Two Required `rndc` Signals

KSK and CSK rollovers **cannot be fully automated without `parental-agents`**. BIND pauses at two points waiting for the operator to confirm parent zone DS changes:

```bash
# 1. After you verify the NEW DS is published in the parent:
rndc dnssec -checkds -key <new-key-id> published example.com

# 2. After you verify the OLD DS has been REMOVED from the parent:
rndc dnssec -checkds -key <old-key-id> withdrawn example.com
```

**How to verify DS state in parent before signaling:**
```bash
dig +dnssec +norecurse +multi @<parent-ns-ip> example.com DS
# The key ID in the RRSIG shows which key is signing the DS
# When the new key ID appears → safe to signal `published`
```

**`parental-agents` — automated DS signaling:**
```nginx
zone "example.com" {
    dnssec-policy default;
    parental-agents { 192.0.2.1; };   /* trusted parent NS IP */
};
```
With `parental-agents`, BIND queries the parent zone directly, detects DS changes, and automatically signals `published`/`withdrawn` without operator intervention.

**Without the `rndc -checkds` signals** (and without `parental-agents`), BIND holds the rollover in `DSState: rumoured` indefinitely. The rollover appears stuck.

Source: [[dnssec-kasp-policy]]

---

## Manual Signing Tools (Legacy / Special Cases)

Use when `dnssec-policy` doesn't fit (e.g., offline KSK, HSM):

```bash
# Generate KSK and ZSK
dnssec-keygen -a ECDSAP256SHA256 example.com          # → ZSK (flags: 256)
dnssec-keygen -a ECDSAP256SHA256 -f KSK example.com   # → KSK (flags: 257)

# Set key lifecycle dates
dnssec-settime -P 20260101 -A 20260115 -I 20270115 -D 20270201 Kexample.com.+013+XXXXX.key

# Sign the zone
dnssec-signzone -t -N INCREMENT -o example.com \
    -f example.com.signed.db example.com.db \
    Kexample.com.+013+ZSK.key Kexample.com.+013+KSK.key

# Generate DS record for parent
dnssec-dsfromkey Kexample.com.+013+KSK.key
```

Key file naming: `K<zone>+<algorithm_number>+<key_tag>.key/.private/.state`
- Algorithm 13 = ECDSAP256SHA256
- Flags 256 in DNSKEY record = ZSK; 257 = KSK or CSK

---

## `dnssec-policy` vs `auto-dnssec`: The Generation Shift

| Feature | `auto-dnssec` (legacy) | `dnssec-policy` / KASP (current) |
|---------|------------------------|----------------------------------|
| Introduced | BIND 9.7 | BIND 9.16 |
| Key generation | Manual | Automated |
| Key rollover | Semi-automated | Fully automated |
| Status | **Removed in BIND 9.19.16+** | Current recommendation |

> ⚠️ `auto-dnssec` is not merely deprecated — it has been **removed** from BIND 9.19.16 and later (including BIND 9.20+). Any zone still using `auto-dnssec` must migrate before upgrading to 9.20.

[[Internet-Systems-Consortium]] explicitly flags all pre-KASP documentation as outdated. Any documentation using `auto-dnssec maintain` or `dnssec-keygen` + `dnssec-signzone` as the primary workflow is describing the legacy approach.

Source: [[dnssec-kasp-policy]]

---

## BIND DNSSEC Tool Reference

| Tool | Purpose |
|------|---------|
| `dnssec-keygen` | Generate key pairs (manual mode) |
| `dnssec-settime` | Set/modify key lifecycle timestamps |
| `dnssec-signzone` | Manually sign a zone file |
| `dnssec-dsfromkey` | Generate DS record from DNSKEY file |
| `dnssec-verify` | Offline verify a signed zone file |
| `named-compilezone` | Convert raw binary zone to human-readable text |
| `rndc signing -list` | Check signing progress for a zone |
| `rndc reconfig` | Apply config changes (trigger zone signing) |
| `rndc nta` | Add/list Negative Trust Anchors |
| `rndc dnssec -checkds` | Manually signal DS published/withdrawn |
| `nsec3hash` | Compute NSEC3 hash for a name (debugging) |

---

## Migration from `auto-dnssec` to `dnssec-policy`

Prerequisites before migrating an existing zone:
1. **Private key format v1.3** — check `.private` file for `Private-key-format: v1.3`. If older, run `dnssec-settime -f <keyfile>` to upgrade.
2. **No rollover in progress** — check with `dig DNSKEY <zone>`. Two keys of the same algorithm/flag means a rollover is underway. Wait for it to complete first.
3. **Match existing keys** — create a custom `dnssec-policy` that matches your current algorithm and KSK/ZSK structure. Switching to `default` on a zone with a KSK/ZSK split immediately triggers an algorithm rollover.

**Example: migrating a RSASHA256 KSK/ZSK zone:**
```nginx
dnssec-policy "myway" {
    keys {
        ksk lifetime unlimited algorithm rsasha256 2048;
        zsk lifetime P60D algorithm rsasha256 1024;
    };
};
zone "example.com" { dnssec-policy myway; };
```

**After migration**, if `DSState` is stuck in `rumoured`:
```bash
rndc dnssec -checkds published example.com
# Tells BIND the DS is already in the parent — unsticks the rollover
```

Source: [[dnssec-kasp-policy]]

---

## SOA Expiry and Signature Validity Alignment

From [[rfc-6781]]:

**Problem**: Secondary servers zone-expire based on the SOA expire timer, independently of RRSIG expiry. If secondaries cannot reach the primary and the RRSIG expires before the SOA expire timer, secondaries serve Bogus (signature-expired) data.

**Guidance**: Set `SOA expire ≈ 1/3 to 1/4 of the signature validity period`.

With BIND `dnssec-policy default` (14-day signature validity):
- Recommended SOA expire: 3–5 days
- This ensures secondaries fail to resolve (zone expires) before signatures expire
- If you see SERVFAIL only from secondary servers → check SOA expire vs RRSIG expiry

---

## DNS Operator Change Procedure

When migrating DNSSEC from one DNS provider to another (from [[rfc-6781]]):

**Cooperating operators** (gaining operator gets full cooperation):
1. Losing operator pre-publishes gaining operator's ZSK + KSK in zone
2. Wait for propagation
3. Re-delegation: update NS + DS records in parent to point to gaining operator
4. After old records expire: remove losing operator's keys

**Non-cooperating operators** (losing operator won't help):
- The only viable path is to make the zone **go Insecure** during the transition:
  1. Registry removes DS → zone goes Insecure
  2. Re-delegate to new operator
  3. New operator re-signs zone
  4. Submit new DS to parent → zone becomes secure again
- Brief Insecure period is unavoidable

---

## Mentioned In

- [[dnssec-and-bind9-isc]] — operational prerequisites, KASP confirmation
- [[dnssec-guide-bind9]] — signing, validation, rollover, troubleshooting
- [[dnssec-kasp-policy]] — `default` policy parameters, inline-signing version gates, KSK rollover signals, `parental-agents`, migration procedure, `auto-dnssec` removal
- [[rfc-6781]] — ZSK/KSK rollover procedure details (Pre-Publish, Double-KSK), SOA expiry guidance, DNS operator change procedure (cooperating vs non-cooperating)
