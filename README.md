# fscli

A cli tool for firestore

## Demo

![fscli-demo](https://github.com/maruware/fscli/assets/1129887/887bbc4c-4f66-40a5-9211-256899abc067)

## Installation

```shell
brew install maruware/tap/fscli
```

or download binary from [Releases](https://github.com/maruware/fscli/releases)

or

```
go install github.com/maruware/fscli/cmd/fscli@latest
```

## Prepare

```sh
$ gcloud auth application-default login
```

## Usage

```sh
$ fscli --project-id my-project
```

| Flag | Description |
|------|-------------|
| `--project-id` | Firebase project ID (required) |
| `--out-mode` | Output format: `table` (default) or `json` |

### Quick Examples

```sql
QUERY users
QUERY users WHERE age = 20
QUERY users SELECT name WHERE age >= 20 ORDER BY name ASC LIMIT 10
GET users/ewpSGf5URC1L1vPENbxh
COUNT users WHERE name = "takashi"
```

## Documentation

- [Operations](docs/operations.md) — `QUERY`, `GET`, `COUNT`, collection paths
- [WHERE Filters](docs/where-filters.md) — Operators (`=`, `!=`, `>`, `<`, `IN`, `ARRAY_CONTAINS`, ...), value types, `TIMESTAMP()`, `__id__`
- [Clauses](docs/clauses.md) — `SELECT`, `ORDER BY`, `LIMIT`
- [Meta Commands](docs/meta-commands.md) — `\d`, `\pager`
- [Output](docs/output.md) — Table / JSON output modes, non-interactive mode

### JSON mode

The `--out-mode json` flag produces a single, valid JSON object (for `GET`) or an array of objects (for `QUERY`), making it easy to parse with tools like `jq`.

**GET command output:**
```json
{
  "id": "documentId",
  "data": { "field": "value" }
}
```

**QUERY command output:**
```json
[
  { "id": "doc1", "data": { "field": "value" } },
  { "id": "doc2", "data": { "field": "value" } }
]
```
