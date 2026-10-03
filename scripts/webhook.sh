#!/usr/bin/env bash
set -euo pipefail

host="${1:-github}"
url="${FUDA_URL:-http://localhost:8080}/api/webhooks/${host}"
body='{"ref":"refs/heads/develop","simulated":true}'
headers=(-H "Content-Type: application/json")

if [[ -n "${FUDA_WEBHOOK_SECRET:-}" ]]; then
  case "$host" in
    github) headers+=(-H "X-Hub-Signature-256: sha256=$(printf '%s' "$body" | openssl dgst -sha256 -hmac "$FUDA_WEBHOOK_SECRET" -r | cut -d' ' -f1)") ;;
    azure) headers+=(-H "X-Fuda-Secret: $FUDA_WEBHOOK_SECRET") ;;
  esac
fi

curl -sS -o /dev/null -w "POST $url -> %{http_code}\n" -X POST "${headers[@]}" --data "$body" "$url"
