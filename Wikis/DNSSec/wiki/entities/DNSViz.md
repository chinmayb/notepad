---
title: "DNSViz"
type: entity
tags: ["tool", "dnssec", "debugging", "visualization", "dns"]
created: 2026-04-14
updated: 2026-04-14
source_count: 1
confidence: high
status: active
related: ["[[dnssec-debugging]]", "[[chain-of-trust]]", "[[dnssec-in-bind]]"]
open_questions:
  - "Does DNSViz work with internal/private zones or only publicly resolvable domains?"
  - "What is the difference between DNSViz online (dnsviz.net) and the open-source CLI tool?"
contradictions: []
---

# DNSViz

> **Type**: Diagnostic tool
> **URL**: https://dnsviz.net/
> **Role**: Visual DNSSEC chain-of-trust analyzer

---

## Description

DNSViz is a tool for visualizing the status of the DNSSEC authentication chain for a given domain. It produces a graphical representation of the full chain — from root trust anchor down through TLD delegation to the authoritative zone — making it easy to spot broken links, expired signatures, algorithm mismatches, and misconfigured DS records.

ISC endorses it as *"the single most useful diagnostic tool for DNSSEC"* (direct quote, [[dnssec-and-bind9-isc]]).

---

## Key Capabilities

- Full chain-of-trust visualization
- Flags broken/invalid RRSIG records
- Highlights missing or mismatched DS records
- Shows algorithm and key tag for each zone
- Available as both a web service and an open-source CLI tool

---

## Mentioned In

- [[dnssec-and-bind9-isc]] — endorsed as primary diagnostic tool
