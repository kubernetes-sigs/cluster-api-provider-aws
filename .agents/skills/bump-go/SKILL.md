---
name: bump-go
description: Bump Go to a new minor or patch version across the entire project. Use when upgrading Go (e.g. from 1.24 to 1.25). Handles go.mod, Makefile, golangci-lint, GitHub Actions, netlify.toml, and cloudbuild.
argument-hint: <go-minor-version> (e.g. 1.26)
allowed-tools: Bash(go:*) Bash(git:*) Bash(grep:*) Bash(find:*) Bash(docker:*) Bash(gh:*) Bash(make:*)
---

Bump Go Version
===============

Target Go minor version: **$ARGUMENTS**

Follow the instructions in docs/book/src/development/bump-go.md exactly, using `$ARGUMENTS` as the target Go minor version.
