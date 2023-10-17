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

### JSON mode

```sh
$ fscli --project-id my-project --out-mode json
> QUERY users
ID: VfsA2DjQOWQmJ1LI8Xee
Data: {"age":20,"name":"shigeru"}
--------------------------------------------------------------------------
ID: ewpSGf5URC1L1vPENbxh
Data: {"age":20,"name":"takashi"}
--------------------------------------------------------------------------
```
