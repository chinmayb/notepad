# Wiki Index

> Master catalog of all pages in this knowledge base.
> **Maintained by the LLM agent.** Updated on every ingest, query filing, and lint pass.
> Last updated: 2026-04-14 (ingest: cs161-dnssec)

---

## How to Use This Index

**New to DNSSEC? Start with [[dnssec-learning-roadmap]]**, not here — it gives the first-principles reading order through the concept pages below. This index is a reference catalog for looking something up once you already know what you're after.

- Start here every session to orient to what exists
- Entries follow the format: `[[PageTitle]] — one-line description (sources: N)`
- Sections mirror the wiki folder structure

---

## Overview

| Page | Description |
|------|-------------|
| [[overview]] | Top-level synthesis; domain = DNSSEC; 11 sources staged, 8 ingested |

---

## Sources
*One page per ingested raw source.*

| Page                                | Raw File                                                              | Type                         | Axes                                                                              |
| ----------------------------------- | --------------------------------------------------------------------- | ---------------------------- | --------------------------------------------------------------------------------- |
| [[01-learn-dnssec]]                 | `raw/01-learn-dnssec.md`                                              | Personal note / goals        | All 6                                                                             |
| [[how-does-dnssec-work-cloudflare]] | `raw/dns-sec-cloudflare&text/How Does DNSSEC Work?.md`                | Tutorial article             | Internals, chain of trust, deployment                                             |
| [[kaminsky-exploit]]                | `raw/dnssec0guide/DNS Vulnerabilities...md`                           | Retrospective article        | First principles (Axis 3)                                                         |
| [[dnssec-and-bind9-isc]]            | `raw/dnssec0guide/DNSSEC and BIND 9.md`                               | Official ISC portal          | BIND ops prerequisites, tooling endorsements                                      |
| [[dnssec-guide-bind9]]              | `raw/dnssec0guide/DNSSEC Guide — BIND 9 9.21.22-dev documentation.md` | Official ISC long-form guide | BIND ops, signing, validation, rollover, debugging                                |
| [[dnssec-kasp-policy]]              | `raw/dnssec0guide/DNSSEC Key and Signing Policy.md`                   | ISC KB article               | `dnssec-policy` migration, KSK rollover signals, KASP states                      |
| [[rfc-6781]]                        | `raw/dns-sec-rfc/RFC 6781 DNSSEC Operational Practices, Version 2.md` | IETF Informational RFC       | Rollover procedures, timing terminology, emergency rollover, stand-by keys        |
| [[cs161-dnssec]]                    | `raw/dns-sec-cloudflare&text/DNSSEC.md`                               | Academic textbook chapter    | Chain of trust walkthrough, `dig` query workflow, DNSKEY flags, OPT pseudosection |

---

## Entities
*People, organizations, products, projects.*

| Page | Description | Status |
|------|-------------|--------|
| [[Cloudflare]] | DNS infrastructure company; DNSSEC educator and operator | Active |
| [[Dan-Kaminsky]] | Security researcher; discovered the 2008 DNS cache poisoning exploit | Active |
| [[Internet-Systems-Consortium]] | Nonprofit; maintains BIND 9; authors of KASP / `dnssec-policy` | Active |
| [[DNSViz]] | Visual DNSSEC chain-of-trust diagnostic tool; ISC-endorsed | Active |

---

## Concepts
*Ideas, frameworks, techniques, mental models, terms.*

| Page | Description | Status |
|------|-------------|--------|
| [[chain-of-trust]] | Cryptographic hierarchy root→answer; DS mechanics; KSK rollover; first-principles motivation | Active |
| [[dnssec-internals]] | RRsets, record types (DNSKEY/RRSIG/DS/NSEC/NSEC3), KSK/ZSK rationale (operational + storage-tier + when-not-to), EDNS0/DO bit | Active |
| [[dns-cache-poisoning]] | The attack DNSSEC prevents; Kaminsky mechanics; port randomization vs cryptographic fix | Active |
| [[dnssec-debugging]] | Tools (dig, delv, DNSViz), failure modes incl. clock skew + EDNS0 stripping | Active |
| [[dnssec-in-bind]] | BIND 9 config, `dnssec-policy` vs `auto-dnssec`, operational prerequisites | Active |
| [[dnssec-vs-dns-resolution]] | What changes at each resolution step under DNSSEC; full 6-step `dig` walkthrough | Active |

---

## Synthesis
*Cross-source analysis, comparisons, structured arguments.*

| Page | Description |
|------|-------------|
| [[dnssec-learning-roadmap]] | Curriculum map: reading order, source coverage by axis, known gaps |

---

## Queries
*Filed answers to questions worth keeping.*

| Page | Question | Date |
|------|----------|------|
| [[lint-2026-04-10]] | Initial wiki bootstrap lint | 2026-04-10 |

---

## Stats

| Metric | Count |
|--------|-------|
| Total wiki pages | 23 |
| Total raw sources staged | 11 |
| Raw sources ingested | 8 |
| Entities | 4 |
| Concepts | 6 (6 active, 0 stubs) |
| Source summaries | 8 |
| Synthesis pages | 1 |
| Filed queries | 1 |
| Open contradictions | 1 (NSEC3 iterations: RFC 6781 vs RFC 9276 — resolved in favor of RFC 9276) |
| Orphan pages | 0 |

---

*To add a page to this index, append to the appropriate section above with the standard format.*
