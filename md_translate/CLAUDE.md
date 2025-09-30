# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`md_translate` is a Rust project for markdown translation. The project is currently in early stages with a basic Cargo workspace structure.

## Development Commands

### Building
```bash
cargo build
```

### Running
```bash
cargo run
```

### Testing
```bash
cargo test
```

### Checking (fast compilation check)
```bash
cargo check
```

## Project Context

This is part of a larger `/tmp` monorepo containing multiple experimental projects. The parent directory includes:
- `cfg` - Configuration-related Rust/JS project
- `downloader` - Go-based parallel file downloader
- `aya` - Another experimental project
- Other small experimental projects

## Custom Claude Code Commands

This repository includes custom slash commands in `.claude/commands/` focused on article workflow management:

- `/manage_article_workflow` - Comprehensive article project management with progress tracking and parallel work optimization
- `/research_article` - Information collection and research execution
- `/plan_article` - Structure design and outline creation
- `/write_article` - Writing execution
- `/review_article` - Review and quality improvement
- `/create_plan`, `/validate_plan`, `/implement_plan` - General planning workflows
- `/research_codebase`, `/research_codebase_nt`, `/research_codebase_generic` - Codebase analysis variants
- `/commit` - Git commit assistance
- `/local_review` - Local code review
- `/debug` - Debugging assistance

These commands follow a systematic workflow for managing writing projects with stage gates, dependency tracking, and quality assurance.

## Git Workflow

This is a shared tmp directory with multiple projects. Changes should be committed thoughtfully:
- Check what's staged before committing (`git status`, `git diff`)
- The main branch is `main`
- Recent commits show work on cfg, aya, and other projects

## Architecture Notes

The Rust project uses Cargo edition "2024" (note: non-standard edition). The codebase is minimal with only `src/main.rs` currently containing a hello world implementation.