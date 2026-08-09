#!/usr/bin/env python3

import argparse
import glob
import os
import platform
import shutil
import subprocess
import sys
from pathlib import Path


class BuildSystem:
    def __init__(self):
        self.wd = Path.cwd()
        self.dst = self.wd / "demo"

        # Set up environment
        self._setup_environment()
        self._determine_pkgconfig_path()

    def _setup_environment(self):
        """Set up environment variables"""
        os.environ["CGO_ENABLED"] = "1"
        os.environ["CC"] = "clang"
        os.environ["CXX"] = "clang++"
        os.environ["CGO_CFLAGS"] = "-O3 -g"
        os.environ["CGO_LDFLAGS"] = "-lc++"
        os.environ["CGO_CXXFLAGS"] = "-std=c++20 -stdlib=libc++"
        if platform.system() == "Darwin":
            flags = "-framework CoreGraphics -framework CoreFoundation -framework Foundation -framework Quartz -framework QuartzCore -lc++"
            os.environ["LDFLAGS"] = flags
            os.environ["CGO_LDFLAGS"] = flags

    def _determine_pkgconfig_path(self):
        """Determine the appropriate PKGCONFIG path based on platform"""
        os_name = platform.system()
        os_arch = platform.machine()
        print(f"Detected OS: {os_name}, Architecture: {os_arch}")
        if os_name == "Windows":
            if os_arch == "AMD64":
                self.pkgconfig = self.wd / "lib/pdfium/pdfium-win-x64"
            else:
                self.pkgconfig = self.wd / "lib/pdfium/pdfium-win-x86"
        elif os_name == "Linux":
            if os_arch == "x86_64":
                self.pkgconfig = self.wd / "lib/pdfium/pdfium-linux-x64"
            elif os_arch == "aarch64":
                self.pkgconfig = self.wd / "lib/pdfium/pdfium-linux-arm64"
            elif os_arch.startswith("arm"):
                self.pkgconfig = self.wd / "lib/pdfium/pdfium-linux-arm"
            elif os_arch.find("86") != -1:
                self.pkgconfig = self.wd / "lib/pdfium/pdfium-linux-x86"
            else:
                print(f"Unsupported architecture: {os_arch}")
                sys.exit(1)
        elif os_name == "Darwin":
            if os_arch == "arm64":
                self.pkgconfig = self.wd / "lib/pdfium/pdfium-mac-arm64"
            else:
                self.pkgconfig = self.wd / "lib/pdfium/pdfium-mac-x64"

        os.environ["PKG_CONFIG_PATH"] = str(self.pkgconfig)

    def _get_build_flags(self):
        """Get Go build flags"""
        return [
            "-ldflags",
            "-s -w",
            "-trimpath",
            "-v",
        ]

    def fix_pc_files(self):
        """Rewrite .pc files in PKGCONFIG directories with absolute prefix paths"""
        # Find all pdfium directories
        pdfium_dirs = glob.glob(str(self.wd / "lib/pdfium/pdfium-*"))

        for pdfium_dir in pdfium_dirs:
            pdfium_path = Path(pdfium_dir)
            pc_files = list(pdfium_path.glob("*.pc"))

            for pc_file in pc_files:
                # Read the current content
                with open(pc_file, "r") as f:
                    content = f.read()

                # Replace prefix line with absolute path
                lines = content.split("\n")
                new_lines = []

                for line in lines:
                    if line.startswith("prefix="):
                        new_lines.append(f"prefix={pdfium_path.absolute()}")
                    else:
                        new_lines.append(line)

                # Write back the modified content
                with open(pc_file, "w") as f:
                    f.write("\n".join(new_lines))

    def run_command(self, cmd, **kwargs):
        """Run a shell command"""
        if isinstance(cmd, list):
            print(f"Running: {' '.join(cmd)}")
        else:
            print(f"Running: {cmd}")

        if isinstance(cmd, str):
            cmd = cmd.split()

        result = subprocess.run(cmd, **kwargs)
        if result.returncode != 0:
            print(f"Command failed with exit code {result.returncode}")
            sys.exit(result.returncode)
        return result

    def help(self):
        """Display help"""
        print("Usage:")
        print("  python build.py <command>")
        print()
        print("Commands:")
        print("  help          - displays help")
        print("  env           - show environment variables")
        print("  tidy          - formats and updates the go.mod file")
        print("  desktop       - compiles the local backend for the desktop platform")
        print("  android       - compiles the local backend for the android platform")
        print("  ios           - compiles the local backend for the ios platform")
        print("  cross         - compiles the desktop backend for many platforms")
        print("  fix-pc        - fix .pc files with absolute paths")

    def env(self):
        """Show environment variables"""
        build_flags = self._get_build_flags()
        pdfium_cflags = ""
        pdfium_ldflags = ""
        try:
            result = subprocess.run(["pkg-config", "--cflags", "pdfium"], capture_output=True, text=True)
            pdfium_cflags = result.stdout.strip()
        except Exception as e:
            print(f"pdfium CFLAGS: (pkg-config not available): {e}")

        try:
            result = subprocess.run(["pkg-config", "--libs", "pdfium"], capture_output=True, text=True)
            pdfium_ldflags = result.stdout.strip()
        except Exception as e:
            print(f"pdfium LDFLAGS: (pkg-config not available): {e}")

        table_data = [
            ("PWD", str(self.wd)),
            ("DST", str(self.dst)),
            ("PKG_CONFIG_PATH", os.environ.get("PKG_CONFIG_PATH", "")),
            ("pdfium CFLAGS", pdfium_cflags),
            ("pdfium LDFLAGS", pdfium_ldflags),
            ("GOBUILDFLAGS", build_flags),
            ("CGO_CFLAGS", os.environ.get("CGO_CFLAGS", "")),
            ("", ""),
        ]

        envvars = []
        for v in os.environ:
            envvars.append(v)
        envvars.sort()
        for v in envvars:
            table_data.append((v, os.environ[v]))

        max_key_width = max(len(key) for key, _ in table_data)
        print(f"{'Variable':<{max_key_width}} | Value")
        print(f"{'-' * max_key_width}-+-{'-' * 50}")
        for key, value in table_data:
            print(f"{key:<{max_key_width}} | {value}")

    def tidy(self):
        """Format and update go.mod file"""
        self.run_command("go mod tidy -v")
        self.run_command("gofumpt -extra -w .")

    def desktop(self):
        """Compile the local backend for desktop platform"""
        self.dst.mkdir(exist_ok=True)
        build_flags = self._get_build_flags()
        cmd = [
            "go",
            "build",
            *build_flags,
            "-o",
            str(self.dst / "desktop"),
            "-tags",
            "desktop",
            str(self.wd / "demo"),
        ]
        self.run_command(cmd)

    def android(self):
        """Compile the local backend for android platform"""
        dst = self.dst / "android"
        dst.mkdir(exist_ok=True)
        build_flags = self._get_build_flags()

        # Build for different architectures
        architectures = [
            ("arm", "pdfium-android-arm"),
            ("arm64", "pdfium-android-arm64"),
            ("386", "pdfium-android-x86"),
            ("amd64", "pdfium-android-x64"),
        ]

        for arch, pdfium_dir in architectures:
            aar_file = dst / f"backend-{arch}.aar"
            jar_file = dst / f"backend-sources-{arch}.jar"

            # Remove existing files
            if aar_file.exists():
                aar_file.unlink()
            if jar_file.exists():
                jar_file.unlink()

            # Set PKG_CONFIG_PATH for this architecture
            pkg_config_path = self.wd / "lib/pdfium" / pdfium_dir
            env = os.environ.copy()
            env["PKG_CONFIG_PATH"] = str(pkg_config_path)

            cmd = [
                "gomobile",
                "bind",
                "-target",
                f"android/{arch}",
                "-androidapi",
                "23",
                "-o",
                str(aar_file),
                *build_flags,
                "-tags",
                "android",
                str(self.wd / "cmd/backend/local/android"),
            ]
            self.run_command(cmd, env=env)

    def ios(self):
        """Compile the local backend for ios platform"""
        self.dst.mkdir(exist_ok=True)
        build_flags = self._get_build_flags()

        xcframework_dir = self.dst / "Backend.xcframework"
        if xcframework_dir.exists():
            shutil.rmtree(xcframework_dir)

        # Set PKG_CONFIG_PATH for iOS
        pkg_config_path = self.wd / "lib/pdfium/pdfium-ios-device-arm64"
        env = os.environ.copy()
        env["PKG_CONFIG_PATH"] = str(pkg_config_path)

        cmd = [
            "gomobile",
            "bind",
            "-target",
            "ios",
            "-o",
            str(xcframework_dir),
            *build_flags,
            "-tags",
            "ios",
            str(self.wd / "cmd/backend/local/ios"),
        ]
        self.run_command(cmd, env=env)

    def cross(self):
        """Compile desktop backend for many platforms"""
        self.dst.mkdir(exist_ok=True)
        build_flags = self._get_build_flags()

        desktop_dir = self.wd / "cmd/backend/local/desktop"
        cmd = [
            "xgo",
            *build_flags,
            "--dest",
            str(self.dst),
            "--pkg",
            "cmd/backend/local/desktop",
            "--targets",
            "darwin/amd64,darwin/arm64,windows/386,windows/amd64",
            str(self.wd),
        ]
        self.run_command(cmd, cwd=desktop_dir)


def main():
    parser = argparse.ArgumentParser(description="Build system for ablibrary.net")
    parser.add_argument("command", nargs="?", default="help", help="Command to execute")

    args = parser.parse_args()

    build_system = BuildSystem()

    # Always fix .pc files first
    build_system.fix_pc_files()

    command = args.command.replace("-", "_")

    if hasattr(build_system, command):
        getattr(build_system, command)()
    else:
        print(f"Unknown command: {args.command}")
        build_system.help()
        sys.exit(1)


if __name__ == "__main__":
    main()
