---
title: "Dan Kaminsky"
type: entity
tags: ["dns", "security", "researcher", "kaminsky"]
created: 2026-04-13
updated: 2026-04-13
source_count: 1
confidence: high
status: active
related: ["[[dns-cache-poisoning]]", "[[kaminsky-exploit]]", "[[chain-of-trust]]"]
entity_type: person
open_questions:
  - "What was Kaminsky's involvement in DNSSEC deployment after the 2008 disclosure?"
  - "What other DNS-related security research did Kaminsky do beyond the 2008 exploit?"
contradictions: []
---

# Dan Kaminsky

> **Type**: Person
> **Domain**: DNS Security, Vulnerability Research

---

## Overview

Dan Kaminsky (1979–2021) was an American security researcher best known for discovering a fundamental flaw in the DNS protocol in 2008 that enabled reliable, large-scale cache poisoning attacks. His responsible disclosure of this vulnerability — coordinating privately with DNS software vendors worldwide before public release — is credited with preventing what could have been one of the most damaging internet security incidents of the decade.

---

## Key Facts

| Attribute | Value |
|-----------|-------|
| Known for | Kaminsky DNS cache poisoning exploit (2008) |
| Disclosure style | Responsible — private coordinated disclosure to vendors before public |
| Impact | Triggered global coordinated patching of DNS resolver software |
| Legacy | Catalyzed renewed DNSSEC adoption; established internet-scale vulnerability response model |

---

## The 2008 DNS Disclosure

Kaminsky discovered that DNS resolvers could be reliably poisoned by flooding them with queries for uncached subdomains of a target domain — generating thousands of transaction ID guessing opportunities per second. This made cache poisoning practical at scale for the first time.

Rather than publish immediately, he:
1. Privately notified DNS software vendors (BIND, Microsoft, Cisco, etc.)
2. Coordinated with CERT and major tech companies
3. Allowed ~30 days for simultaneous global patching
4. Announced the vulnerability class publicly at BlackHat 2008 without initially releasing full technical details

This coordinated disclosure model became a template for handling high-severity internet infrastructure vulnerabilities.

---

## Significance to This Wiki

The Kaminsky exploit is the **"why DNSSEC"** story — the concrete historical event that turned DNSSEC from a theoretical good-practice into an urgent operational priority. Understanding it gives first-principles grounding for every DNSSEC design decision. See [[dns-cache-poisoning]] for the full attack mechanics.

---

## Mentioned In

- [[kaminsky-exploit]] — primary source documenting the exploit and its aftermath

---

*Entity page — Dan Kaminsky passed away in April 2021.*
