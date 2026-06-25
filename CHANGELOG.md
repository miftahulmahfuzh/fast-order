# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2026-06-25

### Added
- Multi-provider LLM support with configurable provider type (e.g. Gemini alongside OpenAI-compatible APIs)
- Kitchen-ticket UI revamp with dedicated design tokens and self-hosted fonts
- Receipt-column shell layout with live station status header
- Printed-ticket output with CSS line-reveal animation and tear caption
- Error stamp and legend on the kitchen-ticket layout
- Borderless monospace input fields with ember focus baseline
- Self-hosted SVG favicon
- Backend end-to-end order test
- Config package tests and LLM assemble logic with unit tests
- Claude context shortcuts and favicon reference in project docs
- Environment variables for the frontend Docker service

### Changed
- Decomposed prompt logic into a dedicated assemble module
- Replaced the slab generate button with a quiet generate text action
- Refactored config handling to support multiple LLM providers

[0.2.0]: https://github.com/miftahfast/fast-order/compare/v0.1.0...v0.2.0

## [0.1.0] - 2026-02-09

### Added
- Render deployment configuration for frontend and backend
- First-touch mode for initial order generation
- Playwright E2E tests for automated testing
- Modern Card UI design for improved user experience
- autoFocus to list menu field for improved UX
- Backend LLM output sanitization for WhatsApp format compliance
- Testing documentation
- Design documentation for first-touch mode and modern card UI

### Changed
- Backend port changed from 8080 to 8089
- Enhanced LLM prompt with better output formatting

### Fixed
- Nginx proxy configuration for Render backend deployment
- Nginx proxy_pass to avoid double /api path
- Proper /api prefix stripping in nginx rewrite
- envsubst variable escaping in CMD
- Shell form CMD for proper envsubst execution
- SHELL-FORMAT usage to preserve nginx variables
- $PORT env var for Render compatibility
- Trailing slash handling for BACKEND_URL in nginx proxy_pass

### Removed
- node_modules from git tracking (frontend)

[0.1.0]: https://github.com/miftahfast/fast-order/compare/4ec722602e751d55319434c8174a4869db4d3b55...v0.1.0
