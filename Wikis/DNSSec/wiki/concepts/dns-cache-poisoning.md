---
title: "DNS Cache Poisoning"
type: concept
tags: ["dns", "security", "attack", "cache-poisoning", "kaminsky"]
created: 2026-04-13
updated: 2026-04-13
source_count: 1
confidence: high
status: active
related: ["[[chain-of-trust]]", "[[dnssec-internals]]", "[[dnssec-vs-dns-resolution]]", "[[dnssec-debugging]]"]
open_questions:
  - "Do modern resolvers still implement source port randomization by default, or has it been superseded?"
  - "Can DNS-over-TLS / DNS-over-HTTPS prevent cache poisoning independently of DNSSEC?"
  - "What is the current estimated probability of a successful Kaminsky-style attack against a port-randomized resolver?"
contradictions: []
---

# DNS Cache Poisoning

> **Domain**: DNS Security / Attack Taxonomy
> **Also known as**: DNS spoofing, DNS hijacking (via cache)

---

## Definition

DNS cache poisoning is an attack in which a malicious actor injects false DNS records into a recursive resolver's cache. Once poisoned, the resolver serves the attacker's fabricated answer to all clients — silently redirecting them to attacker-controlled IP addresses without any visible indication of compromise.

It is the attack class that the Kaminsky exploit (2008) made practically viable at scale, and the **primary threat** that DNSSEC is designed to prevent.

---

## Why DNS Caches Are the Target

Recursive resolvers cache answers to avoid querying authoritative servers on every request. This cache is:
- **Shared** — one resolver serves thousands or millions of users (ISP, enterprise)
- **Trusted implicitly** — clients accept cached answers without re-validating
- **Writable via spoofing** — before DNSSEC, any forged response with the right transaction ID would be accepted

Poisoning the cache once → all users of that resolver are affected until the TTL expires.

---

## The Classical Attack (Pre-Kaminsky)

The baseline cache poisoning attack required:
1. Know the target domain being queried
2. Wait for the resolver to issue a query for that domain
3. Race the real authoritative response with a forged response containing the correct 16-bit transaction ID
4. Win the race before the legitimate response arrives

**Why it was considered impractical**: step 3 required guessing a 16-bit transaction ID (65,536 possibilities) within a narrow timing window. Attacker had to be on-path or very lucky.

---

## The Kaminsky Breakthrough (2008)

[[Dan-Kaminsky]] discovered that the subdomain technique eliminates the timing constraint entirely:

```
Instead of querying: bank.com (might be cached, narrow window)
Query:               aaaaaa.bank.com, bbbbbb.bank.com, cccccc.bank.com ...
```

Each uncached subdomain forces the resolver to issue a fresh outbound query → fresh transaction ID window → fresh poisoning opportunity. The attacker can generate **thousands of attempts per second**, flooding the resolver with forged responses until one transaction ID matches.

**The injected forged answer** doesn't claim `aaaaaa.bank.com` — it poisons the *parent domain* `bank.com` in the AUTHORITY or ANSWER section of the forged response. The subdomain is just the vehicle to generate the query.

**Entropy pre-fix**: 16 bits = 65,536 possibilities
**Expected guesses to succeed**: ~32,768 on average
**Time to succeed**: seconds to minutes at flood rates

See [[kaminsky-exploit]] for the full step-by-step attack sequence.

---

## The Emergency Fix: Source Port Randomization

The coordinated patch released in July 2008 added **UDP source port randomization** to DNS queries:

- Old: fixed or predictable source port + random 16-bit transaction ID = 16 bits of entropy
- New: random source port (16-bit) + random transaction ID (16-bit) = ~32 bits of entropy

**32 bits = ~4 billion combinations** — brute-forcing this at flood rates takes hours to months, not seconds.

This is a *probabilistic mitigation*, not a cryptographic guarantee. A sufficiently patient or resourced attacker can still succeed, especially against resolvers with predictable port allocation.

---

## Why DNSSEC Is the Structural Solution

Port randomization raises the cost of attack. DNSSEC makes the attack **categorically impossible** against validating resolvers:

- The forged response must contain a valid RRSIG over the poisoned records
- The RRSIG requires the zone's private ZSK to generate
- The attacker does not have the private key
- A validating resolver will reject any response with a missing or invalid RRSIG
- **Transaction ID match is irrelevant** — the cryptographic check is independent of transport-layer identifiers

DNSSEC converts cache poisoning from a *probabilistic timing attack* into a *computationally infeasible cryptographic problem*.

---

## Downstream Impact of a Poisoned Cache

Once poisoned, all higher-layer security for affected users is compromised:

| Attack Goal | How Poisoned DNS Enables It |
|-------------|----------------------------|
| Credential theft | Redirect bank.com → fake login page |
| HTTPS bypass | Browser connects to attacker IP *before* TLS — attacker's cert is presented |
| Malware delivery | Redirect software update servers → serve malicious packages |
| Email interception | Poison MX records → reroute email through attacker's servers |
| Certificate mis-issuance | Poison DNS used by CAs for domain validation → trick CA into issuing cert for attacker |

The last entry is particularly significant: a poisoned DNS cache can undermine Certificate Authority domain validation, allowing an attacker to obtain a valid TLS certificate for a domain they don't control.

---

## What Cache Poisoning Does NOT Do

- Does not compromise the zone's authoritative server
- Does not affect users who query the authoritative server directly
- Does not persist after cache TTL expires (but re-attack is trivial)
- Does not affect encrypted DNS (DoH/DoT) at the transport layer — but the resolver itself can still be poisoned at the application layer

---

## Mentioned In

- [[kaminsky-exploit]] — primary source; Kaminsky's specific technique and impact

---

*To be further expanded from: CS161 DNSSEC source (covers DNSSEC as the fix), RFC 4033 (formal threat model)*
