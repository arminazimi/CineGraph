# Three-Month Roadmap

This roadmap covers 2026-10-14 through 2027-01-13. The dates express delivery
intent, not a promise to trade correctness for schedule.

## Month 1 — Foundation

**Window:** 2026-10-14 to 2026-11-13  
**Outcome:** a tested domain foundation that can persist evidence-backed
narrative knowledge locally.

- Review and accept the initial architecture RFC.
- Define representative narrative, temporal, spoiler, and evidence scenarios.
- Establish the Go module, CI, contribution workflow, and release conventions.
- Define core domain types and invariants.
- Design the first SQLite schema and migration mechanism.
- Implement atomic ingestion and retrieval through the Engine.
- Import a small, legally redistributable demonstration dataset.

**Exit criteria:** a local program can create a database, ingest the fixture,
and retrieve claims with their evidence in repeatable tests.

## Month 2 — Reasoning Core

**Window:** 2026-11-14 to 2026-12-13  
**Outcome:** queries correctly account for narrative time, release time,
continuity, conflicting evidence, and spoiler horizons.

- Implement interval and partial-order reasoning.
- Implement continuity-aware relationship and event queries.
- Define and implement the spoiler policy.
- Rank and explain evidence without hiding conflicts.
- Add scenario fixtures for films, episodic series, nonlinear stories, and
  alternate continuities.
- Measure correctness and query latency on the demonstration corpus.

**Exit criteria:** scenario tests show that the same query can safely return
different, explainable results for different spoiler horizons and continuities.

## Month 3 — Agent Interface and v0.1.0

**Window:** 2026-12-14 to 2027-01-13  
**Outcome:** AI agents can install and use CineGraph locally through a stable
MCP-over-STDIO interface.

- Define MCP tools from validated Engine scenarios.
- Implement lifecycle, capability discovery, and structured errors.
- Expose work/entity lookup, relationship, timeline, and evidence queries.
- Add end-to-end STDIO tests and an example agent configuration.
- Document installation, data import, spoiler semantics, and troubleshooting.
- Run release checks and publish v0.1.0 artifacts and checksums.

**Exit criteria:** a clean machine can install CineGraph, import the fixture,
connect an MCP client, run documented queries, and receive evidence-backed,
spoiler-safe results.

## Beyond v0.1.0

Possible later work includes additional import adapters, full-text search,
optional embeddings, richer confidence models, and multi-process coordination.
These are not commitments in this roadmap.
