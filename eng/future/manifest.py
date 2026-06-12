#!/usr/bin/env python3
"""Normalized manifest of a publish dir: timestamp-insensitive hashes so
runs from different times compare equal unless content really changed."""
import hashlib, os, re, subprocess, sys

TS = [
    # ISO datetimes 2026-06-11T09:45:07.123-05:00 / with space / Z
    re.compile(rb"\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}(\.\d+)?([+-]\d{2}:?\d{2}|Z)?"),
    # Thu, Jun 11, 2026 09:45-0500  (kindling page footer style)
    re.compile(rb"(Mon|Tue|Wed|Thu|Fri|Sat|Sun), [A-Z][a-z]{2} \d{1,2}, \d{4} \d{2}:\d{2}([+-]\d{4})?"),
    # bare times like 09:45:07
    re.compile(rb"\b\d{2}:\d{2}:\d{2}\b"),
    # the build embeds the OS username in page footers ("Local Build (jmandel)")
    re.compile(rb"Local Build \([^)]*\)"),
    # render dates like "11 Jun 2026" (expansion-generated lines on valueset pages)
    re.compile(rb"\b\d{1,2} (Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) \d{4}\b"),
    # random UUIDs in generated html (table script ids, image names)
    re.compile(rb"[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"),
    # section numbers: unstable across identical stock runs (HashMap order in vs/cs numbering)
    re.compile(rb'sectioncount">[\d.]+'),
    re.compile(rb'<a name="[\d.]+">'),
    re.compile(rb'#[\d.]+" title="link to here"'),
]
ARCHIVE = (".zip", ".jar", ".pack", ".tgz", ".xlsx")
BINARY = (".png", ".gif", ".jpg", ".jpeg", ".ico", ".eot", ".woff", ".woff2", ".ttf", ".pdf", ".epub", ".exe", ".dll", ".class")

def norm_hash(path, sort_lines=False):
    try:
        with open(path, "rb") as f:
            data = f.read()
    except OSError as e:
        return "READ-ERROR"
    for rx in TS:
        data = rx.sub(b"TS", data)
    if sort_lines:
        # order-insensitive: for files whose only nondeterminism is element ordering, hash the
        # sorted line multiset - pure reordering compares equal, any content change still flags
        data = b"\n".join(sorted(data.split(b"\n")))
    return hashlib.sha256(data).hexdigest()[:16]

def xlsx_sig(path):
    # deep hash of member content with font-metric variance removed: POI auto-sizes column
    # widths using the JVM's installed fonts, so width attributes differ across machines while
    # the actual sheet content is identical. Any real content change still changes the hash.
    import zipfile
    h = hashlib.sha256()
    try:
        with zipfile.ZipFile(path) as z:
            for n in sorted(z.namelist()):
                data = z.read(n)
                data = re.sub(rb'width="[0-9.]+"', b'width="W"', data)
                for rx in TS:
                    data = rx.sub(b"TS", data)
                h.update(n.encode())
                h.update(data)
    except Exception:
        return "XLSX-ERROR"
    return h.hexdigest()[:16]

def archive_sig(path):
    # member names + sizes (zip stores mtimes which always differ). Sorted: archives are
    # written in filesystem-enumeration order, which differs across machines for identical
    # content - member ADD/REMOVE/SIZE changes still change the hash
    try:
        if path.endswith(".tgz"):
            lines = subprocess.run(["tar", "-tzf", path], capture_output=True).stdout.splitlines()
        else:
            p = subprocess.run(["unzip", "-l", path], capture_output=True)
            lines = [b" ".join(l.split()[:1] + l.split()[3:4]) for l in p.stdout.splitlines()[3:-2]]
        out = b"\n".join(sorted(lines))
    except Exception:
        return "ARCHIVE-ERROR"
    return hashlib.sha256(out).hexdigest()[:16]

def main(root, orderlist_path=None):
    ordered_tolerant = set()
    if orderlist_path:
        with open(orderlist_path) as f:
            ordered_tolerant = {l.strip() for l in f if l.strip()}
    # the build embeds its absolute checkout path in published links (htmldiff/jira); make
    # manifests location-independent by normalizing the path (raw and url-encoded forms)
    import urllib.parse
    checkout = os.path.dirname(os.path.abspath(root))
    for pat in (urllib.parse.quote(checkout, safe="").encode(), checkout.encode()):
        TS.append(re.compile(re.escape(pat)))
    files = []
    for dp, dn, fn in os.walk(root):
        for n in fn:
            files.append(os.path.relpath(os.path.join(dp, n), root))
    for rel in sorted(files):
        p = os.path.join(root, rel)
        low = rel.lower()
        if low.endswith(".xlsx"):
            h = "X:" + xlsx_sig(p)
        elif low.endswith(ARCHIVE):
            h = "A:" + archive_sig(p)
        elif low.endswith(BINARY):
            h = "B:" + hashlib.sha256(open(p, "rb").read()).hexdigest()[:16]
        elif rel in ordered_tolerant:
            h = "S:" + norm_hash(p, sort_lines=True)
        else:
            h = "N:" + norm_hash(p)
        print(f"{h}\t{rel}")

if __name__ == "__main__":
    main(sys.argv[1], sys.argv[2] if len(sys.argv) > 2 else None)
