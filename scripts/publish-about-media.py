#!/usr/bin/env python3
"""Publish About WebPs to Cloudflare R2 (bucket dasmlab-home), Surfing-style."""
from __future__ import annotations

import sys
from pathlib import Path

import boto3
from botocore.config import Config

ROOT = Path(__file__).resolve().parents[1]
SRC_DIR = ROOT / "public" / "media" / "hero"
BUCKET = "dasmlab-home"
PREFIX = "about"
CREDS = Path.home() / "r2_creds_dasmlab_surfing"
PUBLIC_BASE = "https://pub-29bde7a836c744729bebe74bfd4008a2.r2.dev/about"


def load_env(path: Path) -> dict[str, str]:
    env: dict[str, str] = {}
    for line in path.read_text().splitlines():
        line = line.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        key, value = line.split("=", 1)
        env[key.strip()] = value.strip().strip('"').strip("'")
    return env


def main() -> int:
    if not CREDS.is_file():
        print(f"ERROR: missing creds file {CREDS}", file=sys.stderr)
        return 1
    env = load_env(CREDS)
    files = sorted(SRC_DIR.glob("*.webp"))
    if not files:
        print(f"ERROR: no WebP files in {SRC_DIR}", file=sys.stderr)
        return 1

    client = boto3.client(
        "s3",
        endpoint_url=f"https://{env['R2_ACCOUNT_ID']}.r2.cloudflarestorage.com",
        aws_access_key_id=env["R2_ACCESS_KEY_ID"],
        aws_secret_access_key=env["R2_SECRET_ACCESS_KEY"],
        region_name="auto",
        config=Config(signature_version="s3v4"),
    )

    for path in files:
        key = f"{PREFIX}/{path.name}"
        client.put_object(
            Bucket=BUCKET,
            Key=key,
            Body=path.read_bytes(),
            ContentType="image/webp",
            CacheControl="public, max-age=31536000, immutable",
        )
        print(f"put s3://{BUCKET}/{key} ({path.stat().st_size} bytes)")
        print(f"    {PUBLIC_BASE}/{path.name}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
