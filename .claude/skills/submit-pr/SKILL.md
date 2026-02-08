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
   - Run `git log main..HEAD --oneline` to see all commits
   - Run `git diff main...HEAD --stat` to see file changes summary
   - Run `git diff main...HEAD` to understand the actual code changes

3. **Create PR title**:
   - Format: `# TICKET-ID: Short summary from branch name`
   - Example: `# MEM-22: Docker Compose and Makefile setup`
   - Convert ticket ID to uppercase
   - Derive summary from branch name (convert hyphens to spaces, capitalize appropriately)

4. **Create PR description**:
   - Write 1-2 sentences summarizing what this PR accomplishes
   - Add bullet points for key changes (3-5 bullets max)
   - Keep it concise and focused

5. **Submit the PR**:
   - Ensure all changes are committed
   - Push the branch: `git push -u origin <branch-name>`
   - Create PR using GitHub CLI:
     ```bash
     gh pr create --title "<title>" --body "<description>"
     ```

6. **Return the PR URL** to the user

## PR Description Template

```
<1-2 sentence summary of what this PR does>

## Changes
- <bullet point 1>
- <bullet point 2>
- <bullet point 3>
```

## Rules

1. **Always check for uncommitted changes**: If there are uncommitted changes, warn the user and ask if they want to commit first

2. **Keep descriptions brief**: Focus on what changed, not implementation details

3. **Derive ticket ID from branch**: The branch name should follow the pattern `<ticket-id>-<description>`

4. **No co-author**: Do not add any co-author lines to commits or PRs

5. **Base branch**: Always target `main` as the base branch unless specified otherwise

6. **No emojis**: Do not use emojis in PR titles, descriptions, or any output

7. **No MCP calls**: Do not fetch Jira ticket details or call any MCP servers - derive all information from git

## Example

**Branch**: `mem-22-docker-compose-and-makefile-setup`

**PR Title**: `# MEM-22: Docker Compose and Makefile setup`

**PR Description**:
```
Adds Docker Compose configuration for local development and production, along with a Makefile for common operations.

## Changes
- Add docker-compose.yml for local development with hot reload
- Add docker-compose.prod.yml for production deployment
- Add Makefile with dev, build, and deploy commands
- Add Dockerfiles for web and api services
```
