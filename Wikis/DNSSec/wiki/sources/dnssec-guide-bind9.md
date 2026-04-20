---
title: "DNSSEC Guide — BIND 9 (Official)"
type: source
tags: ["dnssec", "bind", "isc", "kasp", "dnssec-policy", "key-rollover", "validation", "troubleshooting", "signing"]
created: 2026-04-14
updated: 2026-04-14
source_count: 1
confidence: high
status: active
related: ["[[dnssec-in-bind]]", "[[chain-of-trust]]", "[[dnssec-internals]]", "[[dnssec-debugging]]", "[[Internet-Systems-Consortium]]"]
open_questions:
  - "How does BIND's CSK (Combined Signing Key) differ in operational practice from the KSK/ZSK split?"
  - "What is the exact behavior of `parental-agents` + `checkds` for automated DS upload?"
contradictions: []
---

# Source: DNSSEC Guide — BIND 9 (Official)

> **Raw file**: `raw/dnssec0guide/DNSSEC Guide — BIND 9 9.21.22-dev documentation.md`
> **Source URL**: https://bind9.readthedocs.io/en/latest/dnssec-guide.html
> **Author**: [[Internet-Systems-Consortium]] (originally Josh Kuo, DeepDive Networking)
> **Version**: BIND 9.18.0+ (some features require 9.16+)
> **Type**: Official ISC long-form guide — "why" and "how" of DNSSEC in BIND
> **Coverage axes**: BIND ops (Axis 5), Debugging (Axis 6), Internals (Axis 1 detail)

---

## Summary

The canonical ISC guide to DNSSEC in BIND 9. Covers the full lifecycle: enabling validation on a recursive resolver, signing an authoritative zone with `dnssec-policy`, understanding key types and rollover procedures, troubleshooting common failures, and advanced topics (NSEC/NSEC3 mechanics, offline KSK, DANE). The single richest practical source in the wiki.

---

## Key Claims

### Signing

1. **Signing a zone = 3 config lines** — Adding `dnssec-policy default;` to a zone clause and running `rndc reconfig` is all that is required. BIND generates keys, signs all records, creates NSEC/NSEC3 records, re-signs on expiry, and rolls ZSKs automatically.

2. **`dnssec-policy default` uses a CSK** — The built-in default policy creates a single Combined Signing Key (CSK) that serves both KSK and ZSK roles, using algorithm ECDSAP256SHA256. Unlike a manual KSK/ZSK split, the CSK is never automatically rolled (requires operator intervention for rollover).

3. **Inline signing keeps zone files clean** — The original `example.com.db` is untouched. BIND generates `example.com.db.signed` automatically. The absence of `*.signed.jnl` is a diagnostic indicator that the zone was NOT signed despite appearing to load successfully.

4. **Zone size multiplies after signing** — Signed zones are typically 2-3x larger due to RRSIG, NSEC/NSEC3, and DNSKEY records. Each record in the zone gets its own RRSIG.

### Validation

5. **`dnssec-validation auto` is the default in BIND 9.18+** — No config change needed. BIND ships with the root trust anchor and updates it automatically via RFC 5011 (`managed-keys`). Setting `dnssec-validation yes` requires manual trust anchor management.

