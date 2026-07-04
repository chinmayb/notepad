# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Vault Is

This is an Obsidian vault: a personal, multi-topic learning wiki — a "second brain" that accumulates knowledge across whatever the owner is currently studying. It is not a software project; there are no build, lint, or test commands. Each top-level directory is an independent topic space (`DNSSec/`, `CarnaticMusic/`, `Protoactor/`, `LLMS/`, ...), and new ones get added whenever a new subject grabs interest. Treat each topic directory as self-contained — don't assume one topic's conventions apply to another unless this file says so.

## The LLM Wiki Pattern

Some topic directories are built on the "LLM Wiki" pattern from Andrej Karpathy: https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f. The core idea: raw sources are curated by the human, and an LLM agent does all the summarizing, cross-referencing, and filing, so the knowledge base compounds across sessions instead of scattering across chat history.

**This section is the canonical schema.** A topic directory using this pattern has its own `<Topic>/CLAUDE.md`, which is a domain-specific *overlay* on top of what's below — it must not restate these rules, only note where the topic deviates (domain framing, what counts as an entity vs. a concept, naming quirks). If a topic's `CLAUDE.md` conflicts with this one, the topic file wins *for that topic only*.

### Directory layout of a wiki-pattern topic

```
<Topic>/
├── CLAUDE.md              ← domain overlay (see "Adding a New Topic" below)
├── raw/                   ← IMMUTABLE source documents; agent reads, never modifies
│   ├── assets/            ← images/PDFs referenced by raw sources
│   ├── README.md          ← conventions for adding sources
│   └── <sources>          ← articles, papers, transcripts, notes (markdown preferred)
└── wiki/                  ← entirely LLM-owned and maintained
    ├── index.md           ← master catalog of all wiki pages, updated on every operation
    ├── log.md             ← append-only chronological record of all operations
    ├── overview.md        ← evolving thesis / top-level synthesis
    ├── entities/          ← pages for people, organizations, products, projects
    ├── concepts/          ← pages for ideas, frameworks, techniques, terms
    ├── sources/           ← one summary page per ingested raw source
    ├── synthesis/         ← cross-source analysis, comparisons, structured arguments
    │   └── <topic>-learning-roadmap.md   ← REQUIRED — first-principles reading order (see below)
    ├── queries/           ← filed answers to questions worth persisting
    └── templates/         ← page templates (starting points — don't restructure them)
```

### Page frontmatter (every wiki page)

```yaml
---
title: "Page Title"
type: entity | concept | source | synthesis | query | overview
tags: []                   # topic tags
created: YYYY-MM-DD
updated: YYYY-MM-DD
source_count: 0             # number of raw sources that contributed to this page
confidence: high | medium | speculative
status: active | superseded | archived
related: []                 # wikilinks to strongly related pages, e.g. ["[[PageName]]"]
open_questions: []          # unresolved questions surfaced by this page
contradictions: []          # note where sources disagree, e.g. "Source A says X, Source B says Y"
---
```

- `confidence: speculative` — claim derived from weak/single evidence
- `confidence: medium` — corroborated by 2+ sources but not definitive
- `confidence: high` — well-evidenced across multiple quality sources
- `status: superseded` — must include a `superseded_by` field pointing to the replacement page
- `status: archived` — kept for historical record, no longer actively maintained
- Always use `[[WikiLink]]` syntax for internal references, never bare text

### Operations

**1. INGEST — Processing a new source**

1. Read the raw source fully; summarize the 5 key takeaways and discuss before writing
2. Write `wiki/sources/<slug>.md` — summary, falsifiable key claims, quotes worth preserving, entities/concepts mentioned
3. Update or create the affected `wiki/entities/` pages
4. Update or create the affected `wiki/concepts/` pages, cross-referencing existing related concepts
5. Flag contradictions: add a `> ⚠️ Contradiction:` callout inline and a `contradictions:` entry in frontmatter on the affected page(s)
6. Update `wiki/index.md` under the appropriate section
7. Append an entry to `wiki/log.md`: `## [YYYY-MM-DD] ingest | <Source Title>` with source type, pages touched, entities/concepts touched, contradictions found
8. Update `wiki/overview.md` if the source materially changes the synthesis, or increment supporting evidence if it just confirms an existing claim

