# Contributing to CineGraph

Thank you for helping build CineGraph.

## Before starting

- Read the [glossary](GLOSSARY.md) and use its domain language.
- Check existing issues and discuss large interface or schema changes first.
- Keep changes focused on one observable outcome.
- Do not include copyrighted film scripts, subtitles, or non-redistributable
  datasets in fixtures.

## Development workflow

1. Fork the repository and create a focused branch.
2. Add or update tests for the behavior being changed.
3. Run the local checks:

   ```sh
   gofmt -w .
   go vet ./...
   go test ./...
   ```

4. Explain the motivation, behavior, and verification in the pull request.

## Design changes

Changes to the Engine interface, persistent schema, spoiler semantics, or MCP
tool contract should begin with a concrete scenario. Update an existing RFC or
propose a focused RFC when the decision is difficult to reverse and involves a
real trade-off.

## Reporting security issues

Please do not open a public issue for a suspected vulnerability. Use GitHub's
private vulnerability reporting feature when it becomes available for this
repository.

