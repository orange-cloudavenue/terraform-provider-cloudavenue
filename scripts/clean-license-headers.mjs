#!/usr/bin/env node
// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
//
// SPDX-License-Identifier: MPL-2.0

// Collapse piled/duplicated SPDX license headers into one canonical header.
//
// license-eye's `header fix` only ever *prepends* a header when the existing
// one does not match the configured pattern, which historically piled several
// headers on top of each file. This script removes every leading SPDX comment
// block (any style: `//`, `/* */`, `#`) and writes back exactly one canonical
// header per file, preserving shebangs, build tags and the earliest year.

import { readdirSync, readFileSync, writeFileSync, statSync } from 'node:fs';
import { join, extname, relative } from 'node:path';

const LICENSES = {
  'MPL-2.0': {
    identifier: 'MPL-2.0',
    // Minimal SPDX/REUSE best-practice header: copyright line, blank comment
    // line, identifier line. The LICENSE file at the repo root carries the
    // full text, so no distribution boilerplate is embedded per file.
    body: [
      'SPDX-FileCopyrightText: Copyright (c) {year} Orange',
      '',
      'SPDX-License-Identifier: MPL-2.0',
    ],
  },
  'Apache-2.0': {
    identifier: 'Apache-2.0',
    body: [
      'SPDX-FileCopyrightText: Copyright (c) {year} Orange',
      '',
      'SPDX-License-Identifier: Apache-2.0',
    ],
  },
};

// Comment style per extension family; canonical form chosen to match the
// dominant style in the org (line comments) and license-eye style ids.
const SLASH_EXTENSIONS = new Set(['.go', '.js', '.mjs', '.cjs', '.jsx', '.ts', '.tsx']);
const HASH_EXTENSIONS = new Set(['.yml', '.yaml', '.sh', '.bash', '.zsh', '.py', '.tf', '.hcl']);
const SPDX_MARKER = /SPDX-(FileCopyrightText|License-Identifier):/;
const YEAR_RE = /Copyright \(c\) (\d{4}(?:-\d{4})?)/;

const IGNORED_DIRS = new Set(['.git', 'node_modules', '.venv', 'venv', 'dist', 'build', '.terraform', '.taskfiles']);
const IGNORED_FILES = new Set(['LICENSE', 'LICENCE', 'NOTICE']);
const IGNORED_EXTENSIONS = new Set(['.md', '.json', '.svg', '.txt', '.png', '.jpg', '.gif', '.mod', '.sum', '.lock', '.tmpl', '.tpl']);
const IGNORED_BASENAMES = [/\.gitkeep$/];

function commentStyleFor(ext) {
  if (SLASH_EXTENSIONS.has(ext)) return 'slash';
  if (HASH_EXTENSIONS.has(ext)) return 'hash';
  return null;
}

// Parse consecutive comment lines starting at `i` as one header block.
// Returns { start, end, raw } (end = index after the block) or null.
function parseLineBlock(lines, i, prefix) {
  const raw = [];
  let j = i;
  while (j < lines.length && lines[j].startsWith(prefix)) {
    raw.push(lines[j]);
    j += 1;
  }
  return j > i ? { start: i, end: j, raw } : null;
}

// Parse a /* ... */ block starting at index i (line must start with "/*").
function parseBlockComment(lines, i) {
  const raw = [];
  let j = i;
  for (; j < lines.length; j += 1) {
    raw.push(lines[j]);
    if (lines[j].trim().endsWith('*/')) {
      j += 1;
      return { start: i, end: j, raw };
    }
  }
  return null; // unterminated: not a header
}

// Collect leading SPDX header blocks. Returns:
//   headerCount  — number of consecutive SPDX comment blocks found (>1 == piled)
//   years        — copyright years parsed from those blocks
//   keepRaw      — verbatim lines before the first SPDX line (shebang handled
//                  separately): sync banner, build tags, encoding lines
//   bodyStart    — index (in `lines`) just past the last consumed header
function scanHeaders(lines) {
  let i = 0;
  const years = [];
  let headerCount = 0;
  let firstSpdxLine = -1;

  while (i < lines.length) {
    const line = lines[i];
    if (line.trim() === '') {
      // A blank line only continues the scan if a header was already seen
      // and another comment block follows immediately after; otherwise stop.
      const next = i + 1;
      if (headerCount === 0 || next >= lines.length) break;
      if (!startsComment(lines[next])) break;
      i = next;
      continue;
    }
    let block = null;
    if (line.startsWith('//')) block = parseLineBlock(lines, i, '//');
    else if (line.startsWith('#')) block = parseLineBlock(lines, i, '#');
    else if (line.startsWith('/*')) block = parseBlockComment(lines, i);

    if (!block) break;
    const text = block.raw.join('\n');
    if (!SPDX_MARKER.test(text)) {
      // A non-license comment (sync banner, build tag) before the first
      // header is preserved and we keep scanning; after a header it ends the
      // scan so ordinary doc/build comments are never consumed.
      if (headerCount === 0) {
        i = block.end;
        while (i < lines.length && lines[i].trim() === '') i += 1;
        continue;
      }
      break;
    }
    if (firstSpdxLine === -1) {
      const idx = block.raw.findIndex((l) => SPDX_MARKER.test(l));
      firstSpdxLine = block.start + idx;
    }
    headerCount += 1;
    const year = YEAR_RE.exec(text);
    if (year) years.push(year[1]);
    i = block.end;
  }
  return { headerCount, years, bodyStart: i, keepRaw: firstSpdxLine < 0 ? [] : lines.slice(0, firstSpdxLine) };
}

