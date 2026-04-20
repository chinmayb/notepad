---
title: "DNSSEC Debugging"
type: concept
tags: ["dnssec", "debugging", "tools", "dig", "delv", "dnsviz", "edns0", "clock-skew", "servfail", "security-lameness", "nta"]
created: 2026-04-13
updated: 2026-04-14
source_count: 2
confidence: high
status: active
related: ["[[dnssec-internals]]", "[[chain-of-trust]]", "[[dnssec-in-bind]]", "[[dnssec-vs-dns-resolution]]", "[[DNSViz]]"]
open_questions:
  - "What is the exact sequence of checks `delv` performs and in what order does it report failures?"
  - "What does DNSViz show that `dig`/`delv` don't in terms of intermediate delegation state?"
  - "How do you test DNSSEC validation without affecting production traffic?"
contradictions: []
---

# DNSSEC Debugging

> **Domain**: DNS Operations / Debugging
> **Also known as**: DNSSEC troubleshooting, DNSSEC validation debugging

---

## Definition

DNSSEC debugging is the process of diagnosing failures in the DNSSEC chain of trust — most commonly manifesting as `SERVFAIL` responses from a validating resolver, or as a zone that appears signed but whose signatures are not trusted.

The key insight: **DNSSEC failures are binary** — the chain either validates or it doesn't. A resolver never returns an unvalidated answer; it returns SERVFAIL. This makes failures highly visible but the root cause can be at any link in the chain.

---

## Step 0: Determine If It's a DNSSEC Failure (Not a Network Problem)

```bash
# Step 1: Observe the failure
dig @10.53.0.3 www.example.org A
# → status: SERVFAIL

# Step 2: Re-run with +cd (Checking Disabled) to bypass DNSSEC validation
dig @10.53.0.3 www.example.org A +cd
# → status: NOERROR  ← if this succeeds, it's a DNSSEC validation failure
# → status: SERVFAIL ← if this also fails, it's not a DNSSEC problem

# Step 3: Determine query path (authoritative vs recursive)
# Look for 'aa' flag = talking to authoritative server
# Look for 'ra' flag = talking to recursive server (cached answers possible)
```

Source: [[dnssec-guide-bind9]]

---

## The `ad` Flag: The Validation Signal

In `dig` output, the **`ad` (Authenticated Data) flag** in the header flags section is the signal that DNSSEC validation succeeded:

```
;; flags: qr rd ra ad; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 1
                  ^^
                  ad = DNSSEC validated
```

| Flag present | Meaning |
|-------------|---------|
| `ad` in flags | DNSSEC validation passed — answer is authenticated |
| No `ad` flag | Not validated (zone unsigned, or validation disabled, or failure) |
| `SERVFAIL` + no `ad` | DNSSEC validation failed |
| `SERVFAIL` + `+cd` succeeds | Confirmed DNSSEC validation failure |

---

## Tool Overview

| Tool | When to Use | What It Does |
|------|-------------|--------------|
| `dig +dnssec` | First tool | Shows RRSIG records alongside answers; sets DO bit |
| `dig +cd` | Bypass validation | Checking Disabled — lets you see raw zone data even when validation fails |
| `dig +multiline` | Human-readable detail | Shows DNSKEY algorithm, key tag, RRSIG timestamps in readable form |
| `delv` | Chain validation trace | Performs full chain validation, explains where/why it fails with `+vtrace` |
| `dnssec-verify` | Offline zone file check | Verifies a signed zone file without DNS queries |
| `dnsviz` / [[DNSViz]] | Visual chain debugging | Renders full delegation chain graphically; highlights breaks. **ISC-endorsed as the single most useful DNSSEC diagnostic tool.** |
| `named-checkzone` | Pre-flight syntax check | Catches zone file errors before BIND loads them |
| `nsec3hash` | NSEC3 debugging | Computes the NSEC3 hash for a name to verify NSEC3 chain integrity |

---

## Common Failure Modes

