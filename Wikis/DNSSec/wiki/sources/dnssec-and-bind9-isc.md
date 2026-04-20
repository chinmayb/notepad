---
title: "DNSSEC and BIND 9 (ISC)"
type: source
tags: ["dnssec", "bind", "isc", "kasp", "edns0", "dnsviz", "operational"]
created: 2026-04-14
updated: 2026-04-14
source_count: 1
confidence: high
status: active
related: ["[[dnssec-in-bind]]", "[[dnssec-internals]]", "[[dnssec-debugging]]", "[[Internet-Systems-Consortium]]"]
open_questions:
  - "Is the DNSSEC Multi-Signer Model (RFC 8901) relevant to the target environment?"
contradictions: []
---

# Source: DNSSEC and BIND 9 (ISC)

> **Raw file**: `raw/dnssec0guide/DNSSEC and BIND 9.md`
> **Source URL**: https://www.isc.org/dnssec/
> **Author**: [[Internet-Systems-Consortium]] (ISC)
> **Published**: 2025-12-30
> **Type**: Official overview / resource portal
> **Coverage axes**: BIND operations (Axis 5), Debugging tooling (Axis 6 pointer)

---

## Summary

Short official page from ISC (the maintainers of BIND 9) serving as both a conceptual introduction to DNSSEC and a curated index of authoritative documentation, presentations, and tools. Not a deep technical document — its value is (1) endorsing the canonical tooling and documentation hierarchy, (2) surfacing operational prerequisites often overlooked, and (3) confirming `dnssec-policy` / KASP as the current recommended approach.

---

## Key Claims

1. **KASP (`dnssec-policy`) supersedes `auto-dnssec`** — The source explicitly notes that older documentation (e.g., the 2020 webinar series, NANOG50 talk) "pre-dated the current Key and Signing Policy feature (KASP)." All current ISC guidance favors `dnssec-policy`.

2. **EDNS0 is a hard prerequisite** — DNSSEC requires EDNS0 to support larger DNS message sizes and for the `DO` (DNSSEC OK) EDNS header bit. Network devices that silently drop EDNS0 packets will silently break DNSSEC.

3. **Clock skew causes SERVFAIL** — DNSSEC is more sensitive to time sync issues than plain DNS. System clocks that are significantly out of sync can cause valid RRSIG records to appear expired, producing spurious SERVFAIL. NTP sync is a DNSSEC operational requirement.

4. **DNSViz is the primary diagnostic tool** — Direct ISC quote: *"DNSVIZ remains the single most useful diagnostic tool for DNSSEC."* Strong institutional endorsement of this tool for chain visualization and debugging.

5. **Multi-Signer Model is an active development area** — ISC references multiple presentations on the DNSSEC Multi-Signer Model (RFC 8901), suggesting it is operationally relevant for BIND deployments with multiple authoritative operators.

6. **Secondary nameservers must also support DNSSEC** — When hosting DNSSEC-signed zones, all secondaries must have DNSSEC enabled. A secondary that doesn't support it will serve unsigned data.

7. **BIND 9 supports both roles** — Validating resolver and authoritative signing server. All BIND 9 versions are DNSSEC-capable.

---

## Direct Quotes Worth Preserving

> "DNSSEC doesn't provide a secure tunnel; it doesn't encrypt or hide DNS data. It was designed with backwards compatibility in mind."

> "DNSVIZ remains the single most useful diagnostic tool for DNSSEC."

> "DNSSEC is more sensitive to time issues (i.e. system clocks being really out of sync) than plain DNS; make sure your system clocks are reasonably accurate."

> "Note that this [2020 webinar series] pre-dated the current Key and Signing Policy feature (KASP)."

---

## Resource Index (from source)

| Resource | Type | Relevance |
|----------|------|-----------|
| BIND 9 ARM (ReadTheDocs) | Official reference | Primary BIND config reference |
| ISC KB: DNSSEC Key and Signing Policy | KB article | KASP config detail |
| ISC KB: DNSSEC Key States | KB article | Key lifecycle states |
| DNSViz (dnsviz.net) | Diagnostic tool | Chain visualization — **highest priority debug tool** |
| DELV (ISC KB) | Tool reference | Validation chain checker |
| BIND 9 DNSSEC Guide (appendix) | Long-form guide | "Why" and "how" — our `raw/dnssec0guide/DNSSEC Guide` source |
| 2020 webinar series (10 parts) | Video | Pre-KASP; still useful for theory |
| NANOG50 talk (Alan Clegg) | Presentation | Pre-KASP deployment walkthrough |

---

## Entities and Concepts Mentioned

**Entities**: [[Internet-Systems-Consortium]], [[DNSViz]]
**Concepts**: [[dnssec-in-bind]] (KASP, `dnssec-policy`, `auto-dnssec`), [[dnssec-internals]] (EDNS0, DO bit, RRSIG, DNSKEY, DS, NSEC), [[dnssec-debugging]] (DNSViz, DELV, clock skew as failure mode)

---

## Evidence Quality Notes

- **High authority**: ISC maintains BIND 9. This is the primary source of record for BIND DNSSEC guidance.
- **Not a tutorial**: This page is a portal/index, not a technical walkthrough. It points to richer documents (the BIND Guide, KASP KB article) that are staged in `raw/`.
- **Date**: Published 2025-12-30 — current.
