# RFC 0001: Initial Architecture

- Status: Draft
- Created: 2026-10-14
- Target: v0.1.0

## Summary

CineGraph will be a local-first Go engine that stores narrative and temporal
knowledge in SQLite and exposes spoiler-aware, evidence-backed queries to AI
agents through MCP over STDIO.

This RFC fixes the initial module seams and invariants while deliberately
leaving the graph schema and MCP tool vocabulary open for focused follow-up
RFCs.

## Motivation

General-purpose knowledge stores do not naturally distinguish story chronology
from release chronology. They also tend to return plausible statements without
showing their evidence or checking whether the statement reveals information
beyond a user's progress. CineGraph treats time, evidence, continuity, and
spoiler safety as core parts of every query.

## Principles

1. Local operation must not require a hosted account or network connection.
2. A result without its evidence and disclosure decision is incomplete.
3. Story Time and Release Time must never be represented by one ambiguous
   timestamp.
4. Spoiler enforcement belongs inside the query path, not in each caller.
5. Domain behavior must be testable without starting an MCP process or opening
   a production database.
6. Public interfaces remain small; implementation detail stays local to the
   owning module.

## Proposed modules

### Engine

The Engine is the primary interface for callers. It accepts structured commands
and queries and returns results containing facts, evidence, and disclosure
metadata. It owns orchestration so callers do not need to coordinate temporal,
spoiler, and evidence rules themselves.

Initial conceptual operations:

- ingest evidence-backed claims;
- resolve entities and works;
- query narrative relationships and events;
- explain why a result was included, withheld, or marked uncertain.

The exact Go interface will be proposed only after representative scenarios
exist. This avoids freezing storage-shaped methods into the public interface.

### Narrative model

The narrative model owns Works, Entities, Narrative Facts, Narrative Events,
Continuities, and their invariants. It does not know about SQLite or MCP.

### Temporal reasoning

Temporal reasoning evaluates intervals and ordering in both Story Time and
Release Time. Missing or disputed dates remain explicit rather than being
silently replaced by an arbitrary order.

### Spoiler policy

The spoiler policy decides whether a candidate result can be disclosed for a
given Work, Continuity, and Spoiler Horizon. It returns a decision with a
reason, allowing the same behavior to serve MCP, tests, and future adapters.

### Evidence tracking

Evidence tracking ranks and groups Claims without collapsing disagreement.
Conflicting Claims remain inspectable and retain their sources.

### Storage

The storage seam persists and retrieves domain records and supports atomic
ingestion. SQLite is the first adapter. The initial release will not add a
second database adapter merely to create a hypothetical abstraction.

### MCP adapter

The MCP adapter translates STDIO requests into Engine operations and maps
results to stable, agent-readable responses. Protocol details do not enter the
domain model.

## Data invariants

- Every Claim refers to at least one Evidence record.
- Every disclosure decision is evaluated against an explicit Spoiler Horizon;
  an omitted horizon defaults to the safest behavior, not unrestricted access.
- Story Time and Release Time use distinct fields and types.
- Claims may disagree; ingestion never overwrites disagreement merely because
  a newer Claim arrives.
- Continuity is explicit whenever a Work can belong to multiple continuities.
- Persistent identifiers remain stable across aliases and corrected metadata.

## Initial package direction

```text
cmd/cinegraph/       process entry point
internal/engine/     orchestration and primary interface
internal/narrative/  domain types and invariants
internal/temporal/   interval and ordering rules
internal/spoiler/    disclosure decisions
internal/evidence/   claim and source evaluation
internal/sqlite/     SQLite adapter and migrations
internal/mcp/        STDIO protocol adapter
```

Packages will be created as their first behavior is implemented; empty package
scaffolding is intentionally avoided.

## Error model

Errors must distinguish invalid input, unresolved references, conflicts,
storage failure, and protocol failure. Expected uncertainty in story knowledge
is represented in results rather than surfaced as an operational error.

## Testing strategy

- Domain invariants use table-driven unit tests.
- Engine scenarios exercise the public interface with deterministic adapters.
- SQLite tests run against temporary databases and real migrations.
- MCP tests communicate with the process over STDIO.
- Golden fixtures cover spoiler-sensitive responses and conflicting evidence.

## Security and privacy

- No telemetry or network access is enabled by default.
- The database path is explicit and local.
- Imported sources are treated as untrusted input.
- MCP responses do not reveal hidden facts while explaining a spoiler decision.
- Resource limits will protect parsing and query execution from oversized input.

## Deferred decisions

- Exact SQLite schema and migration tool
- Claim confidence model
- MCP tool names and response schemas
- Import formats and upstream datasets
- Full-text or vector search
- Multi-process write coordination

Each deferred decision needs scenario evidence before becoming a committed
interface.

## Acceptance criteria

This RFC is ready to accept when follow-up issues define representative query,
spoiler, evidence-conflict, and temporal scenarios; the proposed seams can
support those scenarios without exposing storage or MCP details.

