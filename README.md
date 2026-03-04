# Hive Adapter Scanner

## Overview

When creating custom adapters for [Hive](https://github.com/isar/hive), you need to assign them a numeric Type ID which is unique across your project.

When you have many adapters (particularly across multiple packages in a Melos workspace), it can be a bit of a headache keeping track of which ID has been assigned to various classes.

That's where this tool comes in; it'll scan your codebase and list the adapters you're using.

Example output:

```
+-------------------+--------------+---------+----------------------------------------------------------------+
| Adapter Class     | Generic Type | Type ID | File                                                           |
+-------------------+--------------+---------+----------------------------------------------------------------+
| UserAdapter       | User         | 1       | packages/auth/lib/src/storage/adapters/user_adapter.dart       |
| RoleAdapter       | Role         | 1       | packages/auth/lib/src/storage/adapters/role_adapter.dart       |
| BookAdapter       | Book         | 10      | packages/books/lib/storage/adapters/book_adapter.dart          |
| AuthorAdapter     | Author       | 11      | packages/books/lib/storage/adapters/author_adapter.dart        |
| CartAdapter       | Cart         | 20      | packages/ecomm/lib/storage/adapters/cart_adapter.dart          |
+-------------------+--------------+---------+----------------------------------------------------------------+
```

## Usage

Build it:

```bash
go build -o hive_adapter_scanner
```

Run it from your project directory:

```bash
./hive_adapter_scanner
```
