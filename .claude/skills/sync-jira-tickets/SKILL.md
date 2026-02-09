---
name: sync-jira-tickets
description: Read markdown files from docs/jira/ and sync them to Jira using MCP
---

# Sync Jira Tickets Skill

When the user invokes `/sync-jira-tickets`, read Jira ticket markdown files from `docs/jira/` folder and create/update them in Jira using the Atlassian MCP.

## Source Files

- **Location**: `docs/jira/`
- **File pattern**: `MEM-*.md`
- **Format**: Standard Jira ticket markdown format (see create-jira-tickets skill)

## Jira Configuration

Reference from CLAUDE.md:
- **Project Key**: `MEM`
- **Cloud ID**: `45d5c0c7-c3c3-468a-a8fb-7f19720b0424`
- **Default Issue Type**: Story
- **Story Points Field**: `customfield_10016`
- **Sprint Field**: `customfield_10020`

## Workflow

### 1. Discovery Phase
1. List all markdown files in `docs/jira/` folder
2. Parse each file to extract ticket information
3. Check which tickets already exist in Jira (by key)

### 2. Parsing Phase
For each markdown file, extract:
- **Ticket ID**: From filename (e.g., `MEM-001.md` -> `MEM-1`)
- **Title**: From the `# [TICKET-XXX] Title` header
- **Priority**: From the `## Priority` section
- **Story Points**: From the `## Story Points` section (number only)
- **Description**: Full content of `## Description` section
- **Tasks**: List from `## Tasks` section (preserve markdown formatting)
- **Acceptance Criteria**: List from `## Acceptance Criteria` section

### 3. Sync Phase
For each parsed ticket:

**If ticket does NOT exist in Jira:**
1. Create new issue using `mcp__atlassian__createJiraIssue`
2. Set fields:
   - `projectKey`: "MEM"
   - `issueTypeName`: "Story"
   - `summary`: Ticket title (without the [TICKET-XXX] prefix)
   - `description`: Combined description, tasks, and acceptance criteria in markdown
   - `additional_fields`: `{"customfield_10016": <story_points>}`

**If ticket ALREADY exists in Jira:**
1. Report it as "already exists" (skip by default)
2. If user requests `--force` or `--update`, use `mcp__atlassian__editJiraIssue` to update

### 4. Cleanup Phase
After successful sync:
1. Delete markdown files for tickets that were successfully created in Jira
2. Keep markdown files for tickets that failed or were skipped
3. Use Bash tool with `rm` command to delete the files

### 5. Report Phase
After sync, report:
- Tickets created (with Jira URLs)
- Tickets skipped (already exist)
- Tickets failed (with error details)
- Files deleted from docs/jira/
- Total story points synced

## Commands

- `/sync-jira-tickets` - Sync all new tickets (skip existing)
- `/sync-jira-tickets --all` - Sync all tickets, updating existing ones
- `/sync-jira-tickets MEM-001 MEM-002` - Sync specific tickets only
- `/sync-jira-tickets --dry-run` - Show what would be synced without making changes

## Description Format for Jira

When creating the Jira description, combine sections:

```markdown
{Description paragraphs}

## Tasks

{Tasks list with checkboxes}

## Acceptance Criteria

{Acceptance criteria list with checkboxes}
```

## Error Handling

- If MCP is not connected, prompt user to run `/mcp` first
- If a file cannot be parsed, report the error and continue with other files
- If Jira API fails, report the error with details and continue

## Example Usage

```
User: /sync-jira-tickets

Claude: Found 5 markdown files in docs/jira/:
- MEM-001.md
- MEM-002.md
- MEM-003.md
- MEM-004.md
- MEM-005.md

Syncing to Jira...

Created:
- MEM-1: Project Scaffolding (8 pts) - https://israelmendez232.atlassian.net/browse/MEM-1
- MEM-2: Database Schema (10 pts) - https://israelmendez232.atlassian.net/browse/MEM-2

Skipped (already exist):
- MEM-3, MEM-4, MEM-5

Deleted from docs/jira/:
- MEM-001.md
- MEM-002.md

Summary: 2 created, 3 skipped, 2 files deleted, 18 story points synced
```