# Clauses

## SELECT

Specify which fields to include in the output. Works with `QUERY` and `GET` operations.

```
SELECT <field1> [, <field2>, ...]
```

The document `ID` column is always included.

### Examples

```sql
QUERY users SELECT name
QUERY users SELECT name, age, email
GET users/abc123 SELECT name, age
```

## ORDER BY

Sort results by one or more fields. Works with `QUERY` operation only.

```
ORDER BY <field1> [ASC|DESC] [, <field2> [ASC|DESC], ...]
```

- **ASC** — Ascending order (default if omitted)
- **DESC** — Descending order

### Examples

```sql
-- Ascending (default)
QUERY users ORDER BY name

-- Descending
QUERY users ORDER BY age DESC

-- Multiple fields
QUERY users ORDER BY age ASC, name DESC
```

## LIMIT

Limit the number of returned documents. Works with `QUERY` operation only.

```
LIMIT <count>
```

### Examples

```sql
QUERY users LIMIT 10
QUERY users ORDER BY age DESC LIMIT 5
```

## Combining Clauses

Clauses can be combined in a single query:

```sql
QUERY users SELECT name WHERE age >= 20 ORDER BY name ASC LIMIT 10
```
