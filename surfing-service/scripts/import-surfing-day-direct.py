#!/usr/bin/env python3
"""Bulk-import local media into surfing-service PVC via preload + media copy."""
import hashlib
import json
import os
import subprocess
import sys
from datetime import datetime, timezone
from pathlib import Path

NS = os.environ.get("SURFING_NAMESPACE", "surfing-service-system")
KUBECONFIG = os.environ.get("KUBECONFIG") or os.path.expanduser("~/.kube/config")
OC = ["oc", f"--kubeconfig={KUBECONFIG}", "-n", NS]

IMAGE_EXTS = {".jpg", ".jpeg", ".png", ".gif", ".webp", ".heic", ".heif"}
VIDEO_EXTS = {".mp4", ".mov", ".m4v", ".webm", ".avi", ".mkv"}


def run(cmd, check=True):
    return subprocess.run(cmd, check=check)


def oc(*args, check=True, capture=False):
    kw = {"check": check, "text": True}
    if capture:
        kw["capture_output"] = True
    return subprocess.run(OC + list(args), **kw)


def pod_name():
    r = oc("get", "pods", "-l", "app=surfing-service", "-o", "jsonpath={.items[0].metadata.name}", capture=True)
    name = r.stdout.strip()
    if not name:
        raise SystemExit("no surfing-service pod found")
    return name


def stable_id(kind, *parts):
    raw = f"preload:{kind}:" + ":".join(p.lower() for p in parts)
    return hashlib.sha256(raw.encode()).hexdigest()[:32]


def media_type(path: Path):
    ext = path.suffix.lower()
    if ext in IMAGE_EXTS:
        return "image", ext
    if ext in VIDEO_EXTS:
        return "video", ext
    return None, None


def load_manifest(pod):
    r = oc("exec", pod, "--", "cat", "/data/manifest.json", check=False, capture=True)
    if r.returncode != 0 or not r.stdout.strip():
        return []
    return json.loads(r.stdout)


def main():
    if len(sys.argv) < 3:
        raise SystemExit("usage: import-surfing-day-direct.py TITLE SOURCE_DIR [date] [location]")

    title = sys.argv[1]
    source = Path(sys.argv[2])
    date = sys.argv[3] if len(sys.argv) > 3 else "2026-01-01"
    location = sys.argv[4] if len(sys.argv) > 4 else "Bonaire"

    if not source.is_dir():
        raise SystemExit(f"not a directory: {source}")

    pod = pod_name()
    preload = f"/data/preload/{title}"
    print(f"pod={pod} preload={preload}")

    oc("exec", pod, "--", "mkdir", "-p", preload)
    pod = pod_name()  # refresh in case of restart
    print("rsync source -> preload ...")
    run(["oc", f"--kubeconfig={KUBECONFIG}", "rsync", f"{source}/", f"{NS}/{pod}:{preload}/"])

    meta = {"title": title, "date": date, "location": location}
    meta_tmp = Path("/tmp/surfing-day.json")
    meta_tmp.write_text(json.dumps(meta, indent=2))
    run(OC + ["cp", str(meta_tmp), f"{pod}:{preload}/day.json"])
    meta_tmp.unlink(missing_ok=True)

    days = load_manifest(pod)
    day_id = stable_id("day", title)
    day = next((d for d in days if d.get("id") == day_id), None)
    if not day:
        day = {
            "id": day_id,
            "title": title,
            "date": date,
            "location": location,
            "created_at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
            "media": [],
        }
        days.append(day)

    existing = {m["id"] for m in day.get("media", [])}
    files = sorted(p for p in source.iterdir() if p.is_file())
    print(f"processing {len(files)} files ...")

    shell_lines = ["#!/bin/sh", "set -e", "mkdir -p /data/media"]
    added = 0
    for path in files:
        mtype, ext = media_type(path)
        if not mtype:
            continue
        media_id = stable_id("media", title, path.name)
        if media_id in existing:
            continue
        src = f"{preload}/{path.name}"
        dst = f"/data/media/{media_id}{ext}"
        shell_lines.append(f'cp "{src}" "{dst}"')
        day["media"].append({
            "id": media_id,
            "filename": path.name,
            "media_type": mtype,
            "url": f"/serve?id={media_id}",
            "created_at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        })
        existing.add(media_id)
        added += 1

    script_tmp = Path("/tmp/surfing-import.sh")
    script_tmp.write_text("\n".join(shell_lines) + "\n")
    run(OC + ["cp", str(script_tmp), f"{pod}:/tmp/import.sh"])
    script_tmp.unlink(missing_ok=True)
    print(f"copying {added} files into /data/media ...")
    oc("exec", pod, "--", "sh", "/tmp/import.sh")

    manifest_tmp = Path("/tmp/surfing-manifest.json")
    manifest_tmp.write_text(json.dumps(days, indent=2))
    run(OC + ["cp", str(manifest_tmp), f"{pod}:/data/manifest.json"])
    manifest_tmp.unlink(missing_ok=True)

    print(f"imported {added} file(s); restarting surfing-service ...")
    oc("rollout", "restart", "deployment/surfing-service")
    oc("rollout", "status", "deployment/surfing-service", "--timeout=600s")
    print("done")


if __name__ == "__main__":
    main()
