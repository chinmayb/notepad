# CLAUDE.md — DNSSec overlay

Schema: see [../CLAUDE.md](../CLAUDE.md) for the full LLM Wiki pattern (frontmatter, INGEST/QUERY/LINT/EVOLVE, naming, cross-reference rules). This file only covers what's specific to DNSSec.

## Domain

DNSSEC (DNS Security Extensions) — the cryptographic chain-of-trust that lets a resolver validate that a DNS answer really came from the zone's authoritative server. Covers protocol internals (RRsets, DNSKEY/RRSIG/DS/NSEC/NSEC3), BIND 9 operational practice (`dnssec-policy`/KASP, key rollover), the history of the attack it defends against (Kaminsky cache poisoning), and the tooling used to debug it (`dig`, `delv`, DNSViz).

## Domain-specific framing

- **Entities** are the people, organizations, and tools in the DNSSEC ecosystem — e.g. Dan Kaminsky, Internet Systems Consortium, DNSViz, Cloudflare.
- **Concepts** are the protocol mechanics and operational practices — e.g. chain-of-trust, KSK/ZSK rollover, NSEC/NSEC3, `dnssec-policy`/KASP, cache poisoning.
- RFCs are treated as sources like any other and cited by number (e.g. RFC 6781, RFC 9276) — when RFCs disagree (a newer RFC supersedes operational guidance in an older one), that's a contradiction to flag, not silently resolve.

*Schema version: 1.0 | Created: 2026-04-10 | Vault: DNSSec*
