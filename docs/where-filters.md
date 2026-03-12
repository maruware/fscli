# WHERE Filters

The `WHERE` clause filters documents based on field values. It is supported by `QUERY` and `COUNT` operations.

```
WHERE <field> <operator> <value> [AND <field> <operator> <value> ...]
```

Multiple conditions can be combined with `AND`.

## Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `=` | Equal | `WHERE age = 20` |
| `!=` | Not equal | `WHERE status != "inactive"` |
| `>` | Greater than | `WHERE age > 25` |
| `>=` | Greater than or equal | `WHERE age >= 25` |
| `<` | Less than | `WHERE age < 30` |
| `<=` | Less than or equal | `WHERE age <= 30` |
| `IN` | Value in array | `WHERE status IN ["active", "pending"]` |
| `ARRAY_CONTAINS` | Array field contains value | `WHERE tags ARRAY_CONTAINS "tech"` |
| `ARRAY_CONTAINS_ANY` | Array field contains any of values | `WHERE tags ARRAY_CONTAINS_ANY ["tech", "design"]` |

## Value Types

### String

Quoted with single or double quotes.

```sql
WHERE name = "takashi"
WHERE name = 'takashi'
```

### Integer

```sql
WHERE age = 20
WHERE age > -5
```

### Float

```sql
WHERE score = 3.14
WHERE rating >= 4.5
```

### Timestamp

Use the `TIMESTAMP()` function. Supports three formats:

```sql
-- Date only
WHERE created_at > TIMESTAMP("2023-01-01")

-- Date and time
WHERE created_at > TIMESTAMP("2023-01-01T12:00:00")

-- RFC3339
WHERE created_at > TIMESTAMP("2023-01-01T12:00:00Z")
```

### Array

Square bracket syntax for `IN` and `ARRAY_CONTAINS_ANY` operators.

```sql
WHERE status IN ["active", "pending"]
WHERE tags ARRAY_CONTAINS_ANY [1, 2, 3]
```

Arrays can contain mixed types:

```sql
WHERE field IN [1, "text", 3.14]
```

## Document ID Filtering

Use the special field name `__id__` to filter by document ID.

```sql
-- Get a specific document by ID
QUERY users WHERE __id__ = "abc123"

-- Get multiple documents by IDs
QUERY users WHERE __id__ IN ["id1", "id2", "id3"]
```

## Multiple Conditions

Combine multiple conditions with `AND`:

```sql
QUERY users WHERE name = "takashi" AND age = 20
QUERY users WHERE age >= 20 AND age < 30 AND status = "active"
```
