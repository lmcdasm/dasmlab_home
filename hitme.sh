#!/bin/bash
# a simple testing script for pulling a website 
CHROME_BIN=$(command -v google-chrome-stable || command -v google-chrome || command -v chromium || command -v chromium-browser); \
for i in {1..30}; do
  $CHROME_BIN --headless=new --incognito --disable-cache \
  --user-data-dir="$(mktemp -d)" --window-size=1366,768 \
  --virtual-time-budget=8000 --screenshot='/tmp/shot-'$i'.png' 'https://dasmlab.org'
done

