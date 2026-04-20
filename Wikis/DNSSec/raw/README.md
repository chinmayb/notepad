# Raw Sources

> This folder contains immutable source documents.
> **The LLM agent reads from here but NEVER modifies these files.**

---

## What Goes Here

Drop any of the following into this folder:
- **Articles / blog posts** — clipped to markdown (use Obsidian Web Clipper)
- **Research papers** — PDF or markdown export
- **Podcast/video transcripts** — plain text or markdown
- **Book chapters / excerpts** — markdown
- **Your own notes / journal entries** — markdown
- **Meeting transcripts / Slack exports** — markdown
- **Data files** — CSV, JSON (reference from wiki, don't embed raw data in wiki pages)

Images go in `raw/assets/` — after clipping an article in Obsidian, use the "Download attachments" hotkey (Ctrl+Shift+D) to pull all images to local disk.

---

## Naming Conventions

| Source Type | Filename Convention | Example |
|------------|---------------------|---------|
| Article | `YYYY-MM-DD-kebab-title.md` | `2026-03-15-attention-is-all-you-need.md` |
| Paper | `YYYY-author-short-title.md` | `2017-vaswani-transformer.md` |
| Transcript | `YYYY-MM-DD-source-topic.md` | `2026-01-10-lex-fridman-hinton.md` |
| Personal note | `YYYY-MM-DD-note-topic.md` | `2026-04-10-note-observations-rlhf.md` |
| Book chapter | `book-title-chNN.md` | `poor-charlies-almanack-ch03.md` |

---

## Source Metadata Header (Recommended)

Paste this at the top of every source file before ingesting:

```markdown
---
source_title: "Full Title Here"
source_type: article | paper | transcript | book_chapter | personal_note | data
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

*(none yet)*

---

*Files in this folder are source of truth. Never delete or modify them after ingestion.*
