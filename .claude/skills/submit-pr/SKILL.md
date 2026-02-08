---
name: submit-pr
description: Summarize changes and submit a pull request using GitHub CLI
args: none
---

# Submit Pull Request Skill

When the user invokes `/submit-pr`, summarize the current branch changes and create a pull request using the GitHub CLI.

## Workflow

1. **Get current branch info**:
   - Run `git branch --show-current` to get the branch name
   - Extract the ticket ID from the branch name (e.g., `mem-22-docker-compose` → `MEM-22`)

2. **Gather changes**:
   - Run `git status` to check for uncommitted changes
   - Run `git diff main...HEAD --stat` to see file changes summary
   - Run `git diff main...HEAD` to understand the actual code changes

3. **Git commit** (if there are changes):
   - If there are staged/unstaged changes, commit them with a summary
   - If there are no files to add, skip this step and move forward

4. **Git push**:
   - Push the branch: `git push -u origin <branch-name>`

5. **Create PR with gh**:
   - Format title: `TICKET-ID: Short summary from branch name` (no # in title)
   - Example title: `MEM-22: Docker Compose and Makefile setup`
   - Description starts with the title as H1 (with #)
   - Create PR using GitHub CLI:
     ```bash
     gh pr create --title "<title>" --body "<description>"
     ```

6. **Return the PR URL** to the user

## PR Description Template

```
# TICKET-ID: Short summary (same as PR title)

<1-2 sentence summary of what this PR does>

## Changes
- <bullet point 1>
- <bullet point 2>
- <bullet point 3>
```

## Rules

1. **Commit if needed, skip if not**: If there are uncommitted changes, commit them. If there are no files to add, move forward with the PR

2. **Keep descriptions brief**: Focus on what changed, not implementation details

3. **Derive ticket ID from branch**: The branch name should follow the pattern `<ticket-id>-<description>`

4. **No co-author**: Do not add any co-author lines to commits or PRs

5. **Base branch**: Always target `main` as the base branch unless specified otherwise

6. **No emojis**: Do not use emojis in PR titles, descriptions, or any output

7. **No MCP calls**: Do not fetch Jira ticket details or call any MCP servers - derive all information from git

## Example

**Branch**: `mem-22-docker-compose-and-makefile-setup`

**PR Title**: `MEM-22: Docker Compose and Makefile setup`

**PR Description**:
```
# MEM-22: Docker Compose and Makefile setup

Adds Docker Compose configuration for local development and production, along with a Makefile for common operations.

## Changes
- Add docker-compose.yml for local development with hot reload
- Add docker-compose.prod.yml for production deployment
- Add Makefile with dev, build, and deploy commands
- Add Dockerfiles for web and api services
```
