# Raw Sources

> This folder contains immutable source documents.
> **The LLM agent reads from here but NEVER modifies these files.**

---

## What Goes Here

Drop any of the following into this folder:
- **Articles / blog posts** — clipped to markdown (use Obsidian Web Clipper)
- **Research papers / journal articles** — PDF or markdown export (e.g. JMA, Sangeet Natak)
- **Concert / lecture-demonstration transcripts** — plain text or markdown
- **Book chapters / excerpts** — markdown (e.g. Sambamoorthy's *South Indian Music*, Subbarama Dikshitar's *Sangita Sampradaya Pradarsini*)
- **Your own notes / journal entries / class notes** — markdown
- **Liner notes, sleeve essays, festival booklets** — markdown
- **Audio / video links + your transcribed observations** — markdown referencing files in `assets/`
- **Notation files** — reference from wiki, don't embed raw notation in wiki pages

Images and audio go in `raw/assets/` — after clipping an article in Obsidian, use the "Download attachments" hotkey (Ctrl+Shift+D) to pull all images to local disk.

---

## Naming Conventions

| Source Type | Filename Convention | Example |
|------------|---------------------|---------|
| Article | `YYYY-MM-DD-kebab-title.md` | `2024-09-12-the-lineage-of-tyagaraja.md` |
| Paper | `YYYY-author-short-title.md` | `1998-viswanathan-raga-grammar.md` |
| Transcript | `YYYY-MM-DD-source-topic.md` | `2025-03-04-tm-krishna-lecdem-todi.md` |
| Personal note | `YYYY-MM-DD-note-topic.md` | `2026-05-08-note-kalyani-listening.md` |
| Book chapter | `book-title-chNN.md` | `sambamoorthy-south-indian-music-ch04.md` |
| Liner notes | `YYYY-artist-album.md` | `1976-msss-meera-bhajans.md` |

---

## Source Metadata Header (Recommended)

Paste this at the top of every source file before ingesting:

```markdown
---
source_title: "Full Title Here"
source_type: article | paper | transcript | book_chapter | personal_note | liner_notes | data
source_url: "https://..."       # original URL if applicable
source_date: YYYY-MM-DD         # publication/creation date
source_author: "Name(s)"
added_date: YYYY-MM-DD          # date you added it to this vault
ingested: false                  # set to true after LLM ingestion
---
```

The LLM agent will set `ingested: true` after processing.

---

## Ingestion Workflow

1. Drop the file here (with metadata header above)
2. Tell the agent: `ingest raw/<filename.md>`
3. The agent reads, discusses key takeaways, then writes to `wiki/`
4. Check the updated pages in Obsidian's graph view

**One source at a time is recommended** — staying involved in the ingest lets you guide emphasis and catch misinterpretations immediately.

---

## Current Sources

- `karpathy-llm-wiki-gist.md` — Andrej Karpathy's gist describing the LLM Wiki pattern (the meta-source that bootstrapped this vault)

---

*Files in this folder are source of truth. Never delete or modify them after ingestion.*
