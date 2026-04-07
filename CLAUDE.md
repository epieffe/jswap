# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Jswap is a cross-platform CLI tool (Windows, Linux, macOS) for downloading and switching between Java JDK versions. It uses the [Adoptium API](https://api.adoptium.net/v3) to fetch Eclipse Temurin OpenJDK releases. Written in Go using the Cobra CLI framework.

## Build Commands

- `make` — build for current OS/architecture (output in `build/`)
- `make linux-amd64` / `make mac-amd64` / `make mac-arm64` / `make win-amd64` — cross-compile
- `make win-installer` — build Windows NSIS installer (requires nsis + EnVar plugin)
- `make all` — build all platforms + Windows installer
- `make clean` — remove `build/` directory

Version is injected at build time via `-ldflags` from `git describe --tags`.

There are no tests in this project.

## Architecture

**Entry point:** `jswap.go` → `cmd.Execute()`

**`cmd/`** — Cobra command definitions. Each file defines one subcommand (`get`, `set`, `ls`, `releases`, `rm`). Commands parse arguments and delegate to `internal/jdk`.

**`internal/jdk/`** — Core business logic.
- `jdk.go` — All JDK operations: install, list, set current, remove. Manages the `current-jdk` symlink via `internal/file.Link()`.
- `config.go` — Reads/writes `jswap.json` (the local state file tracking installed JDKs and which is current). Stored in `~/.jswap/` on Unix, `%LocalAppData%/Jswap/` on Windows.

**`internal/jdk/adoptium/`** — Adoptium API client.
- `adoptium.go` — API calls (release names, latest asset, specific release) and archive download+extraction logic.
- `types.go` — JSON response types for the Adoptium API.
- `system/` — Build-tag-constrained files (`system_linux_amd64.go`, `system_darwin_*.go`, `system_windows_amd64.go`) that define `OS` and `ARCH` constants matching Adoptium API values.

**`internal/file/`** — Filesystem utilities.
- `path.go` — Platform-aware paths (`JswapData()`, `JavaHome()`, `TempDir()`).
- `link_unix.go` / `link_windows.go` — Platform-specific symlink creation. Windows falls back to junctions, then full copy if symlinks fail.
- `archive.go` — Extracts `.tar.gz` and `.zip` archives.

**`internal/web/`** — HTTP utilities: generic JSON fetcher (`FetchJson[T]`) and file downloader with progress output.

## Key Design Patterns

- **Platform abstraction via build tags:** OS/architecture constants in `internal/jdk/adoptium/system/` and symlink behavior in `internal/file/` use Go build tags, not runtime checks.
- **No admin required:** Data stored in user-local directories (`~/.jswap` or `%LocalAppData%/Jswap`). JDK switching works by symlinking `current-jdk` → the selected JDK path.
- **Generics for HTTP:** `web.FetchJson[T]` uses Go generics to deserialize API responses.
