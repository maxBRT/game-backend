# AI Context

This document provides operational context for AI agents working on this project. It defines the Atlas planning framework, file formats, and workflows for effective collaboration.

## Overview

This project uses **Atlas**, a planning framework that stores all planning data as markdown files with YAML frontmatter in the `.atlas` directory. This keeps planning artifacts version-controlled alongside code.

## Directory Structure

```
.atlas/
├── application.md      # Project overview, goals, architecture
├── context.md          # This file - AI operational context
├── specs/              # Feature specifications
│   ├── 001-feature.md
│   └── 002-feature.md
└── tasks/              # Granular development tasks
    ├── T-101.md
    └── T-102.md
```

## File Formats

### Application Spec (`application.md`)

The central document describing the project. Contains:
- Problem statement
- Proposed solution
- Core features
- Technical architecture
- Data model
- UI/UX requirements

### Feature Specification (`specs/*.md`)

Detailed design documents for individual features. Naming convention: `NNN-feature-name.md`

**Structure:**
```markdown
# Feature Name

## 1. Description
[What this feature does and why]

## 2. UI/UX
[User interface requirements and interactions]

## 3. Backend API
[API endpoints, request/response formats]

## 4. File System
[Any file operations or storage requirements]
```

### Task (`tasks/*.md`)

Granular, actionable work items. Naming convention: `T-NNN.md`

**Required YAML frontmatter:**
```yaml
---
id: T-101
type: task
status: todo          # todo | in-progress | done
priority: high        # low | medium | high
parent_spec: 001-feature.md
---
```

**Body structure:**
```markdown
# Task Title

[Brief description of what needs to be done]

## Sub-tasks:

- [ ] **Sub-task name**: Description of the sub-task
- [ ] **Another sub-task**: More details here
```

## AI Workflows

### Creating Tasks from a Feature Spec

When asked to break down a feature specification into tasks:

1. **Read the spec thoroughly** - Understand all requirements including UI, API, and file system interactions
2. **Identify work boundaries** - Separate backend work (API endpoints, data logic) from frontend work (components, state, interactions)
3. **Create atomic tasks** - Each task should be completable in a single focused session
4. **Link to parent spec** - Always set `parent_spec` in the frontmatter
5. **Use checklists** - Break each task into verifiable sub-tasks

### Updating Task Status

When working on or completing tasks:
- Change `status: todo` to `status: in-progress` when starting work
- Change `status: in-progress` to `status: done` when complete
- Mark sub-task checkboxes as complete: `- [x]`

### Creating New Specs

When asked to create a feature specification:
1. Use the next available number prefix (e.g., `008-feature.md`)
2. Follow the standard spec structure (Description, UI/UX, Backend API, File System)
3. Be specific about API contracts and data formats
4. Consider edge cases and error handling

## Valid Field Values

| Field | Valid Values |
|-------|--------------|
| `type` | `task` |
| `status` | `todo`, `in-progress`, `done` |
| `priority` | `low`, `medium`, `high` |
| `parent_spec` | Filename from `specs/` directory |

## Best Practices

- **Task granularity**: A task should take roughly 1-4 hours of focused work
- **Clear acceptance criteria**: Each sub-task should be verifiable
- **Dependency awareness**: Note if a task depends on another being completed first
- **Incremental progress**: Prefer multiple small tasks over one large task
