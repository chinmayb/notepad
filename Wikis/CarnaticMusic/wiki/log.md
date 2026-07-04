# Wiki Log

> Append-only chronological record of all wiki operations.
> Format: `## [YYYY-MM-DD] <operation> | <title>`
> Grep-friendly: `grep "^## \[" wiki/log.md | tail -10`

---

## [2026-05-08] ingest | Karpathy LLM Wiki Gist (meta-source)

- **Operation**: ingest
- **Source**: `raw/karpathy-llm-wiki-gist.md`
- **Source type**: Meta — pattern definition (gist by Andrej Karpathy)
- **Pages created**:
  - `wiki/sources/karpathy-llm-wiki-gist.md` — source summary with key claims, architecture (3 layers), operations (ingest/query/lint)
  - `wiki/entities/Andrej-Karpathy.md` — entity stub for the gist's author
  - `wiki/concepts/llm-wiki-pattern.md` — concept page describing the pattern itself
- **Pages updated**:
  - `wiki/index.md` — meta-source, meta-entity, and meta-concept added; stats initialized
  - `wiki/overview.md` — domain set to Carnatic Music; meta-source listed; thesis intentionally empty pending domain sources
- **Entities touched**: Andrej-Karpathy (created)
- **Concepts touched**: llm-wiki-pattern (created)
- **Contradictions found**: none
- **Notes**: This is a meta-ingest — the source describes the wiki pattern itself, not Carnatic music. Treated as a foundational entry so the wiki is non-empty and the schema/conventions are testable on a real source. Future ingests will be Carnatic-specific (book chapters, lec-dems, journal articles, listening notes). The owner should now define learning axes (see `overview.md` open questions) and stage the first domain sources.

---

## [2026-05-08] init | Wiki bootstrapped

- **Operation**: Initial setup
- **Pattern source**: [Karpathy's LLM Wiki gist](https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f)
- **Pages created**:
  - `CLAUDE.md` — schema and agent instructions (domain: Carnatic Music)
  - `raw/README.md` — raw source conventions and ingestion workflow
  - `wiki/index.md` — master catalog (initialized empty)
  - `wiki/log.md` — this file
  - `wiki/overview.md` — top-level synthesis (empty thesis pending sources)
- **Templates created**: `source-template.md`, `entity-template.md`, `concept-template.md`, `synthesis-template.md`, `query-template.md` — all five customized lightly for Carnatic music (raga/tala fields in concept template; bani/sampradaya in source/entity templates)
- **Directories created**: `raw/`, `raw/assets/`, `wiki/entities/`, `wiki/concepts/`, `wiki/sources/`, `wiki/synthesis/`, `wiki/queries/`, `wiki/templates/`
- **Raw sources**: 0 (the Karpathy gist was added in the same session — see next entry)
- **Notes**: Vault structure mirrors the existing DNSSec vault in this same workspace (sibling folder), which was previously bootstrapped from the same Karpathy pattern. Schema version 1.0. Domain: Carnatic Music — South Indian classical music.

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