**2. QUERY — Answering a question**

1. Read `wiki/index.md` to find relevant pages, then read those pages fully
2. Synthesize an answer with `[[WikiLink]]` citations
3. If the answer required non-trivial synthesis (3+ pages connected, a revealed pattern, a comparison), ask whether to file it as `wiki/queries/<slug>.md`
4. Output formats: default narrative markdown, `table:` prefix for a comparison table, `marp:` prefix for a slide deck, `chart:` prefix for a chart description

**3. LINT — Wiki health check**

Check for: orphan pages (no inbound links), broken `[[links]]`, stale `confidence: high` claims sourced only from old material, entities/concepts mentioned inline but missing their own page, asymmetric `related:` links, unsourced claims in `synthesis/` pages, and `contradictions:` open more than 30 days. File the report to `wiki/queries/lint-YYYY-MM-DD.md`.

**4. EVOLVE — Schema evolution**

When the schema doesn't fit a new document type or workflow: propose the change → get approval → update the relevant `CLAUDE.md` (this file for shared schema changes, the topic's file for domain-specific changes) → retroactively update affected pages.

### Naming conventions

| Type | Location | Filename convention |
|---|---|---|
| Source summary | `wiki/sources/` | `kebab-case-title.md` |
| Entity | `wiki/entities/` | `FirstLast.md` or `OrgName.md` |
| Concept | `wiki/concepts/` | `kebab-case-concept.md` |
| Synthesis | `wiki/synthesis/` | `topic-comparison.md` or `topic-analysis.md` |
| Query answer | `wiki/queries/` | `YYYY-MM-DD-short-question.md` |
| Lint report | `wiki/queries/` | `lint-YYYY-MM-DD.md` |

### The Learning Roadmap (required)

Every wiki-pattern topic must maintain exactly one `wiki/synthesis/<topic>-learning-roadmap.md`. It is the designated entry point for a human who wants to know "where do I start reading this topic" — the fix for a wiki that has grown enough pages that it's no longer obvious which file to open first.

- **Job**: lay out the topic's concepts in the order they build on each other, starting from first principles (the motivating problem or the foundational idea) through to advanced/operational material. The specific phases and how many there are is domain-specific and chosen per topic — this rule fixes the *role* of the page, not its content shape.
- **Link-accurate**: it only links to concept pages that actually exist. If a planned page ends up folded into a different page instead of being built standalone, correct the roadmap in the same edit — no aspirational or dead links.
- **Discoverable**: linked prominently near the top of `wiki/index.md` (not buried in the Synthesis table) and referenced in `wiki/overview.md`'s "Current Focus" section.
- **Maintained on ingest**: INGEST step 6 (updating `index.md`) also means updating the topic's learning roadmap whenever a source adds, merges, or reshapes a concept page.
- **Bootstrapped from** `wiki/templates/learning-roadmap-template.md`.
- A topic with zero or one concept pages doesn't need a populated roadmap yet — an empty/skeleton roadmap is fine; don't fabricate a reading order richer than the ingested content supports (see "Flag insufficient sources" below).

### Cross-reference and contradiction rules

- Every entity or concept first mentioned in any wiki page must be wikilinked (`[[PageName]]`); if the page doesn't exist yet, create a stub and note it in the log
- Source pages list every entity/concept they contribute to; entity/concept pages list every source that mentions them
- The `related:` field is for *strong* relationships only (same domain, direct interaction, definitional dependency)
- Contradictions are first-class citizens — never silently overwrite a conflicting claim. When source B contradicts source A: keep both claims visible, add the `> ⚠️ Contradiction:` callout, record it in `contradictions:` frontmatter on both the affected page and both source pages, lower `confidence` on the contested claim if warranted, and flag it for the human to resolve if it needs a judgment call

### Rules that apply everywhere