| Symptom | Likely Cause | First Check |
|---------|-------------|-------------|
| `SERVFAIL` on all queries for a zone | Broken chain of trust | `dig +cd` to bypass; compare DS key tag vs active DNSKEY tag |
| `SERVFAIL` only after key change | DS/KSK mismatch (security lameness) | `dig <zone> DS` vs `dig +cd <zone> DNSKEY +multiline` — compare key tags |
| `SERVFAIL` appearing gradually over time | RRSIG expiry | `dig +dnssec <zone> SOA +multiline` — check RRSIG expiration timestamps |
| DNSSEC working locally but failing in production | EDNS0 stripped by network device | Test with `dig +edns=0` vs plain `dig`; firewall/LB may be filtering |
| Intermittent `SERVFAIL` on correctly-signed zone | Clock skew | Check NTP sync; RRSIG timestamps can appear expired with skewed clocks |
| Zone appears signed but parent doesn't trust it | DS not submitted to parent yet | `dig <zone> DS` — if empty, DS hasn't propagated |
| `SERVFAIL` for NXDOMAIN replies only | NSEC/NSEC3 misconfiguration | Check NSEC3 parameters; `nsec3hash` to verify hash chain |
| All DNS lookups `SERVFAIL` on a resolver | Invalid trust anchor | `delv` trace will show "no DNSKEY matching DS" for root |
| Algorithm-related `SERVFAIL` | Algorithm not supported by resolver | Check algorithm number in DNSKEY; RFC 8624 for support matrix |

---

## Security Lameness (Most Common Production Failure)

**Security lameness** = the parent DS record points to a key that doesn't exist in the child zone (or vice versa). The entire child zone becomes unresolvable.

**BIND log signature**:
```
named[126063]: validating example.net/DNSKEY: no valid signature found (DS)
named[126063]: no valid RRSIG resolving 'example.net/DNSKEY/IN': 10.53.0.2#53
named[126063]: broken trust chain resolving 'www.example.net/A/IN': 10.53.0.2#53
```

**Triage procedure**:
```bash
# 1. Get the DS record from the parent zone
dig @10.53.0.3 example.net. DS
# Note the key tag (e.g., DS 14956 ...)

# 2. Get the DNSKEY records from the child zone (bypass validation)
dig @10.53.0.3 example.net. DNSKEY +dnssec +cd +multiline
# Look for DNSKEY records; note their key IDs

# 3. Compare: DS key tag MUST match one of the DNSKEY key IDs
# Mismatch = security lameness. Resolution: re-upload DS records to parent.
```

Source: [[dnssec-guide-bind9]]

---

## Clock Skew: The Silent Killer

RRSIG records have two timestamps:
```
ftp.isc.org. 300 IN RRSIG A 13 3 300 (
    20200401191851   ← Signature Expiration (UTC April 1, 2020 19:18:51)
    20200302184340   ← Signature Inception (UTC March 2, 2020 18:43:40)
    27566 isc.org.
    ...
```

The resolver's clock **must fall within** this window. If the clock is:
- Too far in the **future**: `RRSIG has expired` → SERVFAIL
- Too far in the **past**: `RRSIG validity period has not begun` → SERVFAIL

**BIND log messages for clock errors**:
```
# Clock too far in future (RRSIG appears expired):
validating example.com/DNSKEY: verify failed due to bad signature (keyid=19036): RRSIG has expired

# Clock too far in past (RRSIG not yet valid):
validating example.com/DNSKEY: verify failed due to bad signature (keyid=4521): RRSIG validity period has not begun
```

**Fix**: Ensure NTP is running and synced. See also: [[dnssec-in-bind]] Prerequisites.

---

## BIND Debug Logging for DNSSEC

> ⚠️ Enable only temporarily — dramatically increases server load.

```nginx
# named.conf — log DNSSEC validation chain at debug level 3
logging {
    channel dnssec_log {
        file "/var/log/dnssec.log";
        severity debug 3;
        print-category yes;
    };
    category dnssec { dnssec_log; };
};
```

