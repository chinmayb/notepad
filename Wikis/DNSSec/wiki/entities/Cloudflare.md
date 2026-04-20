---
title: "Cloudflare"
type: entity
tags: ["dns", "cdn", "dnssec", "infrastructure"]
created: 2026-04-13
updated: 2026-04-13
source_count: 1
confidence: high
status: active
related: ["[[how-does-dnssec-work-cloudflare]]", "[[chain-of-trust]]"]
entity_type: organization
open_questions: []
contradictions: []
---

# Cloudflare

> **Type**: Organization
> **Domain**: DNS infrastructure, CDN, network security

---

## Overview

Cloudflare is a major Internet infrastructure company providing DNS hosting, CDN, DDoS protection, and security services. Relevant to this wiki as the author of DNSSEC educational content and as a party involved in efforts to automate DNSSEC DS record management between DNS operators and registrars.

---

## Key Facts

| Attribute | Value |
|-----------|-------|
| Type | Private company |
| Known for | DNS hosting, CDN, 1.1.1.1 resolver, DNSSEC education |
| DNSSEC role | DNS operator + educator + automation protocol proposer |

---

## Relevance to DNSSEC

- Author of the `How Does DNSSEC Work?` tutorial (one of the clearest operator-facing explanations)
- Co-authored an Internet Draft with CIRA (.ca registry) proposing a protocol for DNS operators to automate DS record updates with registrars — would eliminate manual DS record submission
- Operates a DNSSEC-enabled resolver (1.1.1.1) and offers DNSSEC signing for customer zones

---

## Mentioned In

- [[how-does-dnssec-work-cloudflare]] — primary source; Cloudflare is the author

---

*Relevant as source author and DNSSEC automation advocate.*