function startsComment(line) {
  return line.startsWith('//') || line.startsWith('#') || line.startsWith('/*');
}

function canonicalHeader(style, licenseId, year) {
  const body = LICENSES[licenseId].body.map((l) => l.replace('{year}', year));
  if (style === 'hash') return body.map((l) => (l ? `# ${l}` : '#')).join('\n');
  return body.map((l) => (l ? `// ${l}` : '//')).join('\n');
}

// Keep the earliest creation year across consumed headers ("2023-2026" wins
// over "2026"); fall back to the current year for files with no year found.
function pickYear(years) {
  if (years.length === 0) return String(new Date().getFullYear());
  return years.reduce((a, b) => (Number(b.split('-')[0]) < Number(a.split('-')[0]) ? b : a));
}

function cleanFile(path, licenseId, apply) {
  const ext = extname(path).toLowerCase();
  const style = commentStyleFor(ext);
  if (!style) return null;

  const original = readFileSync(path, 'utf8');
  const lines = original.split('\n');

  let shebang = '';
  let start = 0;
  if (lines[0] && lines[0].startsWith('#!')) {
    shebang = lines[0];
    start = 1;
  }

  const { headerCount, years, bodyStart, keepRaw } = scanHeaders(lines.slice(start));
  if (headerCount === 0) return null; // no license header here: never touch

  // `bodyStart` / `keepRaw` are relative to the shebang-stripped slice.
  const rest = lines.slice(start + bodyStart);
  while (rest.length > 0 && rest[0].trim() === '') rest.shift();

  const rebuilt = [];
  if (shebang) rebuilt.push(shebang);
  rebuilt.push(...keepRaw); // sync banner, build tags — preserved verbatim
  rebuilt.push(canonicalHeader(style, licenseId, pickYear(years)));
  rebuilt.push('');
  rebuilt.push(...rest);
  const normalized = rebuilt.join('\n').replace(/\n+$/, '\n');

  const changed = normalized !== original;
  if (changed && apply) writeFileSync(path, normalized);
  return changed ? { path, headerCount } : null;
}

function* walk(dir) {
  for (const entry of readdirSync(dir)) {
    if (IGNORED_DIRS.has(entry) || entry.startsWith('.')) continue;
    const full = join(dir, entry);
    const st = statSync(full);
    if (st.isDirectory()) {
      yield* walk(full);
    } else if (st.isFile()) {
      if (IGNORED_FILES.has(entry)) continue;
      if (IGNORED_EXTENSIONS.has(extname(entry).toLowerCase())) continue;
      if (IGNORED_BASENAMES.some((re) => re.test(entry))) continue;
      if (entry.endsWith('.njk') || entry.endsWith('.njk.tmpl')) continue; // hub templates
      yield full;
    }
  }
}

function main(argv) {
  const args = argv.slice(2);
  const apply = !args.includes('--check');
  const licenseFlag = args.indexOf('--license');
  const licenseId = licenseFlag >= 0 ? args[licenseFlag + 1] : 'MPL-2.0';
  if (!LICENSES[licenseId]) {
    console.error(`unknown license id: ${licenseId} (supported: ${Object.keys(LICENSES).join(', ')})`);
    return 2;
  }
  const rootFlag = args.indexOf('--root');
  const root = rootFlag >= 0 ? args[rootFlag + 1] : '.';

  let cleaned = 0;
  let piled = 0;
  for (const file of walk(root)) {
    const res = cleanFile(file, licenseId, apply);
    if (!res) continue;
    cleaned += 1;
    if (res.headerCount > 1) piled += 1;
    console.log(`${apply ? 'cleaned' : 'would clean'}: ${relative(root, file)} (${res.headerCount} header block${res.headerCount > 1 ? 's' : ''})`);
  }
  console.log(`${apply ? 'cleaned' : 'found'} ${cleaned} file(s), ${piled} with piled headers`);
  // In --check mode a non-zero exit signals files still need cleaning.
  return apply ? 0 : (cleaned > 0 ? 1 : 0);
}

// Only run when invoked directly so tests can import the functions.
if (process.argv[1] && process.argv[1].endsWith('clean-license-headers.mjs')) {
  process.exit(main(process.argv));
}

export { scanHeaders, cleanFile, canonicalHeader, commentStyleFor };
