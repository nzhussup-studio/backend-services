#!/usr/bin/env python3

import json
import os
import subprocess
from pathlib import PurePosixPath
from typing import Any


ZERO_SHA = "0" * 40


def load_services(path: str = ".github/services.json") -> list[dict[str, Any]]:
    with open(path, "r", encoding="utf-8") as file:
        services = json.load(file)["services"]

    seen_paths: set[str] = set()
    for service in services:
        for field in ("name", "path", "deployment", "image_name"):
            if not service.get(field):
                raise ValueError(f"service entry is missing required field: {field}")
        service_path = service["path"].strip("/")
        if service_path in seen_paths:
            raise ValueError(f"duplicate service path: {service_path}")
        service["path"] = service_path
        seen_paths.add(service_path)
    return services


def diff_range() -> tuple[str, str]:
    event_name = os.environ.get("EVENT_NAME", "")
    head = os.environ.get("HEAD_SHA") or "HEAD"

    if event_name == "pull_request" and os.environ.get("BASE_SHA"):
        return f"{os.environ['BASE_SHA']}...{head}", "pull request"

    before = os.environ.get("BEFORE_SHA", "")
    if event_name == "push" and before and before != ZERO_SHA:
        return f"{before}..{head}", "push"

    # workflow_dispatch has no before SHA. Limit it to the latest commit so a
    # manual retry retains the same changed-service behavior as normal CI.
    return f"{head}^..{head}", "latest commit"


def changed_files(revision_range: str) -> list[str]:
    result = subprocess.run(
        ["git", "diff", "--name-only", "--diff-filter=ACMRTUXB", revision_range, "--"],
        check=True,
        capture_output=True,
        text=True,
    )
    return [line.strip() for line in result.stdout.splitlines() if line.strip()]


def service_changed(service_path: str, files: list[str]) -> bool:
    root = PurePosixPath(service_path)
    return any(PurePosixPath(file) == root or root in PurePosixPath(file).parents for file in files)


def select_services(services: list[dict[str, Any]], files: list[str]) -> list[dict[str, Any]]:
    return [service for service in services if service_changed(service["path"], files)]


def compact_json(value: Any) -> str:
    return json.dumps(value, separators=(",", ":"))


def main() -> None:
    services = load_services()
    revision_range, range_source = diff_range()
    files = changed_files(revision_range)
    changed = select_services(services, files)
    checked = [service for service in changed if service.get("ci_mode", "full") == "full"]
    deployments = [service["deployment"] for service in changed]

    print(f"Diff source: {range_source} ({revision_range})")
    print(f"Changed files: {compact_json(files)}")
    print(f"Selected services: {compact_json([service['name'] for service in changed])}")

    with open(os.environ["GITHUB_OUTPUT"], "a", encoding="utf-8") as output:
        output.write(f"build_changed={'true' if changed else 'false'}\n")
        output.write(f"build_matrix={compact_json(changed)}\n")
        output.write(f"checked_changed={'true' if checked else 'false'}\n")
        output.write(f"checked_matrix={compact_json(checked)}\n")
        output.write(f"deployments_joined={','.join(deployments)}\n")


if __name__ == "__main__":
    main()
