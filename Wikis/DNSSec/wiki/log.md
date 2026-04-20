# Wiki Log

> Append-only chronological record of all wiki operations.
> Format: `## [YYYY-MM-DD] <operation> | <title>`
> Grep-friendly: `grep "^## \[" wiki/log.md | tail -10`

---

## [2026-04-20] lint | Wiki Review — Metadata & Consistency Fixes

- **Operation**: lint + fix
- **Pages created**:
  - `wiki/queries/lint-2026-04-10.md` — bootstrap lint stub (file was missing; referenced in index and log)
- **Pages updated**:
  - `wiki/overview.md` — source_count 1→8; added `[[dns-cache-poisoning]]` to related; removed stale KASP open question; removed duplicate ❌ rows for CS161 and RFC 6781 from source inventory; updated date 2026-04-14→2026-04-20
  - `wiki/synthesis/dnssec-learning-roadmap.md` — source inventory table updated (8 sources now ✅); Bottom Line paragraph rewritten to reflect current state; source_count 5→8; updated date 2026-04-13→2026-04-20
  - `wiki/index.md` — overview description fixed ("7 ingested"→"8 ingested"); Cloudflare status "Stub"→"Active"
  - `wiki/entities/Cloudflare.md` — removed vestigial "Stub entity" footer (frontmatter already active)
- **Entities touched**: Cloudflare (status aligned)
- **Contradictions found**: none
- **Issues fixed**: 8 (1 missing file, 3 stale metadata, 2 duplicate rows, 1 stale open question, 1 entity status mismatch)

---

## [2026-04-14] ingest | DNSSEC — CS161 UC Berkeley Textbook

- **Operation**: ingest
- **Source**: `raw/dns-sec-cloudflare&text/DNSSEC.md`
- **Source type**: Academic textbook chapter — UC Berkeley CS161 Computer Security course; pedagogical focus, `dig`-based walkthrough with real output
- **Pages created**: (source summary was written in prior session)
  - `wiki/sources/cs161-dnssec.md` — source summary with 14 key claims (previously created)
- **Pages updated**:
  - `wiki/concepts/dnssec-vs-dns-resolution.md` — **stub → active** (confidence: speculative → high): full 6-step `dig` walkthrough for `eecs.berkeley.edu` with per-step records returned and verification logic; OPT pseudo-record / ADDITIONAL count +1 explanation; DNSKEY flag values table (256=ZSK, 257=KSK, SEP bit explanation); chain visualization diagram; offline signing motivation (online signing = DoS vulnerability); NSEC pre-computed ranges explanation; "What Stays the Same" and caching sections retained; source_count 0→2; status stub→active; confidence speculative→high
  - `wiki/concepts/dnssec-internals.md` — added: DNSKEY Flag Values section (flag 256=ZSK, 257=KSK, SEP bit, example dig output); OPT Pseudo-Record expanded (ADDITIONAL count +1 anomaly with example dig output block; `flags: do` explanation); offline signing DoS motivation clarified with explicit attacker model; tags updated (dnskey-flags, opt-record, offline-signing); source_count 5→6; Mentioned In updated
  - `wiki/concepts/chain-of-trust.md` — added: "KSK as Parent, ZSK as Child Within One Name Server" section (CS161 pedagogical framing; inter-zone vs intra-zone endorsement comparison table; CS161 direct quote); source_count 4→5; tags updated; Mentioned In updated
  - `wiki/index.md` — added `cs161-dnssec` source entry; updated `dnssec-vs-dns-resolution` status stub→active; updated stats (23 total pages, 8 ingested, 6 active concepts, 0 stubs)
  - `wiki/overview.md` — CS161 source marked ✅ in source inventory; concept table updated
- **Entities touched**: none new
- **Concepts touched**: dnssec-vs-dns-resolution (stub→active, heavily rewritten), dnssec-internals (DNSKEY flags + OPT record + offline signing), chain-of-trust (KSK/ZSK pedagogy section added)
- **Contradictions found**: none — CS161 content is consistent with all existing wiki sources
- **Notes**: CS161 uniquely contributes: (1) the most concrete step-by-step query walkthrough in the wiki — each delegation step now has the exact records returned and the verification logic; (2) the KSK/ZSK "parent-child within one name server" framing elegantly unifies inter-zone and intra-zone trust delegation into a single mental model; (3) the OPT pseudo-record ADDITIONAL count +1 anomaly is now explained — this is a common source of confusion when reading `dig` output; (4) DNSKEY flag values (256/257, SEP bit) are now explicitly documented with an example; (5) offline signing's DoS motivation is now precise — the attacker model is "flood with NXDOMAIN queries → exhaust CPU on on-demand signing."

