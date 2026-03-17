# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.0]

### Added

- Add `--allow-extra-frontmatter` flag to suppress warnings for non-spec
  frontmatter fields ([#27]). Useful for teams that embed custom metadata
  (e.g. internal tags or routing hints) alongside standard skill fields.
- Add `--allow-flat-layouts` flag to support skills that keep all files at
  the root instead of using `references/`, `scripts/`, and `assets/`
  subdirectories ([#23]). When enabled, root-level files are treated as
  standard content for token counting and orphan detection rather than
  flagged as extraneous.

### Changed

- Both new flags are available on `validate structure` and `check` commands.

## [1.2.1]

### Fixed

- Fix false positive in comma-separated keyword stuffing heuristic on
  multi-sentence descriptions with inline enumeration lists ([#26]).
  The heuristic now splits descriptions into sentences before checking,
  so commas in separate sentences are no longer counted together.

### Changed

- Extract keyword stuffing thresholds into named constants for easier tuning.

## [1.2.0]

### Changed

- Bump default OpenAI model to GPT 5.2.
- Add CI and review-skill examples to `examples/`.

## [1.1.0]

### Changed

- Increase model name truncation limit in eval compare report.

## [1.0.0]

First stable release. Includes the complete CLI and importable library packages.

### CLI

- `validate structure` — spec compliance, frontmatter, token counts, code fence
  integrity, internal link validation, orphan file detection, keyword stuffing
- `validate links` — external HTTP/HTTPS link validation with template URL support
- `analyze content` — content quality metrics (density, specificity, imperative ratio)
- `analyze contamination` — cross-language contamination detection and scoring
- `check` — run all deterministic checks with `--only`/`--skip` filtering
- `score evaluate` — LLM-as-judge scoring (Anthropic and OpenAI-compatible providers)
- `score report` — view and compare cached LLM scores across models
- Output formats: text, JSON, markdown
- GitHub Actions annotations via `--emit-annotations`
- `--strict` mode for CI (treats warnings as errors)
- Multi-skill directory detection
- Pre-commit hook support for all major agent platforms
- Homebrew install via `agent-ecosystem/tap`

### Library

- `orchestrate` — high-level validation coordination
- `evaluate` — LLM scoring orchestration with caching and progress reporting
- `judge` — LLM client abstraction and scoring (EXPERIMENTAL)
- `structure`, `content`, `contamination`, `links` — individual analysis packages
- `skill` — SKILL.md parsing (frontmatter + body)
- `skillcheck` — skill detection and reference file analysis
- `report` — output formatting (text, JSON, markdown, GitHub annotations)
- `types` — shared data types (`Report`, `Result`, `Level`, etc.)
- `judge.LLMClient` interface for custom LLM providers

[1.3.0]: https://github.com/agent-ecosystem/skill-validator/compare/v1.2.1...v1.3.0
[1.2.1]: https://github.com/agent-ecosystem/skill-validator/compare/v1.2.0...v1.2.1
[1.2.0]: https://github.com/agent-ecosystem/skill-validator/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/agent-ecosystem/skill-validator/compare/v1.0.0...v1.1.0
[#26]: https://github.com/agent-ecosystem/skill-validator/issues/26
[#23]: https://github.com/agent-ecosystem/skill-validator/issues/23
[#27]: https://github.com/agent-ecosystem/skill-validator/issues/27
