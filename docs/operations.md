# Operations

fscli supports three operations: `QUERY`, `GET`, and `COUNT`.

## QUERY

Query documents in a collection. Supports `SELECT`, `WHERE`, `ORDER BY`, and `LIMIT` clauses.

```
QUERY <collection_path> [SELECT ...] [WHERE ...] [ORDER BY ...] [LIMIT ...]
QUERY COLLECTION_GROUP <collection_id> [SELECT ...] [WHERE ...] [ORDER BY ...] [LIMIT ...]
```

### Examples

```sql
-- Query all documents
QUERY users

-- Query with filter
QUERY users WHERE age = 20

-- Subcollection query
QUERY users/abc123/posts

-- Collection group query
QUERY COLLECTION_GROUP posts WHERE title = "post-1-1"
```

## GET

Fetch a single document by its full path.

```
GET <document_path> [SELECT ...]
```

The document path must have an odd number of segments (e.g., `collection/docId`).

### Examples

```sql
-- Get a single document
GET users/ewpSGf5URC1L1vPENbxh

-- Get with field selection
GET users/ewpSGf5URC1L1vPENbxh SELECT name, age
```

## COUNT

Count documents in a collection. Supports `WHERE` clause for filtering.

```
COUNT <collection_path> [WHERE ...]
COUNT COLLECTION_GROUP <collection_id> [WHERE ...]
```

Returns a single integer.

### Examples

```sql
-- Count all documents
COUNT users

-- Count with filter
COUNT users WHERE name = "takashi"

-- Collection group count
COUNT COLLECTION_GROUP posts WHERE title = "post-1-1"
```

## Collection Path

Collection paths support nested subcollections using the format:

```
collection/docId/subcollection/docId/subcollection/...
```

- **Collection path** (even number of segments): `users`, `users/abc123/posts`
- **Document path** (odd number of segments): `users/abc123`, `users/abc123/posts/xyz789`

Leading slashes are automatically stripped.

## Collection Group

`COLLECTION_GROUP` targets all collections with the same collection ID across the database hierarchy.

- Use a single collection ID (for example, `posts`)
- Slash-separated paths are not allowed (for example, `users/posts` is invalid)
- Supported with `QUERY` and `COUNT`