6. **The `ad` flag is the validation signal** — `dig` output with `flags: qr rd ra ad` means DNSSEC validation succeeded. Absence of `ad` = not validated (either the zone isn't signed or validation is disabled). A SERVFAIL where `dig +cd` succeeds = DNSSEC validation failure.

7. **RRSIG validity is timestamp-gated** — Each RRSIG has `Signature Inception` and `Signature Expiration` timestamps. The resolver's clock must fall within this window. Clock skew outside the window = SERVFAIL even on correctly-signed zones.

### Key Rollover

8. **Rollover complexity comes from caching** — Recursive resolvers cache DNSKEY records at their TTL. Swapping keys instantly breaks validation for all resolvers with the old key cached. All rollover procedures wait for TTL expiry before transitioning.

9. **Three KSK rollover methods**:
   - **Double-KSK**: Add new KSK → update parent DS → remove old KSK (one parent interaction; bigger DNSKEY RRset during transition)
   - **Double-DS**: Add new DS to parent → swap KSK at child → remove old DS from parent (two parent interactions; smaller DNSKEY RRset)
   - **Double-RRset**: Fastest rollover; most complex; two parent interactions + bigger RRset
   
10. **`dnssec-policy` automates all of this** — If using `dnssec-policy`, the only manual task left is uploading DS records to the parent (and even that can be automated with CDS/CDNSKEY + `parental-agents`).

### Troubleshooting

11. **Security lameness is the #1 failure mode** — When the parent DS record doesn't match any DNSKEY in the child zone, the zone becomes entirely unresolvable. BIND log message: `"broken trust chain resolving..."`. Triage: `dig example.net DS` to get parent's view, then `dig +cd example.net DNSKEY` to get child's keys, compare key tags.

12. **BIND debug logging level 3 for DNSSEC traces** — Enable `category dnssec { severity debug 3; }` to see full validation chain in logs. Warning: this is CPU-intensive; not for production.

13. **`*.signed.jnl` file = zone is signed** — Missing `.signed.jnl` despite BIND appearing to process the zone normally is the diagnostic indicator that keys weren't loaded.

14. **Negative Trust Anchors (NTA)** — `rndc nta example.com` temporarily disables DNSSEC validation for a zone for 1 hour (default). A safety valve for when a third-party zone is misconfigured and causing SERVFAIL.

---

## Direct Quotes Worth Preserving

> "DNSSEC does not provide a secure tunnel; it does not encrypt or hide DNS data. It was designed with backwards compatibility in mind."

> "Once the zone is signed, BIND takes care of everything else. As for updating your zone file, you can continue to update it the same way as prior to signing your zone; the normal work flow of editing a zone file and using the rndc command to reload the zone still works as usual."

> "Do not be tempted to disable DNSSEC validation just because some names are failing to resolve. Remember, DNSSEC protects your DNS lookup from hacking."

> "If using `dnssec-policy`, you don't really need to concern yourself with the details of a key rollover: BIND 9 takes care of it all for you."

> "The signed DNS data is stored in `example.com.db.signed` and in the associated journal file [...] unless you see the `*.signed.jnl` file, your zone has not been signed."

---

## Config Recipes (Key Snippets)

### Signing a zone (minimal)
```nginx
zone "example.com" in {
    type primary;
    file "example.com.zone";
    dnssec-policy default;
};
```

### Custom policy
```nginx
dnssec-policy standard {
    dnskey-ttl 600;
    keys {
        ksk lifetime 365d algorithm ecdsap256sha256;
        zsk lifetime 60d algorithm ecdsap256sha256;
    };
    max-zone-ttl 600;
    parent-ds-ttl 600;
    parent-propagation-delay 2h;
    publish-safety 7d;
    retire-safety 7d;
    signatures-refresh 5d;
    signatures-validity 15d;
    signatures-validity-dnskey 15d;
    zone-propagation-delay 2h;
};
```

### Enabling validation (resolver)
```nginx
options {
    dnssec-validation auto;  /* default in BIND 9.18+ — usually no change needed */
};
```

### DNSSEC debug logging
```nginx
logging {
    channel dnssec_log {
        file "/var/log/dnssec.log";
        severity debug 3;
    };
    category dnssec { dnssec_log; };
};
```

---

## BIND DNSSEC Tool Reference

| Tool | Purpose |
|------|---------|
| `dnssec-keygen` | Generate KSK/ZSK key pairs (manual mode) |
| `dnssec-settime` | Set key lifecycle timestamps (Publish/Activate/Inactive/Delete) |
| `dnssec-signzone` | Manually sign a zone file |
| `dnssec-dsfromkey` | Generate DS record from a DNSKEY file |
| `dnssec-verify` | Offline verification of a signed zone file |
| `named-compilezone` | Convert raw zone format to human-readable text |
| `rndc signing -list` | Check signing progress for a zone |
| `rndc nta` | Add/list Negative Trust Anchors |
| `rndc dnssec -checkds` | Manually signal DS published/withdrawn to BIND |
| `nsec3hash` | Compute NSEC3 hash for a name (debugging) |
| `dnssec-ksr` | Generate Key Signing Requests for offline KSK |

---

## Entities and Concepts Mentioned

**Entities**: [[Internet-Systems-Consortium]], [[DNSViz]]
**Concepts**: [[dnssec-in-bind]] (primary — `dnssec-policy`, signing, rollover, inline signing), [[chain-of-trust]] (12-step validation, trust anchors, DS mechanics), [[dnssec-internals]] (CSK, key metadata, NSEC3 mechanics), [[dnssec-debugging]] (SERVFAIL triage, `ad` flag, security lameness, NTA, debug logging)

---

## Evidence Quality Notes

- **Highest authority**: This is the official ISC guide written by the BIND maintainers.
- **Version note**: Written for BIND 9.18+; most content applies from 9.16+. Pre-9.16 `auto-dnssec` content is present but explicitly deprecated in favor of `dnssec-policy`.
- **Completeness**: Covers signing, validation, troubleshooting, key management, and advanced topics. The single most comprehensive BIND DNSSEC document in the wiki.