---

## [2026-04-14] ingest | RFC 6781 — DNSSEC Operational Practices, Version 2

- **Operation**: ingest
- **Source**: `raw/dns-sec-rfc/RFC 6781 DNSSEC Operational Practices, Version 2.md`
- **Source type**: IETF Informational RFC (December 2012) — operational guidance for zone administrators; obsoletes RFC 4641
- **Pages created**:
  - `wiki/sources/rfc-6781.md` — source summary with 20 key claims, formal timing definitions, rollover procedures, emergency rollover options, stand-by keys, NSEC3 contradiction noted
- **Pages updated**:
  - `wiki/concepts/dnssec-internals.md` — added: formal signature timing terminology section (validity period / publication period / key effectivity period / Re-Sign Period / Refresh Period) with BIND mapping; SOA expiry guidance; NSEC3 iterations contradiction callout; source_count 4→5; contradictions[] updated; tags updated; open_questions pruned
  - `wiki/concepts/chain-of-trust.md` — added: Emergency KSK Rollover section (three options: keep chain / go Bogus / go Insecure); Stand-By Keys section; source_count 3→4; open_questions pruned; tags updated
  - `wiki/concepts/dnssec-in-bind.md` — added: SOA expiry alignment section (SOA expire ≈ 1/3–1/4 signature validity); DNS operator change procedure (cooperating vs non-cooperating); source_count updated; Mentioned In updated
  - `wiki/index.md` — added `rfc-6781` source entry; stats updated (22 total pages, 7 ingested, 1 open contradiction)
  - `wiki/overview.md` — source inventory updated (rfc-6781 marked ✅)
- **Entities touched**: none new
- **Concepts touched**: dnssec-internals, chain-of-trust, dnssec-in-bind (all updated)
- **Contradictions found**: ⚠️ **RFC 6781 §5.3.2 vs RFC 9276**: RFC 6781 recommends 100 NSEC3 iterations; RFC 9276 (2022) supersedes with 0 iterations. Flagged in `dnssec-internals.md` frontmatter and inline. Wiki guidance follows RFC 9276.
- **Notes**: This RFC is the canonical source for the rollover procedures that BIND 9's KASP implements (Pre-Publish ZSK, Double-KSK). Key new additions to wiki: (1) formal timing terminology is now defined — these terms appear in BIND config parameters without explanation elsewhere; (2) emergency rollover has three options (keep chain / go Bogus / go Insecure) — the "go Insecure" path via removing the DS is a valid last resort not previously documented; (3) stand-by keys are a pre-positioning strategy for fast recovery; (4) DNS operator migration procedure for the non-cooperating case (zone must go insecure). The NSEC3 iterations contradiction with RFC 9276 is the wiki's first flagged contradiction between sources.

---

## [2026-04-14] ingest | DNSSEC Key and Signing Policy (ISC KB)

- **Operation**: ingest
- **Source**: `raw/dnssec0guide/DNSSEC Key and Signing Policy.md`
- **Source type**: ISC knowledge base article — migration and rollover operational guide for `dnssec-policy`
- **Pages created**:
  - `wiki/sources/dnssec-kasp-policy.md` — source summary with 14 key claims, config snippets, evidence quality notes
- **Pages updated**:
  - `wiki/concepts/dnssec-in-bind.md` — added: exact `default` policy parameter table (dnskey-ttl, signatures-refresh, etc.); inline-signing version gate table (9.18 vs 9.20); 9.20 upgrade pitfall with 3 mitigations; KSK/CSK rollover two `rndc -checkds` signals section; `parental-agents` config snippet; `auto-dnssec` status upgraded from "deprecated" to "**removed in 9.19.16+**" with warning callout; migration from `auto-dnssec` procedure section; Mentioned In updated; source_count 2→3; open_questions pruned (now answered)
  - `wiki/concepts/dnssec-internals.md` — added: KASP key states section (hidden/rumoured/omnipresent/unretentive) with per-RR-type tracking table and `.state` file explanation; `DSState: rumoured stuck` fix documented; source_count 3→4; tags updated
  - `wiki/index.md` — added `dnssec-kasp-policy` source entry; stats updated (21 total pages, 6 ingested)
  - `wiki/overview.md` — source inventory updated (dnssec-kasp-policy marked ✅)