- Never modify files in `raw/` — they are immutable source of truth
- Never delete a wiki page — archive it (`status: archived`) instead
- Never silently resolve a contradiction — surface it
- Never write a wiki page without full frontmatter
- Never leave `index.md` or `log.md` stale after an operation
- Never create a duplicate page — check the index first
- Never commit changes to git unless explicitly asked
- **Flag insufficient sources.** If a raw source is too thin to responsibly support a full page (a fragment, a one-line note, an ambiguous or unverifiable claim), say so instead of padding it out. File a stub with `confidence: speculative` and an `open_questions:` entry, or ask for more material before writing a page that overstates what the source actually supports.

### Session start protocol (per topic)

1. Read this file and the topic's own `CLAUDE.md`
2. Read that topic's `wiki/index.md` to orient to what exists
3. Read the last 5 entries of `wiki/log.md` — `grep "^## \[" wiki/log.md | tail -10` is the fast way to do this
4. Ask what to work on, or proceed if instructions were already given

## Vault Index

[`index.md`](index.md) at the vault root is the human navigation entry point across topics — one row per topic, linking to that topic's learning roadmap (or equivalent entry point for lighter patterns). It answers "which file do I open to start learning about X," distinct from a topic's own `wiki/index.md`, which is a reference catalog for looking things up *within* that topic once you're already there.

## Topic Directories Today

| Directory | Pattern | Status |
|---|---|---|
| `DNSSec/` | Full LLM-wiki pattern (`raw/` + `wiki/`) | Mature — 23 wiki pages, 8 of 11 staged sources ingested |
| `CarnaticMusic/` | Full LLM-wiki pattern | Just bootstrapped — only the meta-source (the Karpathy gist itself) ingested so far |
| `Protoactor/` | Lightweight: numbered `modules/NN-Topic.md` files + a `Readme.md` index with a learning-path table | Standalone reference notes, not agent-maintained, no `CLAUDE.md`. Its learning-path table already plays the same role as a learning roadmap for this lighter pattern — not a gap. |
| `LLMS/` | Bare notes file | Not yet structured |
| `Clippings/` | Empty | Obsidian Web Clipper's default save target |

## Adding a New Topic

Choose the pattern deliberately — don't force every topic into the full wiki:

- **Use the full LLM-wiki pattern** (`raw/` + `wiki/`) when the topic will accumulate many sources over time and benefits from cross-referenced entities/concepts worth querying later. Bootstrap it by copying the `wiki/templates/` and `raw/README.md` conventions from `DNSSec/` or `CarnaticMusic/`, create its `<topic>-learning-roadmap.md` from `wiki/templates/learning-roadmap-template.md` (a skeleton is fine before most sources are ingested), then write a short domain-overlay `CLAUDE.md` (template below).
- **Use a lighter structure** — like `Protoactor/modules/`, a flat sequence of numbered files with a `Readme.md` index — for linear, once-through, course-style learning where you're moving through a fixed curriculum rather than building a long-lived reference. No `CLAUDE.md` is required for this style.
- A topic can start light and graduate to the full pattern later if it accumulates enough sources to be worth the overhead.
- Add a row for the new topic to the root [`index.md`](index.md) pointing to its entry point (its learning roadmap once one exists, or its equivalent for a lighter pattern).

### Domain-overlay CLAUDE.md template

```markdown
# CLAUDE.md — <Topic> overlay

Schema: see [../CLAUDE.md](../CLAUDE.md) for the full LLM Wiki pattern (frontmatter, INGEST/QUERY/LINT/EVOLVE, naming, cross-reference rules). This file only covers what's specific to <Topic>.

## Domain

<one paragraph: what this topic covers>

## Domain-specific framing

<what counts as an entity vs. concept here; naming quirks; anything that changes how contradictions are judged in this domain>
```

## Known Gaps / Flags

- The Karpathy gist is the only source describing this pattern. It was reachable via fetch and sufficient to write this file — it's also already ingested as a wiki source at `CarnaticMusic/wiki/sources/karpathy-llm-wiki-gist.md`. The full pattern is captured above, so nothing here depends on the live gist staying up if it ever moves or disappears.
