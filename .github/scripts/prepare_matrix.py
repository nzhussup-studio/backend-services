#!/usr/bin/env python3

import json
import os


with open(".github/services.json", "r", encoding="utf-8") as f:
    services = json.load(f)["services"]

changed_files = [
    line.strip()
    for line in os.environ.get("CHANGED_FILES", "").splitlines()
    if line.strip()
]

full_rebuild_paths = (
    ".github/services.json",
    ".github/templates/",
    ".github/scripts/",
    ".github/workflows/cicd.yml",
)

if os.environ.get("GITHUB_EVENT_NAME") == "workflow_dispatch":
    selected = services
elif any(
    path == ".github/services.json"
    or any(path.startswith(prefix) for prefix in full_rebuild_paths[1:])
    for path in changed_files
):
    selected = services
else:
    selected = []
    for service in services:
        prefix = service["path"].rstrip("/") + "/"
        if any(path.startswith(prefix) for path in changed_files):
            selected.append(service)

checked = [service for service in selected if service.get("ci_mode", "full") == "full"]
deployments = [service["deployment"] for service in selected]

with open(os.environ["GITHUB_OUTPUT"], "a", encoding="utf-8") as f:
    f.write(f"build_changed={'true' if selected else 'false'}\n")
    f.write(f"build_matrix={json.dumps(selected, separators=(',', ':'))}\n")
    f.write(f"checked_changed={'true' if checked else 'false'}\n")
    f.write(f"checked_matrix={json.dumps(checked, separators=(',', ':'))}\n")
    f.write(f"deployments={json.dumps(deployments, separators=(',', ':'))}\n")
    f.write(f"deployments_joined={' '.join(deployments)}\n")
