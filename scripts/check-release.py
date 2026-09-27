#!/usr/bin/env python3
"""Check snapshot archives and smoke-test the native binary without publishing."""

import hashlib
import json
import os
from pathlib import Path
import platform
import subprocess
import tarfile
import tempfile
import zipfile


def main():
    dist = Path(__file__).resolve().parent.parent / "dist"
    checksums = {}
    for line in (dist / "checksums.txt").read_text().splitlines():
        digest, name = line.split(maxsplit=1)
        checksums[name.lstrip("*")] = digest

    native_os = {"Darwin": "darwin", "Linux": "linux", "Windows": "windows"}.get(platform.system())
    native_arch = {"arm64": "arm64", "aarch64": "arm64", "x86_64": "amd64", "AMD64": "amd64"}.get(platform.machine())
    native_binary = None
    for target_os in ("darwin", "linux", "windows"):
        for arch in ("arm64", "amd64"):
            extension = "zip" if target_os == "windows" else "tar.gz"
            archives = list(dist.glob(f"opensource_*_{target_os}_{arch}.{extension}"))
            if len(archives) != 1:
                raise ValueError(f"Expected one archive for {target_os}/{arch}, got {len(archives)}")
            archive = archives[0]
            if hashlib.sha256(archive.read_bytes()).hexdigest() != checksums.get(archive.name):
                raise ValueError(f"Missing or incorrect checksum: {archive.name}")
            binary = "opensource.exe" if target_os == "windows" else "opensource"
            if extension == "zip":
                with zipfile.ZipFile(archive) as bundle:
                    names = bundle.namelist()
                    data = bundle.read(binary)
            else:
                with tarfile.open(archive) as bundle:
                    names = bundle.getnames()
                    member = bundle.getmember(binary)
                    if not member.isfile() or not member.mode & 0o111:
                        raise ValueError(f"Binary is not executable: {archive.name}")
                    data = bundle.extractfile(member).read()
            if not {binary, "LICENSE", "README.md"}.issubset(names):
                raise ValueError(f"Missing binary or documentation: {archive.name}")
            if (target_os, arch) == (native_os, native_arch):
                native_binary = data
            print(f"Verified {archive.name}")

    if native_binary is None:
        raise ValueError("No release binary for this host's platform")

    version = json.loads((dist / "metadata.json").read_text())["version"]
    with tempfile.TemporaryDirectory(prefix="opensource-release-") as temp:
        root = Path(temp)
        binary = root / ("opensource.exe" if native_os == "windows" else "opensource")
        binary.write_bytes(native_binary)
        binary.chmod(0o755)
        home = root / "home"
        home.mkdir()
        project = root / "project"
        project.mkdir()
        env = dict(os.environ, HOME=str(home), USERPROFILE=str(home))

        def run(*args, **kwargs):
            return subprocess.run([str(binary), *args], cwd=project, env=env,
                                  check=True, capture_output=True, text=True, **kwargs).stdout

        if run("--version").strip() != version:
            raise ValueError("Embedded binary version does not match release metadata")
        run("--setup", input="Release Test\nrelease@example.com\n")
        run("--license", "mit", "--append", "README.md")
        license_text = (project / "LICENSE").read_text()
        if "Release Test <release@example.com>" not in license_text or "Permission is hereby granted" not in license_text:
            raise ValueError("Packaged binary failed license generation")
        if "## License" not in (project / "README.md").read_text():
            raise ValueError("Packaged binary failed README append")

    print("PASS: archives, checksums, version, setup, and license generation verified.")
    print("Signing, notarization, and Homebrew installation are checked by the tagged release workflow.")


if __name__ == "__main__":
    main()