- **Entities touched**: none new
- **Concepts touched**: dnssec-in-bind (heavily updated), dnssec-internals (KASP states added)
- **Contradictions found**: none; `auto-dnssec` status upgraded from "deprecated" to "removed" — this is a clarification, not a contradiction (the BIND guide said deprecated; this source confirms it was then removed in 9.19.16)
- **Notes**: Key contributions from this source: (1) `auto-dnssec` is now confirmed REMOVED (not just deprecated) in BIND 9.20+; (2) `default` policy parameters are now precisely documented (14-day sig validity, 5-day refresh, 1-hour TTL, 90-day key purge); (3) the two `rndc dnssec -checkds` signals are the critical operational step for KSK/CSK rollover that was missing from the wiki; (4) KASP four-state machine (hidden/rumoured/omnipresent/unretentive) is now documented — this is the runtime propagation tracking layer, distinct from key lifecycle timestamps; (5) inline-signing version gate between 9.18 and 9.20 is now captured.

---

## [2026-04-14] ingest | DNSSEC Guide — BIND 9 (Official)

- **Operation**: ingest (completion of steps 1–8)
- **Source**: `raw/dnssec0guide/DNSSEC Guide — BIND 9 9.21.22-dev documentation.md`
- **Source type**: Official ISC long-form guide — BIND 9 DNSSEC reference (written by Josh Kuo, DeepDive Networking)
- **Pages created**:
  - `wiki/sources/dnssec-guide-bind9.md` — full source summary with key claims, config recipes, tool reference table
- **Pages updated**:
  - `wiki/concepts/dnssec-in-bind.md` — fully rewritten (stub→active, confidence: high): signing with `dnssec-policy default`, inline signing, `rndc` workflow, validation (`dnssec-validation auto`), three KSK rollover methods, full BIND tool reference table; source_count 1→2
  - `wiki/concepts/dnssec-debugging.md` — fully rewritten (stub→active, confidence: high): `ad` flag as validation signal, security lameness triage (DS/DNSKEY key tag comparison), clock skew log messages, NTA (`rndc nta`), BIND debug logging (`category dnssec { severity debug 3; }`), ordered triage methodology; source_count 1→2
  - `wiki/concepts/chain-of-trust.md` — added: 12-step validation walkthrough, trust anchor config table (`dnssec-validation auto` vs `yes`), three rollover methods comparison table, CSK variant section, clock skew row in "What Breaks the Chain" table; fixed duplicate content bug; source_count 2→3; confidence medium→high
  - `wiki/concepts/dnssec-internals.md` — added: CSK key type to key types table with "When to use" column, key lifecycle metadata table (Created/Publish/Activate/Inactive/Delete with rationale for Inactive→Delete gap), NSEC vs NSEC3 decision guidance table, RFC 9276 NSEC3 parameter guidance (0 iterations, empty salt); removed duplicate NSEC section; source_count 2→3; confidence medium→high
  - `wiki/index.md` — added `dnssec-guide-bind9` source entry; updated stats (20 total pages, 5 ingested)
  - `wiki/overview.md` — updated source inventory, source counts, and current focus
- **Entities touched**: none new
- **Concepts touched**: dnssec-in-bind, dnssec-debugging, chain-of-trust, dnssec-internals (all updated)
- **Contradictions found**: none
- **Notes**: This is the richest source ingested to date. Key contributions: (1) CSK type clarified — `dnssec-policy default` does NOT use the traditional KSK/ZSK split, it uses a Combined Signing Key; (2) key lifecycle state machine now fully documented (Created→Publish→Activate→Inactive→Delete) with the critical Inactive→Delete gap explained; (3) three KSK rollover methods (Double-KSK, Double-DS, Double-RRset) documented with trade-off table; (4) security lameness established as #1 production failure mode; (5) 12-step validation walkthrough adds concrete operational understanding of the chain of trust. Both `dnssec-in-bind` and `dnssec-debugging` are now `confidence: high`. All 6 concept pages are active (no stubs except `dnssec-vs-dns-resolution`).