Example successful validation log output (chain for `ftp.isc.org`):
```
validating ./DNSKEY: verify rdataset (keyid=20326): success
validating ./DNSKEY: marking as secure (DS)
validating org/DS: verify rdataset (keyid=33853): success
validating org/DNSKEY: verify rdataset (keyid=9795): success
validating isc.org/DS: verify rdataset (keyid=33209): success
validating isc.org/DNSKEY: verify rdataset (keyid=7250): success
validating ftp.isc.org/A: verify rdataset (keyid=27566): success
validating ftp.isc.org/A: marking as secure, noqname proof not needed
```

---

## Negative Trust Anchors (NTA): Break-Glass for Misconfigured Third-Party Zones

When a third-party zone has broken DNSSEC and is causing SERVFAIL, use an NTA as a temporary workaround:

```bash
# Add NTA for 1 hour (default)
rndc nta example.com
# → "Negative trust anchor added: example.com/_default, expires 19-Mar-2020 19:57:42"

# List active NTAs
rndc nta -dump

# BIND polls the zone every 5 minutes; NTA auto-expires when zone validates correctly
```

NTAs are a **safety valve** — not a long-term solution. They bypass validation for the entire subtree of the named zone.

Source: [[dnssec-guide-bind9]]

---

## `dig` DNSSEC Cheatsheet

```bash
# Check if a zone is signed (look for DNSKEY record)
dig +dnssec example.com DNSKEY +multiline

# Check if a record is validated (look for 'ad' flag)
dig @resolver example.com A +dnssec

# Bypass validation to see raw zone data
dig @resolver example.com A +cd +dnssec

# Check DS record at parent
dig example.com DS

# Check RRSIG expiration dates (use +multiline for timestamps)
dig @ns1.example.com example.com SOA +dnssec +multiline

# Query authoritative directly, bypass resolver cache
dig @ns1.example.com example.com DNSKEY +dnssec +multiline

# Confirm DNSSEC validation failure vs network failure
dig @resolver www.dnssec-failed.org A           # → SERVFAIL (with validation)
dig @resolver www.dnssec-failed.org A +cd       # → NOERROR (without validation)
```

---

## `delv` Cheatsheet

*To be expanded from `delv` manpage source.*

```bash
# Full chain validation with trace output
delv +vtrace example.com A

# Use a specific trust anchor (for testing before DS is in parent)
delv @192.168.1.13 -a /tmp/example.key +root=example.com example.com. SOA +multiline
# Look for: "; fully validated"
```

---

## Debugging Methodology (Ordered Triage)

1. **Is it a DNSSEC failure?** — `dig +cd` to bypass validation. If `+cd` succeeds and plain query fails → yes, DNSSEC problem.
2. **Determine query path** — `aa` flag = authoritative; `ra` flag = recursive. Route to the right server.
3. **Check the parent's DS** — `dig <zone> DS` — is it present? Does the key tag match an active DNSKEY?
4. **Check RRSIG validity windows** — `dig @auth <zone> SOA +dnssec +multiline` — are signatures expired or future-dated?
5. **Check for clock skew** — Verify NTP sync on both resolver and authoritative server.
6. **Check EDNS0** — Is DNSSEC traffic making it through firewalls? `dig +edns=0` vs plain `dig`.
7. **Enable BIND debug logging** — `category dnssec { severity debug 3; }` to trace the full validation chain.
8. **Visualize with [[DNSViz]]** — If above steps don't reveal the issue, paste the domain into dnsviz.net for a full graphical chain analysis.

---

## Mentioned In

- [[dnssec-and-bind9-isc]] — DNSViz endorsed as primary tool; clock skew and EDNS0 as failure modes
- [[dnssec-guide-bind9]] — primary source; security lameness, SERVFAIL triage, `ad` flag, NTA, debug logging, clock skew log messages

---

*Further expansion pending: `dig and delv.md`, `Ubuntu Manpage delv.md`*
