#!/bin/sh
set -e

KEY_FILE="/data/id_ed25519"

# Wait up to 30s for hbbs to generate the key on shared volume
MAX_WAIT=30
WAITED=0
while [ ! -f "$KEY_FILE" ] && [ $WAITED -lt $MAX_WAIT ]; do
    echo "Waiting for hbbs to generate key... (${WAITED}s)"
    sleep 3
    WAITED=$((WAITED + 3))
done

if [ -f "$KEY_FILE" ]; then
    KEY=$(cat "$KEY_FILE")
    echo "Using relay key from shared volume: $KEY_FILE"
else
    KEY="${RELAY_KEY}"
    if [ -n "$KEY" ]; then
        # hbbs doesn't write key file when -k is given a specific value,
        # so write it ourselves to the shared volume for future reference.
        echo "$KEY" > "$KEY_FILE"
        echo "Wrote RELAY_KEY to shared volume: $KEY_FILE"
    else
        echo "Key file not found after ${MAX_WAIT}s, and RELAY_KEY is empty"
    fi
fi

if [ -z "$KEY" ]; then
    echo "ERROR: No relay key available."
    echo "Ensure /data/id_ed25519 exists (shared with hbbs) or set RELAY_KEY."
    exit 1
fi

exec /usr/local/bin/hbbr -k "$KEY"