---

## [2026-04-14] ingest | DNSSEC and BIND 9 (ISC)

- **Operation**: ingest
- **Source**: `raw/dnssec0guide/DNSSEC and BIND 9.md`
- **Source type**: Official ISC portal page — resource index + operational overview
- **Pages created**:
  - `wiki/sources/dnssec-and-bind9-isc.md` — source summary with key claims, resource index, evidence quality notes
  - `wiki/entities/Internet-Systems-Consortium.md` — entity page for ISC (BIND 9 maintainer)
  - `wiki/entities/DNSViz.md` — entity page for the diagnostic tool (ISC-endorsed)
- **Pages updated**:
  - `wiki/concepts/dnssec-in-bind.md` — added Operational Prerequisites section (EDNS0, NTP, secondaries); added KASP vs `auto-dnssec` comparison table; promoted from stub to active; source_count 0→1; confidence speculative→medium
  - `wiki/concepts/dnssec-internals.md` — expanded EDNS0/DO bit section with enterprise EDNS0-filtering risk note; source_count 1→2; removed EDNS0 open question (now answered)
  - `wiki/concepts/dnssec-debugging.md` — updated DNSViz entry with ISC endorsement quote; added EDNS0 stripping + clock skew failure modes to the failure table; updated Mentioned In; source_count 0→1; confidence speculative→medium
  - `wiki/index.md` — sources, entities, concepts, stats updated (19 total pages, 4 ingested, 5 active concepts)
- **Entities touched**: Internet-Systems-Consortium (created), DNSViz (created)
- **Concepts touched**: dnssec-in-bind (stub→active), dnssec-internals (updated), dnssec-debugging (stub→active)
- **Contradictions found**: none
- **Notes**: Source is a short portal/index page, not a deep technical document. Primary value: (1) ISC institutional confirmation that `dnssec-policy`/KASP is the current recommended path and that all pre-9.16 documentation should be treated as legacy, (2) two critical operational prerequisites now in the wiki — EDNS0 filtering and NTP clock sync — both are common real-world DNSSEC failure modes, (3) DNSViz now has its own entity page with ISC's direct endorsement. Phase 2 (BIND ops) is now underway; the next source (`DNSSEC Guide — BIND 9`) is the key document.

---


## [2026-04-13] ingest | DNS Vulnerabilities: The Kaminsky Exploit

- **Operation**: ingest
- **Source**: `raw/dnssec0guide/DNS Vulnerabilities Lessons from the Kaminsky Exploit and Its Lasting Impact on Internet Security.md`
- **Source type**: Retrospective article — dn.org
- **Pages created**:
  - `wiki/sources/kaminsky-exploit.md` — full source summary with step-by-step attack mechanics
  - `wiki/concepts/dns-cache-poisoning.md` — new concept page (active, confidence: high)
  - `wiki/entities/Dan-Kaminsky.md` — entity page
- **Pages updated**:
  - `wiki/concepts/chain-of-trust.md` — added "Why This Chain Had to Be Built" section linking the Kaminsky attack to DNSSEC's design rationale; source_count 1→2; added `dns-cache-poisoning` to related
  - `wiki/index.md` — sources, entities, concepts, stats updated
  - `wiki/overview.md` — source inventory updated
- **Entities touched**: Dan-Kaminsky (created)
- **Concepts touched**: dns-cache-poisoning (created, active), chain-of-trust (updated with motivation)
- **Contradictions found**: none
- **Notes**: This source completes Phase 1 of the learning roadmap (first principles). Key value: the precise Kaminsky subdomain-flooding technique is now documented step-by-step, explaining exactly why DNSSEC's cryptographic approach is the only structural solution. The attack-to-fix mapping table (port randomization vs DNSSEC per threat vector) is filed in the source page.

---

## [2026-04-13] ingest | How Does DNSSEC Work? (Cloudflare)

- **Operation**: ingest
- **Source**: `raw/dns-sec-cloudflare&text/How Does DNSSEC Work?.md`
- **Source type**: Tutorial article — Cloudflare Learning Center
- **Pages created**:
  - `wiki/sources/how-does-dnssec-work-cloudflare.md`
  - `wiki/entities/Cloudflare.md`
