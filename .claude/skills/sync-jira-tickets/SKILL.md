---
name: sync-jira-tickets
description: Read markdown files from docs/jira/ and sync them to Jira using MCP
---

# Sync Jira Tickets Skill

When the user invokes `/sync-jira-tickets`, read Jira ticket markdown files from `docs/jira/` folder and create/update them in Jira using the Atlassian MCP.

## Source Files

- **Location**: `docs/jira/`
- **File pattern**: `*.md`
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

### 2. Parsing Phase
For each markdown file, extract:
- **Title**: From the `# [TICKET-XXX] Title` or `# Title` header
- **Priority**: From the `## Priority` section
- **Story Points**: From the `## Story Points` section (number only)
- **Description**: Full content of `## Description` section
- **Tasks**: List from `## Tasks` section (preserve markdown formatting)
- **Acceptance Criteria**: List from `## Acceptance Criteria` section

### 3. Sync Phase
For each parsed ticket:

1. Create new issue using `mcp__atlassian__createJiraIssue`
2. Set fields:
   - `projectKey`: "MEM"
   - `issueTypeName`: "Story"
   - `summary`: Ticket title (without the [TICKET-XXX] prefix if present)
   - `description`: Combined description, tasks, and acceptance criteria in markdown
   - `additional_fields`: `{"customfield_10016": <story_points>}`

### 4. Cleanup Phase
After successful sync:
1. Delete markdown files for tickets that were successfully created in Jira
2. Keep markdown files for tickets that failed or were skipped
3. Use Bash tool with `rm` command to delete the files

### 5. Report Phase
After sync, report:
- Tickets created (with Jira URLs)
- Tickets failed (with error details)
- Files deleted from docs/jira/
- Total story points synced

## Commands

- `/sync-jira-tickets` - Sync all tickets from docs/jira/ to Jira
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

Claude: Found 3 markdown files in docs/jira/:
- Backend Integration Tests.md
- Backend Subdeck Support.md
- Frontend Subdeck Visualization.md

Syncing to Jira...

Created:
- MEM-27: Backend Integration Tests with Database (3 pts) - https://israelmendez232.atlassian.net/browse/MEM-27
- MEM-28: Backend Subdeck Support - Nested Deck Hierarchy (3 pts) - https://israelmendez232.atlassian.net/browse/MEM-28
- MEM-29: Frontend Subdeck Visualization and Navigation (3 pts) - https://israelmendez232.atlassian.net/browse/MEM-29

Deleted from docs/jira/:
- Backend Integration Tests.md
- Backend Subdeck Support.md
- Frontend Subdeck Visualization.md

Summary: 3 created, 3 files deleted, 9 story points synced
```