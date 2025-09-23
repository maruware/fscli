# fscli

A cli tool for firestore

## Demo

![fscli-demo](https://github.com/maruware/fscli/assets/1129887/887bbc4c-4f66-40a5-9211-256899abc067)

## Installation

Download binary from [Releases](https://github.com/maruware/fscli/releases)

or

```
go install github.com/maruware/fscli/cmd/fscli@latest
```

## Prepare

```sh
$ gcloud auth application-default login
```

## Usage

### Basic

```sh
$ fscli --project-id my-project
> QUERY users
+----------------------+---------+-----+
|          ID          |  name   | age |
+----------------------+---------+-----+
| VfsA2DjQOWQmJ1LI8Xee | shigeru |  20 |
| ewpSGf5URC1L1vPENbxh | takashi |  20 |
+----------------------+---------+-----+
> QUERY users WHERE age = 20
+----------------------+---------+-----+
|          ID          |  name   | age |
+----------------------+---------+-----+
| VfsA2DjQOWQmJ1LI8Xee | shigeru |  20 |
| ewpSGf5URC1L1vPENbxh | takashi |  20 |
+----------------------+---------+-----+
> QUERY users WHERE name = "takashi" AND age = 20;
+----------------------+---------+-----+
|          ID          |  name   | age |
+----------------------+---------+-----+
| ewpSGf5URC1L1vPENbxh | takashi |  20 |
+----------------------+---------+-----+
> GET users/ewpSGf5URC1L1vPENbxh
+----------------------+---------+-----+
|          ID          |  name   | age |
+----------------------+---------+-----+
| ewpSGf5URC1L1vPENbxh | takashi |  20 |
+----------------------+---------+-----+
> QUERY users SELECT name WHERE age = 20 ORDER BY name ASC LIMIT 1
+----------------------+---------+
|          ID          |  name   |
+----------------------+---------+
| VfsA2DjQOWQmJ1LI8Xee | shigeru |
+----------------------+---------+
> QUERY posts WHERE tags ARRAY_CONTAINS 'tech'
+----------------------+------------+----------------+
|          ID          |   title    |      tags      |
+----------------------+------------+----------------+
| nOsNxixUQ1rqNwVJz56O | First post | [tech finance] |
+----------------------+------------+----------------+
> COUNT users
2
> COUNT users WHERE name = "takashi"
1
> \d
users
posts
groups
> \d groups/yvPWOCcd4CvfUr2POuXk
members
```

### Non-Interactive Mode

`fscli` can be used in non-interactive mode by piping commands to it. This is useful for scripting and automation.

**Example: Get a single document and extract a field**
```sh
$ echo "GET users/user123" | fscli --project-id my-project --out-mode json | jq '.data.name'
"Test User"
```

**Example: Query a collection and get the ID of the first result**
```sh
$ echo "QUERY users WHERE age > 25" | fscli --project-id my-project --out-mode json | jq '.[0].id'
"user123"
```

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
