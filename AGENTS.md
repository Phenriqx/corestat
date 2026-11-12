# Agent instructions

These instructions apply to the Corestat repository.

## Project context

- Corestat is a Linux desktop system monitor built with Go, Wails, and Svelte.
- Treat `DESIGN.md` as the source of truth for the intended interface layout. Where it leaves a decision open, ask the maintainer before establishing a new layout direction.

## Working rules

- Run the relevant tests after completing a task. Use `make test` for Go changes. For frontend work, use an appropriate Makefile target once one exists; ask before adding or using a new project command when the Makefile does not cover it. Report what was run and any verification gap.
- Use the Makefile as the entry point for project-related terminal commands, including building, running, and testing the application. Keep new recurring project commands in the Makefile.
- Consult the maintainer before major changes or changes with significant project impact. Examples include architecture or public Go–frontend contract changes, dependency changes, platform scope changes, and changes to process-control behavior. Explain the proposed change and its impact before proceeding.
- Ask permission before modifying sensitive files, including credentials, keys, secrets, environment files, and configuration that contains sensitive values. Do not expose their contents in logs or responses.
- Follow clean architecture principles in Go code: keep responsibilities clear, separate system-data access from application decisions and Wails bindings where practical, make dependencies explicit, and keep functions focused and testable. Apply these principles incrementally; a broad restructuring requires consultation under the rule above.
- Preserve existing user work and avoid unrelated edits.

## Interface work

- Read `DESIGN.md` before changing layout or visual hierarchy.
- Distinguish the current interface from the intended design. Do not treat the current Svelte markup as an approved future layout.
- Update `DESIGN.md` when the maintainer approves a new layout decision, then implement against that decision.
