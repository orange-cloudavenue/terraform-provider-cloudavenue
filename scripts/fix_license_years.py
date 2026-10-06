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

"""Rewrite SPDX copyright years from git history.

Usage:
  python3 scripts/fix_license_years.py FILE [FILE ...]
  python3 scripts/fix_license_years.py DIR --recursive
  python3 scripts/fix_license_years.py PATH --check

For each tracked file, script reads first and latest commit years from git log,
computes YYYY or YYYY-YYYY, then rewrites first matching SPDX copyright line in
line-comment or block-comment form.
"""

from __future__ import annotations

import argparse
import re
import subprocess
import sys
from pathlib import Path


SKIP_DIR_NAMES = {".git", ".opencode", ".kilo"}
COPYRIGHT_RE = re.compile(
    r"^(?P<prefix>\s*(?://|\*))\s*SPDX-FileCopyrightText:\s+Copyright \(c\) "
    r"(?P<years>\d{4}(?:-\d{4})?) Orange(?P<suffix>\s*)$",
    re.MULTILINE,
)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Rewrite SPDX copyright years from git history."
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


def git_root() -> Path:
    result = subprocess.run(
        ["git", "rev-parse", "--show-toplevel"],
        check=True,
        capture_output=True,
        text=True,
    )
    return Path(result.stdout.strip())


def is_skipped(path: Path) -> bool:
    return any(part in SKIP_DIR_NAMES for part in path.parts)


def iter_targets(root: Path, paths: list[str], recursive: bool) -> list[Path]:
    targets: list[Path] = []

    for raw_path in paths:
        path = Path(raw_path)

        if not path.exists():
            print(f"skip missing path: {path}", file=sys.stderr)
            continue

        if path.is_file():
            if not is_skipped(path):
                targets.append(path.resolve())
            continue

        if not path.is_dir():
            print(f"skip unsupported path: {path}", file=sys.stderr)
            continue

        iterator = path.rglob("*") if recursive else path.iterdir()
        for child in iterator:
            if child.is_file() and not is_skipped(child):
                targets.append(child.resolve())

    root_resolved = root.resolve()
    return sorted(
        {
            target
            for target in targets
            if target.is_relative_to(root_resolved)
        }
    )


def is_tracked(root: Path, path: Path) -> bool:
    rel = path.relative_to(root)
    result = subprocess.run(
        ["git", "ls-files", "--error-unmatch", str(rel)],
        cwd=root,
        capture_output=True,
        text=True,
    )
    return result.returncode == 0


def get_year_token(root: Path, path: Path) -> str | None:
    rel = path.relative_to(root)
    result = subprocess.run(
        ["git", "log", "--follow", "--format=%ad", "--date=format:%Y", "--", str(rel)],
        cwd=root,
        check=True,
        capture_output=True,
        text=True,
    )
    years = [line.strip() for line in result.stdout.splitlines() if line.strip()]
    if not years:
        return None

    newest = years[0]
    oldest = years[-1]
    if newest == oldest:
        return newest
    return f"{oldest}-{newest}"


def rewrite_first_spdx_line(text: str, year_token: str) -> str:
    def replace(match: re.Match[str]) -> str:
        return (
            f"{match.group('prefix')} SPDX-FileCopyrightText: Copyright (c) "
            f"{year_token} Orange{match.group('suffix')}"
        )

    return COPYRIGHT_RE.sub(replace, text, count=1)


def process_file(root: Path, path: Path, check: bool) -> bool:
    if not is_tracked(root, path):
        return False

    year_token = get_year_token(root, path)
    if year_token is None:
        return False

    try:
        original = path.read_text(encoding="utf-8")
    except UnicodeDecodeError:
        print(f"skip non-utf8 file: {path}", file=sys.stderr)
        return False

    updated = rewrite_first_spdx_line(original, year_token)
    if updated == original:
        return False

    rel = path.relative_to(root)
    if check:
        print(rel)
        return True

    path.write_text(updated, encoding="utf-8")
    print(f"updated {rel}")
    return True


def main() -> int:
    args = parse_args()
    root = git_root()
    changed = False

    for path in iter_targets(root, args.paths, args.recursive):
        changed = process_file(root, path, args.check) or changed

    return 1 if args.check and changed else 0


if __name__ == "__main__":
    raise SystemExit(main())
