#!/usr/bin/env python3
# SPDX-FileCopyrightText: Copyright (c) 2026 Orange
# SPDX-License-Identifier: MPL-2.0
#
# This software is distributed under the MPL-2.0 license.
# the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
# or see the "LICENSE" file for more details.

# SPDX-FileCopyrightText: Copyright (c) 2026 Orange
# SPDX-License-Identifier: Mozilla Public License 2.0
#
# This software is distributed under the MPL-2.0 license.
# the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
# or see the "LICENSE" file for more details.

"""Remove duplicate leading license headers from files.

Usage:
  python3 scripts/dedupe_license_headers.py FILE [FILE ...]
  python3 scripts/dedupe_license_headers.py DIR --recursive
  python3 scripts/dedupe_license_headers.py PATH --check

Detects repeated leading license header blocks at start of file, keeps first
effective block, removes later duplicates. Supports:
  1) line comment SPDX header blocks using // comments
  2) block comment MPL header blocks using /* ... */
"""

from __future__ import annotations

import argparse
import sys
from pathlib import Path


SKIP_DIR_NAMES = {".git", ".opencode", ".kilo"}


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Collapse duplicate leading license header blocks at file start."
    )
    parser.add_argument("paths", nargs="+", help="Files or directories to process")
    parser.add_argument(
        "-r",
        "--recursive",
        action="store_true",
        help="Recurse into directories",
    )
    parser.add_argument(
        "--check",
        action="store_true",
        help="Report files needing changes without rewriting",
    )
    return parser.parse_args()


def iter_targets(paths: list[str], recursive: bool) -> list[Path]:
    targets: list[Path] = []

    for raw_path in paths:
        path = Path(raw_path)

        if not path.exists():
            print(f"skip missing path: {path}", file=sys.stderr)
            continue

        if path.is_file():
            if any(part in SKIP_DIR_NAMES for part in path.parts):
                continue
            targets.append(path)
            continue

        if not path.is_dir():
            print(f"skip unsupported path: {path}", file=sys.stderr)
            continue

        iterator = path.rglob("*") if recursive else path.iterdir()
        for child in iterator:
            if child.is_file():
                if any(part in SKIP_DIR_NAMES for part in child.parts):
                    continue
                targets.append(child)

    return sorted(set(targets))


def count_leading_blank_lines(text: str) -> int:
    count = 0
    for line in text.splitlines(keepends=True):
        if line.strip() != "":
            break
        count += len(line)
    return count


def consume_line_header(text: str, start: int) -> tuple[str, int] | None:
    rest = text[start:]
    lines = rest.splitlines(keepends=True)
    if len(lines) < 2:
        return None

    first = lines[0]
    second = lines[1]

    if not first.startswith("// SPDX-FileCopyrightText:"):
        return None
    if not second.startswith("// SPDX-License-Identifier:"):
        return None

    consumed = len(first) + len(second)

    saw_mpl_explanatory_line = False
    for line in lines[2:]:
        stripped = line.strip()

        if stripped == "":
            consumed += len(line)
            continue

        if stripped == "//":
            consumed += len(line)
            continue

        if stripped.startswith("// This software is distributed under the MPL-2.0 license."):
            saw_mpl_explanatory_line = True
            consumed += len(line)
            continue

        if saw_mpl_explanatory_line:
            if stripped.startswith("// SPDX-"):
                break
            if stripped.startswith("// the text of which is available at "):
                consumed += len(line)
                continue
            if stripped.startswith('// or see the "LICENSE" file for more details.'):
                consumed += len(line)
                continue
            if stripped == "//":
                consumed += len(line)
                continue
            break

        break

    return text[start : start + consumed], start + consumed


def consume_block_header(text: str, start: int) -> tuple[str, int] | None:
    rest = text[start:]
    if not rest.startswith("/*"):
        return None

    end = rest.find("*/")
    if end == -1:
        return None

    block_end = end + 2
    if block_end < len(rest) and rest[block_end:block_end + 1] == "\r":
        block_end += 1
    if block_end < len(rest) and rest[block_end:block_end + 1] == "\n":
        block_end += 1

    block = rest[:block_end]
    if "SPDX-FileCopyrightText:" not in block:
        return None
    if (
        "SPDX-License-Identifier: Mozilla Public License 2.0" not in block
        and "SPDX-License-Identifier: MPL-2.0" not in block
    ):
        return None

    return block, start + block_end


def consume_header(text: str, start: int) -> tuple[str, int] | None:
    for consumer in (consume_line_header, consume_block_header):
        match = consumer(text, start)
        if match is not None:
            return match
    return None


def dedupe_leading_headers(text: str) -> str:
    pos = 0
    headers: list[str] = []

    while True:
        blank_prefix = count_leading_blank_lines(text[pos:])
        pos += blank_prefix

        match = consume_header(text, pos)
        if match is None:
            pos -= blank_prefix
            break

        header, pos = match
        headers.append(header)

    if len(headers) <= 1:
        return text

    first_header = headers[0]
    normalized = first_header.rstrip("\r\n") + "\n\n"
    remainder = text[pos:].lstrip("\r\n")

    return normalized + remainder


def process_file(path: Path, check: bool) -> bool:
    try:
        original = path.read_text(encoding="utf-8")
    except UnicodeDecodeError:
        print(f"skip non-utf8 file: {path}", file=sys.stderr)
        return False

    updated = dedupe_leading_headers(original)
    if updated == original:
        return False

    if check:
        print(path)
        return True

    path.write_text(updated, encoding="utf-8")
    print(f"updated {path}")
    return True


def main() -> int:
    args = parse_args()
    changed = False

    for path in iter_targets(args.paths, args.recursive):
        changed = process_file(path, args.check) or changed

    return 1 if args.check and changed else 0


if __name__ == "__main__":
    raise SystemExit(main())
