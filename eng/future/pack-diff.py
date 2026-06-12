#!/usr/bin/env python3
"""Canonical diff of two terminology answer packs (dirs or zips).

Pack entries legitimately contain volatile fields (empirically derived by diffing two
recordings of the same server: millisecond step-timings inside diagnostics strings,
expansion.identifier UUIDs, expansion.timestamp). Pack *identity* stays the raw content
hash; pack *equality* for refresh decisions is this canonical comparison, so a re-recording
that changed nothing semantic reports "no change" and the refresh job keeps the old pack.

Usage: pack-diff.py OLD NEW [--markdown]
Exit 0 = canonically identical; exit 3 = real differences (markdown summary on stdout).
"""
import io, json, re, sys, zipfile
from pathlib import Path

VOLATILE = [
    (re.compile(r"\b\d+ms "), "Nms "),                     # server step timings in diagnostics
    (re.compile(r"urn:uuid:[0-9a-f-]{36}"), "urn:uuid:X"), # expansion identifiers
    (re.compile(r"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?(Z|[+-]\d{2}:?\d{2})?"), "TS"),
]
ENTRY_MARKER = "-------------------------------------------------------------------------------------"
BREAK = "####"

def canon(text):
    text = text.replace("\r\n", "\n")  # zip reads preserve CRLF; dir text reads translate it
    for rx, sub in VOLATILE:
        text = rx.sub(sub, text)
    return text

def load(path):
    """page name -> {canonical request -> canonical response}"""
    p = Path(path)
    pages = {}
    def add(name, text):
        if not name.endswith(".cache") or name.startswith("."):
            return
        entries = {}
        for chunk in text.split(ENTRY_MARKER):
            chunk = chunk.strip()
            if not chunk or BREAK not in chunk:
                continue
            req, resp = chunk.split(BREAK, 1)
            entries[canon(req.strip())] = canon(resp.strip())
        pages[name] = entries
    if p.is_dir():
        for f in sorted(p.iterdir()):
            add(f.name, f.read_text(errors="replace")) if f.is_file() else None
    else:
        with zipfile.ZipFile(p) as z:
            for n in sorted(z.namelist()):
                add(Path(n).name, z.read(n).decode(errors="replace"))
    return pages

def main():
    old, new = load(sys.argv[1]), load(sys.argv[2])
    added, removed, changed = [], [], []
    for page in sorted(set(old) | set(new)):
        o, n = old.get(page, {}), new.get(page, {})
        for req in n.keys() - o.keys():
            added.append((page, req))
        for req in o.keys() - n.keys():
            removed.append((page, req))
        for req in o.keys() & n.keys():
            if o[req] != n[req]:
                changed.append((page, req, o[req], n[req]))
    if not (added or removed or changed):
        print("canonically identical (volatile fields normalized: step timings, expansion ids/timestamps)")
        return 0
    print(f"## Answer pack changes\n")
    print(f"| | count |\n|---|---|\n| added | {len(added)} |\n| removed | {len(removed)} |\n| changed | {len(changed)} |\n")
    def head(req):
        return " ".join(req.split())[:140]
    if changed:
        print("### Changed answers (the part that needs review)\n")
        for page, req, o, n in changed[:25]:
            print(f"- **{page}**: `{head(req)}`")
            print(f"  - was: `{head(o)}`")
            print(f"  - now: `{head(n)}`")
        if len(changed) > 25:
            print(f"- … and {len(changed)-25} more")
    if added:
        print("\n### Added\n")
        for page, req in added[:15]:
            print(f"- {page}: `{head(req)}`")
        if len(added) > 15:
            print(f"- … and {len(added)-15} more")
    if removed:
        print("\n### Removed\n")
        for page, req in removed[:15]:
            print(f"- {page}: `{head(req)}`")
    return 3

if __name__ == "__main__":
    sys.exit(main())
