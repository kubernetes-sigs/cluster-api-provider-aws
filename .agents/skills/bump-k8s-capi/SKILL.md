---
name: bump-k8s-capi
description: Bump Kubernetes and Cluster API to new minor versions across the entire project. Use when upgrading k8s (e.g. from 1.34 to 1.35) and CAPI (e.g. from 1.12 to 1.13). Handles go.mod, Makefile, versions.mk, metadata, e2e configs, CRDs, and migration-guide-driven code adaptations.
argument-hint: <k8s-minor> <capi-minor> (e.g. 1.35 1.13)
allowed-tools: Bash(go:*) Bash(git:*) Bash(grep:*) Bash(find:*) Bash(gh:*) Bash(make:*) Bash(curl:*)
---

Bump Kubernetes and Cluster API Versions
=========================================

Target versions: **$ARGUMENTS**

Follow the instructions in docs/book/src/development/bump-k8s-capi.md exactly, using `$ARGUMENTS` as the target Kubernetes and CAPI minor versions.
