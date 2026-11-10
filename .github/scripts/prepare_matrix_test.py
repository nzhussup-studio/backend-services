#!/usr/bin/env python3

import importlib.util
import pathlib
import unittest


SCRIPT = pathlib.Path(__file__).with_name("prepare_matrix.py")
SPEC = importlib.util.spec_from_file_location("prepare_matrix", SCRIPT)
assert SPEC and SPEC.loader
prepare_matrix = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(prepare_matrix)


SERVICES = [
    {"name": "image", "path": "image-service"},
    {"name": "account", "path": "account-service"},
]


class PrepareMatrixTests(unittest.TestCase):
    def test_selects_only_registered_changed_service(self) -> None:
        selected = prepare_matrix.select_services(
            SERVICES,
            ["image-service/internal/service/image.go", "README.md"],
        )
        self.assertEqual(["image"], [service["name"] for service in selected])

    def test_ignores_unregistered_directories(self) -> None:
        self.assertEqual([], prepare_matrix.select_services(SERVICES, ["other-service/main.go"]))

    def test_does_not_match_similar_path_prefixes(self) -> None:
        self.assertEqual([], prepare_matrix.select_services(SERVICES, ["image-service-old/main.go"]))


if __name__ == "__main__":
    unittest.main()
