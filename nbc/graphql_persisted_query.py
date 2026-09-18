import json

from mitmproxy import http


def request(flow):
    query = flow.request.query
    try:
        extensions = json.loads(query.get("extensions") or "null")
    except ValueError:
        return

    persisted = extensions.get("persistedQuery") if isinstance(extensions, dict) else None
    if not isinstance(persisted, dict):
        return

    digest = persisted.get("sha256Hash")
    if not isinstance(digest, str) or not digest:
        return

    persisted["sha256Hash"] = digest[:-1] + ("0" if digest[-1] != "0" else "1")
    query["extensions"] = json.dumps(extensions, separators=(",", ":"))