- **Pages updated**:
  - `wiki/concepts/dnssec-internals.md` — added: RRsets as signing unit, full record type table with parent/child ownership, intra-zone validation sequence, NSEC/NSEC3 detail, updated confidence to medium
  - `wiki/concepts/chain-of-trust.md` — added: full trust delegation chain diagram, DS record mechanics (concrete hash-and-compare procedure), KSK rollover double-DS procedure with TTL rationale, insecure delegation vs broken chain distinction, updated confidence to medium
  - `wiki/index.md` — sources, entities, stats updated
  - `wiki/overview.md` — source inventory, concept status updated
- **Entities touched**: Cloudflare (created)
- **Concepts touched**: dnssec-internals, chain-of-trust (both promoted from stub to active)
- **Contradictions found**: none
- **Notes**: Strong source for the mechanics layer. Particular value: the KSK rollover TTL-aware procedure (add-wait-remove) and the clear articulation of why ZSK rotation is cheap vs KSK rotation expensive. Source has a product pitch at the end (Cloudflare DNSSEC automation) — noted in evidence quality section.

---

## [2026-04-13] ingest | Learning Goals: DNSSEC (01-learn-dnssec.md)

- **Operation**: ingest
- **Source**: `raw/01-learn-dnssec.md`
- **Source type**: Personal note — learning intent / goals document
- **Pages created**:
  - `wiki/sources/01-learn-dnssec.md` — source summary + 6-axis coverage map
  - `wiki/synthesis/dnssec-learning-roadmap.md` — reading order, source inventory, gap analysis
  - `wiki/concepts/dnssec-internals.md` — stub
  - `wiki/concepts/chain-of-trust.md` — stub
  - `wiki/concepts/dnssec-vs-dns-resolution.md` — stub
  - `wiki/concepts/dnssec-in-bind.md` — stub
  - `wiki/concepts/dnssec-debugging.md` — stub
- **Pages updated**:
  - `wiki/overview.md` — domain established as DNSSEC; thesis, themes, source inventory added
  - `wiki/index.md` — all new pages cataloged; stats updated
- **Entities touched**: none
- **Concepts touched**: dnssec-internals, chain-of-trust, dnssec-vs-dns-resolution, dnssec-in-bind, dnssec-debugging (all stubs)
- **Contradictions found**: none
- **Notes**: Source is a goals document, not a content document. Treated as curriculum anchor. Created roadmap synthesis from the goals + prior gap analysis of all staged sources. 10 additional sources in raw/ await ingestion.

---

## [2026-04-10] init | Wiki bootstrapped

- **Operation**: Initial setup
- **Agent**: Sisyphus (OpenCode)
- **Pages created**: CLAUDE.md, wiki/index.md, wiki/log.md, wiki/overview.md
- **Templates created**: source, entity, concept, synthesis, query
- **Raw sources**: 0
- **Notes**: Vault bootstrapped from Karpathy's LLM Wiki pattern (https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f). Ready for source ingestion.

---

<!-- 
APPEND NEW ENTRIES ABOVE THIS LINE
Format for each entry:
## [YYYY-MM-DD] <operation> | <title>
- **Operation**: ingest | query | lint | evolve | update
- **Source**: <filename if ingest>
- **Pages created**: <list>
- **Pages updated**: <list>
- **Entities touched**: <list>
- **Concepts touched**: <list>
- **Contradictions found**: <list or none>
- **Notes**: <freeform>
-->

- **Operation**: Initial setup
- **Agent**: Sisyphus (OpenCode)
- **Pages created**: CLAUDE.md, wiki/index.md, wiki/log.md, wiki/overview.md
- **Templates created**: source, entity, concept, synthesis, query
- **Raw sources**: 0
- **Notes**: Vault bootstrapped from Karpathy's LLM Wiki pattern (https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f). Ready for source ingestion.

---

<!-- 
APPEND NEW ENTRIES ABOVE THIS LINE
Format for each entry:
## [YYYY-MM-DD] <operation> | <title>
- **Operation**: ingest | query | lint | evolve | update
- **Source**: <filename if ingest>
- **Pages created**: <list>
- **Pages updated**: <list>
- **Entities touched**: <list>
- **Concepts touched**: <list>
- **Contradictions found**: <list or none>
- **Notes**: <freeform>
-->
