# CineGraph

CineGraph is an open-source, local-first knowledge engine for AI agents that
reason about films, series, characters, events, and narrative relationships.
It combines a narrative graph with temporal knowledge, spoiler-aware queries,
and evidence tracking in an embedded SQLite database.

> [!IMPORTANT]
> CineGraph is in its foundation phase. The command currently exposes only
> version information; the MCP server and reasoning engine are planned for
> the first public release.

## Design goals

- **Local-first:** useful without a hosted dependency or mandatory account.
- **Agent-native:** a focused MCP interface over STDIO.
- **Spoiler-aware:** every query is evaluated against an explicit spoiler
  policy.
- **Evidence-backed:** claims retain provenance and confidence.
- **Time-aware:** narrative chronology and release chronology stay distinct.
- **Embeddable:** SQLite keeps deployment and data ownership simple.

## Planned shape

```text
AI agent
   │ MCP over STDIO
   ▼
Query and policy module
   │
   ├── narrative and temporal reasoning
   ├── spoiler policy enforcement
   └── evidence selection
   │
   ▼
SQLite storage adapter
```

The initial architecture is described in
[RFC 0001](docs/RFC-0001-architecture.md), and the first three months of work
are tracked in the [roadmap](docs/ROADMAP.md).

## Development

Prerequisites: Go 1.27 or newer.

```sh
go test ./...
go run ./cmd/cinegraph --version
```

See [CONTRIBUTING.md](CONTRIBUTING.md) before opening a change.

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE).

