# Meta Commands

Meta commands start with `\` and provide utility functions outside of Firestore queries.

## \d — List Collections

List collections at the root level or subcollections of a specific document.

```
\d [document_path]
```

### Examples

```
-- List root-level collections
> \d
users
posts
groups

-- List subcollections of a document
> \d groups/yvPWOCcd4CvfUr2POuXk
members
```

## \pager — Toggle Output Paging

Enable or disable paging for large output. Uses the `$PAGER` environment variable, or falls back to `less`.

```
\pager on
\pager off
```

### Examples

```
-- Enable pager
> \pager on

-- Disable pager
> \pager off
```
