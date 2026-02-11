# Repository Guidelines

## Project Structure & Module Organization
- `cmd/usefultools` contains the main GUI application entrypoint.
- `cmd/upgrade` builds the standalone upgrade helper binary.
- `app/` hosts application-level features (UI flows and feature modules).
- `pkg/` and `common/` provide shared libraries (crypto, proxy, config, etc.).
- `utils/` and `helper/` include small reusable helpers.
- `resource/` and `material/` store images, icons, and UI assets used by the Fyne app.
- Build artifacts are written under `bin/` by Makefile targets.

## Build, Test, and Development Commands
- `make build`: builds macOS (arm64 + amd64) and Windows binaries, then zips them.
- `make build-mac-arm64`: native macOS arm64 build into `bin/darwin/arm64/`.
- `make build-mac-amd64`: cross-build macOS amd64 with `fyne-cross`.
- `make build-windows-amd64`: cross-build Windows amd64 with `fyne-cross`.
- `make tidy`: runs `go mod tidy` and vendors modules.
- `go test ./...`: runs all Go tests in the module.

## Coding Style & Naming Conventions
- Go formatting follows `gofmt` (standard Go tooling). Run before committing.
- Package names are short and lowercase; exported identifiers use `CamelCase`.
- Tests use Go’s `_test.go` suffix and `TestXxx` naming.
- No repo-wide linter is configured; keep code idiomatic and consistent with existing files.

## Testing Guidelines
- Tests live alongside code in `*_test.go` files (examples: `pkg/crypto/`, `pkg/proxy/`).
- Prefer table-driven tests for edge cases and protocol handling.
- Run targeted tests with `go test ./pkg/crypto -run TestName` and full suite with `go test ./...`.

## Commit & Pull Request Guidelines
- Commit subjects commonly use prefixes like `feat:`, `fix:`, `refactor:`, `chore:`.
- Both English and Chinese subjects appear; keep the subject concise and imperative.
- For PRs, include: a short summary, linked issue (if any), and screenshots for UI changes.

## Configuration & Tooling Notes
- The app is built with Fyne; ensure `fyne` and `fyne-cross` are installed when packaging.
- Versioning is embedded via Makefile `-ldflags` into `common/config.Version`.
