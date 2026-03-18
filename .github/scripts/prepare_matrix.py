#!/usr/bin/env python3

import json


with open(".github/services.json", "r", encoding="utf-8") as f:
    services = json.load(f)["services"]

checked = [service for service in services if service.get("ci_mode", "full") == "full"]
deployments = [service["deployment"] for service in services]

from os import environ

with open(environ["GITHUB_OUTPUT"], "a", encoding="utf-8") as f:
    f.write(f"build_matrix={json.dumps(services, separators=(',', ':'))}\n")
    f.write(f"checked_changed={'true' if checked else 'false'}\n")
    f.write(f"checked_matrix={json.dumps(checked, separators=(',', ':'))}\n")
    f.write(f"deployments_joined={' '.join(deployments)}\n")
