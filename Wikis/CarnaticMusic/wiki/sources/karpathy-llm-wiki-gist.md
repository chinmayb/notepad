---
title: "LLM Wiki — A pattern for building personal knowledge bases using LLMs"
type: source
tags: ["meta", "llm", "knowledge-management", "wiki-pattern"]
created: 2026-05-08
updated: 2026-05-08
source_count: 1
confidence: high
status: active
related: ["[[llm-wiki-pattern]]", "[[Andrej-Karpathy]]", "[[overview]]"]
source_url: "https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f"
source_author: "Andrej Karpathy"
source_date: 2026-04-01
source_type: article
open_questions:
  - "At what wiki size does the index-only navigation break down and require a real search engine (e.g. qmd)?"
  - "How should images / audio assets be referenced in a Carnatic music context where audio is often the primary source?"
contradictions: []
entities_mentioned:
  - "[[Andrej-Karpathy]]"
concepts_mentioned:
  - "[[llm-wiki-pattern]]"
---

# LLM Wiki — A pattern for building personal knowledge bases using LLMs

> **Author**: Andrej Karpathy · **Date**: ~2026-04 · **Type**: Gist / idea document
> **Original**: [gist.github.com/karpathy/442a6bf555914893e9891c11519de94f](https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f)

---

## Summary

Karpathy proposes a pattern for personal knowledge management built around an **LLM-maintained wiki** that sits between raw source documents and the human reader. Unlike RAG systems (NotebookLM, ChatGPT file uploads), which re-derive answers from raw chunks on every query, this pattern compiles knowledge once into a persistent, interlinked markdown vault and keeps it current as new sources arrive. The wiki is a *compounding artifact* — cross-references, contradictions, and synthesis are already in place by the time you ask a question.

The architecture has three layers: (1) **Raw sources** — immutable input documents the LLM reads but never modifies; (2) **The wiki** — LLM-owned, LLM-maintained markdown files (entity pages, concept pages, source summaries, synthesis, an overview); (3) **The schema** — a configuration document (CLAUDE.md, AGENTS.md) that codifies conventions, frontmatter, and workflows. Three operations drive the system: **Ingest** (process a new source, touching 10–15 wiki pages), **Query** (answer questions with citations; file valuable answers back into the wiki), and **Lint** (periodic health check for orphans, stale claims, contradictions).

Two special files anchor navigation: **`index.md`** (content-oriented catalog organized by category) and **`log.md`** (chronological append-only operation record). Karpathy notes the index-only approach works well at moderate scale (~100 sources, hundreds of pages) and avoids embedding-based RAG infrastructure entirely; beyond that, a tool like [qmd](https://github.com/tobi/qmd) can layer hybrid BM25/vector search on top of the markdown.

The "why this works" argument is that **the bottleneck of personal wikis has always been maintenance, not reading or thinking** — humans abandon wikis because the cross-reference upkeep grows faster than the value. LLMs don't get bored, can touch 15 files in one pass, and reduce maintenance cost to near zero. The human curates and asks; the LLM does the bookkeeping.

---

## Key Claims

1. **RAG re-derives knowledge on every query; this pattern compiles it once and maintains it** — `confidence: high` — *"The LLM is rediscovering knowledge from scratch on every question. There's no accumulation."* The wiki accumulates structure (cross-refs, contradictions, synthesis) so future queries are cheaper and richer.
2. **The schema (CLAUDE.md / AGENTS.md) is the key configuration file** — `confidence: high` — *"It's what makes the LLM a disciplined wiki maintainer rather than a generic chatbot."* Without it the LLM doesn't know the conventions and outputs drift.
3. **Three operations cover the full lifecycle: ingest, query, lint** — `confidence: high` — *"A single source might touch 10-15 wiki pages."* The lint operation is what keeps the wiki healthy as it grows.
4. **An index file is sufficient navigation up to ~100 sources / hundreds of pages** — `confidence: medium` — *"This works surprisingly well at moderate scale ... and avoids the need for embedding-based RAG infrastructure."* Beyond that, layer in a real search tool like qmd.
5. **Good query answers should be filed back as new wiki pages** — `confidence: high` — *"This way your explorations compound in the knowledge base just like ingested sources do."* Otherwise valuable synthesis disappears into chat history.
6. **The bottleneck of personal wikis is maintenance, not reading** — `confidence: high` — *"Humans abandon wikis because the maintenance burden grows faster than the value. LLMs don't get bored."* This is the core "why now" argument.
7. **Obsidian + LLM agent is a productive pairing — Obsidian as IDE, LLM as programmer, wiki as codebase** — `confidence: medium` — workflow recommendation, not a falsifiable claim.
8. **The pattern is intentionally abstract / modular** — `confidence: high` — *"Pick what's useful, ignore what isn't."* No single prescribed directory layout; the agent should co-evolve specifics with the user.

---

## Notable Quotes

> "The wiki is a persistent, compounding artifact." — Karpathy, on the difference from RAG

> "Obsidian is the IDE; the LLM is the programmer; the wiki is the codebase." — Karpathy, on the practical workflow

> "Humans abandon wikis because the maintenance burden grows faster than the value. LLMs don't get bored, don't forget to update a cross-reference, and can touch 15 files in one pass." — Karpathy, the core "why now"

> "The idea is related in spirit to Vannevar Bush's Memex (1945) ... The part he couldn't solve was who does the maintenance. The LLM handles that." — Karpathy, on intellectual lineage

---

## Methodology / Evidence Quality

- **Sample size / scope**: A single idea document; no quantitative evaluation. Karpathy reports personal experience using the pattern but doesn't share metrics.
- **Methodology**: Argument-by-design plus invitation for the reader to instantiate and validate.
- **Potential biases**: The author is a high-context power user of LLMs (former OpenAI / Tesla AI lead). Patterns that work for him may need adaptation for less LLM-fluent users.
- **Replication / corroboration**: Multiple community implementations exist (the gist references "SwarmVault, OmegaWiki, Link" and others); the pattern is being used in the wild but no controlled comparison vs RAG exists.
- **Use-case coverage**: Karpathy lists six use cases (personal, research, books, business, competitive analysis, hobby) — broad, but Carnatic-music-specific affordances (audio-primary sources, oral-tradition variants, school/bani disagreements) are not addressed.

---

## Entities Mentioned

- [[Andrej-Karpathy]] — author of the gist; pattern originator

---

## Concepts Covered

- [[llm-wiki-pattern]] — the central pattern; this source is its definitional reference

---

## Contradictions with Existing Wiki

*None — this is the first source in the wiki.*

---

## My Notes

This source is the **bootstrap reference** for this vault. It is not Carnatic music content. It exists in the wiki because:

1. Future sessions need to know the schema came from somewhere — this is the citation.
2. It seeds the wiki with non-empty entity, concept, and source pages so the conventions can be exercised before real domain sources arrive.
3. The "why this works" argument applies *especially* well to Carnatic music, where most knowledge is dispersed across books, lec-dems, recordings, oral tradition, and bani-specific variants — a domain crying out for cross-referenced synthesis.

When the wiki has 5+ Carnatic-domain sources, this meta-source can be relegated visually (it'll still anchor [[llm-wiki-pattern]] but won't dominate the index).

---

*Page maintained by LLM agent. Source file: `raw/karpathy-llm-wiki-gist.md`*
