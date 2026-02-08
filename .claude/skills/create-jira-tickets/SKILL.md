---
name: create-jira-tickets
description: Generate Jira tickets as individual markdown files in docs/jira/ folder
---

# Create Jira Tickets Skill

When the user invokes `/create-jira-tickets`, generate Jira ticket(s) as individual markdown files in the `docs/jira/` folder.

## Usage

```
/create-jira-tickets                     # Interactive mode - will ask for details
/create-jira-tickets --context <text>    # Create ticket(s) from provided context
```

The `--context` argument allows passing a description, requirements, or any text that should be turned into one or more tickets. If the context is too large for a single ticket (>3 story points), automatically break it into multiple tickets.

## File Structure

Each ticket must be saved as a separate markdown file:
- **Location**: `docs/jira/`
- **Filename**: `{TICKET-ID}.md` (e.g., `MEM-001.md`, `MEM-002.md`)
- **One ticket per file** for easier review and version control

## Ticket Format

Each markdown file must follow this structure:

```markdown
# [TICKET-XXX] Ticket Title

## Priority
[High | Medium | Low]

## Story Points
[1-3]

## Description

[First paragraph: Provide context and background for this ticket. Explain why this work is needed and what problem it solves. Be specific about the current state and the desired outcome.]

[Second paragraph: Detail the technical approach or implementation strategy. Include any relevant constraints, dependencies, or considerations that the developer should be aware of.]

## Tasks

- [ ] Task 1
- [ ] Task 2
- [ ] Task 3
...

## Acceptance Criteria

- [ ] Criterion 1
- [ ] Criterion 2
- [ ] Criterion 3
...
```

## IMPORTANT

**NO NEED TO SEARCH FOR OTHER TICKETS, JUST CREATE THEM.** Do not query Jira or search for existing tickets before creating new ones. Simply create the markdown files directly.

## Rules

1. **No Emojis**: Never use emojis anywhere in the ticket content.

2. **One File Per Ticket**: Each ticket must be saved as its own markdown file in `docs/jira/`.

3. **Sequential Numbering**: Use sequential ticket IDs (MEM-001, MEM-002, etc.). Check existing files in `docs/jira/` to determine the next available number.

4. **Description**:
   - Must contain at least 2 paragraphs
   - First paragraph: Context, background, and business value
   - Second paragraph: Technical approach and considerations

5. **Tasks**:
   - Maximum of 10 tasks per ticket
   - Tasks should be actionable and specific
   - Use imperative mood (e.g., "Implement...", "Create...", "Add...")
   - **Do NOT include story points in individual tasks** - story points are only shown in the Story Points section

6. **Acceptance Criteria**:
   - Maximum of 10 acceptance criteria per ticket
   - Each criterion should be testable and verifiable
   - Use clear, measurable language
   - Cover both functional and non-functional requirements where applicable

7. **Story Points**:
   - **1 story point = 1 business day**
   - **Maximum 3 story points per ticket** - if total exceeds 3, the ticket MUST be broken into multiple smaller tickets
   - Each ticket must be independently deliverable
   - Story points are shown ONLY in the Story Points section, not repeated elsewhere

8. **Priority**:
   - High: Critical for release, blocking other work, or time-sensitive
   - Medium: Important but not blocking, should be done soon
   - Low: Nice to have, can be deferred if needed

## Workflow

1. If `--context` is provided, use that as the input; otherwise ask the user for details
2. Check existing files in `docs/jira/` to determine the next ticket number
3. If scope exceeds 3 story points, break into multiple tickets
4. Generate the ticket(s) in the specified markdown format
5. Save each ticket as a separate file in `docs/jira/{TICKET-ID}.md`
6. Report the created files to the user

## Output

After creating tickets, report:
- List of created files with their paths
- Summary of total story points across all created tickets
