# CLAUDE.md — CarnaticMusic overlay

Schema: see [../CLAUDE.md](../CLAUDE.md) for the full LLM Wiki pattern (frontmatter, INGEST/QUERY/LINT/EVOLVE, naming, cross-reference rules). This file only covers what's specific to Carnatic Music.

## Domain

Carnatic Music — the classical music tradition of South India. The wiki accumulates knowledge about ragas, talas, compositions, composers, performers, instruments, theoretical frameworks, and performance practice.

## Domain-specific framing

- **Ragas and talas** live in `wiki/concepts/`. **Composers, performers, gurus, sabhas, and instruments** live in `wiki/entities/`. A specific kriti/composition can go in either — file it under whichever gives better cross-linking (typically `entities/` as a work, or `concepts/` as an exemplar of a raga/form).
- **Entity naming**: composers and performers use `FirstLast.md` (e.g. `Tyagaraja.md`, `MS-Subbulakshmi.md`); organizations and instruments use `OrgName.md` / `InstrumentName.md`.
- **`related:` convention**: a janya raga lists its parent melakarta; a kriti lists its raga, tala, and composer.
- **Contradictions matter more here than elsewhere** — the same raga, tala, or composition can be legitimately described differently across bani/sampradaya (schools), eras, or regional traditions. Don't treat divergent descriptions as an error to resolve; surface them as a contradiction and let both stand unless there's a clear reason to prefer one.
- The concept template (`wiki/templates/concept-template.md`) is customized with property tables: ragas get arohana/avarohana/vadi-samvadi/gamaka/rasa; talas get aksharas/angas/jati/gati.
- Prefer one canonical concept page per raga (e.g. `[[kalyani]]`) and let janya/melakarta relationships live in `related:` frontmatter rather than duplicating content across pages.

*Schema version: 1.0 | Created: 2026-05-08 | Vault: CarnaticMusic | Pattern: [Karpathy LLM Wiki](https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f)*
