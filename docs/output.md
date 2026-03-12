# Output

## Output Modes

fscli supports two output formats, configurable via the `--out-mode` flag.

### Table (default)

ASCII table format with headers.

```sh
$ fscli --project-id my-project
> QUERY users
+----------------------+---------+-----+
|          ID          |  name   | age |
+----------------------+---------+-----+
| VfsA2DjQOWQmJ1LI8Xee | shigeru |  20 |
| ewpSGf5URC1L1vPENbxh | takashi |  20 |
+----------------------+---------+-----+
```

### JSON

Structured JSON output, suitable for piping to tools like `jq`.

```sh
$ fscli --project-id my-project --out-mode json
```

**GET output:**

```json
{"id": "documentId", "data": {"name": "takashi", "age": 20}}
```

**QUERY output:**

```json
[
  {"id": "doc1", "data": {"name": "shigeru", "age": 20}},
  {"id": "doc2", "data": {"name": "takashi", "age": 20}}
]
```

## Non-Interactive Mode

fscli can be used in non-interactive mode by piping commands via stdin. This is useful for scripting and automation.

```sh
# Get a single document and extract a field
echo "GET users/user123" | fscli --project-id my-project --out-mode json | jq '.data.name'

# Query and get the ID of the first result
echo "QUERY users WHERE age > 25" | fscli --project-id my-project --out-mode json | jq '.[0].id'
```
