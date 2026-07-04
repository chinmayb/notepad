---
title: "LLM Wiki Pattern"
type: concept
tags: ["meta", "knowledge-management", "wiki", "llm"]
created: 2026-05-08
updated: 2026-05-08
source_count: 1
confidence: high
status: active
related: ["[[karpathy-llm-wiki-gist]]", "[[Andrej-Karpathy]]", "[[overview]]"]
open_questions:
  - "What's the empirical break-even between index-only navigation and search-tool-assisted navigation in a Carnatic-music vault (audio-heavy, oral-tradition-heavy)?"
  - "How should the pattern adapt when the primary source is audio rather than text?"
contradictions: []
---

# LLM Wiki Pattern

> **Domain**: Personal knowledge management / LLM tooling (meta — describes how *this vault* works, not Carnatic music)
> **Also known as**: LLM-maintained wiki; agentic knowledge base; persistent-wiki RAG alternative

---

## Definition

**Core definition**: A personal knowledge base architecture in which an LLM agent — not the human — writes and maintains a structured, interlinked collection of markdown files (the "wiki") that sits between raw source documents and the reader. The human curates sources and asks questions; the LLM does the summarizing, cross-referencing, contradiction tracking, and bookkeeping. The wiki is a *persistent, compounding artifact* — knowledge accumulates rather than being re-derived on every query.

**Alternative framings**:
- **vs. RAG**: RAG retrieves raw chunks at query time and re-synthesizes; this pattern compiles knowledge once and incrementally maintains it.
- **vs. traditional personal wiki**: traditional wikis stall because humans don't keep cross-references current; here the LLM handles maintenance.
- **vs. Memex (Vannevar Bush, 1945)**: Bush envisioned associative trails between documents but couldn't solve "who maintains them"; the LLM is the missing piece.

---

## Why It Matters

For this Carnatic music vault: **the domain is unusually well-suited to the pattern**. Carnatic knowledge is dispersed across books, journal articles, lec-dems, recordings, liner notes, and oral tradition; the same raga or kriti is described differently across schools (bani / sampradaya) and eras; and most students accumulate years of fragmented notes. A LLM-maintained wiki can do the cross-bani comparison, the lakshya-vs-lakshana reconciliation, and the composer-disciple-lineage mapping that humans rarely have time to do. The wiki becomes a long-term study companion.

---

## Key Properties

1. **Three layers**: raw sources (immutable) → wiki (LLM-owned, LLM-maintained) → schema (the configuration file, e.g. `CLAUDE.md`).
2. **Three operations**: **ingest** (process a new source, touch 10–15 wiki pages), **query** (answer with citations; file valuable answers back), **lint** (periodic health check).
3. **Two anchor files**: `index.md` (content catalog, by category) and `log.md` (chronological append-only).
4. **Frontmatter discipline**: every page has YAML metadata (title, type, tags, dates, source_count, confidence, status, related, open_questions, contradictions).
5. **Cross-references are first-class**: every entity/concept first mentioned anywhere becomes a `[[WikiLink]]`; orphan pages and asymmetric links are lint targets.
6. **Contradictions are first-class**: never silently overwritten; surfaced via `> ⚠️ Contradiction:` callouts and tracked in frontmatter.
7. **Confidence is graded**: speculative / medium / high based on source count and quality.
8. **Modular**: the abstract pattern is mandatory; specific directory names, page formats, and tooling are co-evolved with the user.

---

## How It Works

### Ingest workflow (8 steps in this vault)
1. Read the raw source fully
2. Discuss key takeaways with the user before writing
3. Write a source summary page
4. Update entity pages for every person/org/work mentioned
5. Update concept pages for every raga/tala/idea mentioned
6. Flag contradictions with existing wiki content
7. Update `index.md`
8. Append a `log.md` entry

### Query workflow
1. Read `index.md` to find candidate pages
2. Read those pages fully
3. Synthesize an answer with `[[WikiLink]]` citations
4. Offer to file the answer as a query page if synthesis was non-trivial

### Lint workflow
Check for: orphans, broken links, stale `confidence: high` claims, mentioned-but-missing pages, asymmetric `related:` links, unsourced synthesis claims, unresolved contradictions older than 30 days.

---

## Examples

- **DNSSec sibling vault** (`../DNSSec/`) — same pattern instantiated for the DNSSEC protocol. 23 wiki pages, 8 ingested sources, 6 active concepts, 1 contradiction tracked. Demonstrates the pattern's effectiveness at moderate scale.
- **This vault** (`./`) — pattern instantiated for Carnatic music, currently bootstrapped with only the meta-source.

---

## Relationship to Other Concepts

- **Contrasts with** RAG (retrieval-augmented generation): RAG retrieves raw chunks per query; this pattern compiles a maintained intermediate layer.
- **Contrasts with** traditional personal wikis: human-maintained wikis stall on cross-reference upkeep; this pattern offloads that to the LLM.
- **Enabled by**: long-context LLMs that can touch 15 files in one pass without losing coherence.
- **Tooling-adjacent**: Obsidian (UI / graph view), git (versioning), Marp (slide output), Dataview (frontmatter queries), qmd (search at scale).

---

## Evidence & Sources

| Source | Stance | Key Claim |
|--------|--------|-----------|
| [[karpathy-llm-wiki-gist]] | defines | The pattern itself; "the wiki is a persistent, compounding artifact" |

---

## Open Questions

- At what wiki size does index-only navigation break and require a real search tool?
- How should the pattern adapt for **audio-primary sources** (recorded concerts, lec-dems) which are central to Carnatic music but cannot be ingested as text directly?
- Should "contradictions across bani / sampradaya" get a dedicated frontmatter field, or is the existing `contradictions:` field sufficient?

---

## Contradictions

*(none yet — this concept is defined by a single source)*

---

*Page maintained by LLM agent.*
