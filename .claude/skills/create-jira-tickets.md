---
name: create-jira-tickets
description: Generate Jira tickets in markdown format with tasks, story points, and acceptance criteria
---

# Create Jira Tickets Skill

When the user invokes `/create-jira-tickets`, generate Jira ticket(s) in markdown format following the structure and rules below.

## Ticket Format

Each ticket must be output in the following markdown structure:

```markdown
# [TICKET-XXX] Ticket Title

## Priority
[High | Medium | Low]

## Story Points
[Total points - sum of task points]

## Description

[First paragraph: Provide context and background for this ticket. Explain why this work is needed and what problem it solves. Be specific about the current state and the desired outcome.]

[Second paragraph: Detail the technical approach or implementation strategy. Include any relevant constraints, dependencies, or considerations that the developer should be aware of.]

## Tasks

- [ ] Task 1 (X pts)
- [ ] Task 2 (X pts)
- [ ] Task 3 (X pts)
...

## Acceptance Criteria

- [ ] Criterion 1
- [ ] Criterion 2
- [ ] Criterion 3
...
```

## Rules

1. **No Emojis**: Never use emojis anywhere in the ticket content.

2. **Description**:
   - Must contain at least 2 paragraphs
   - First paragraph: Context, background, and business value
   - Second paragraph: Technical approach and considerations

3. **Tasks**:
   - Maximum of 10 tasks per ticket
   - Each task must have story points assigned (1, 2, or 3 points max)
   - If a task would require more than 3 points, split it into separate tasks
   - Tasks should be actionable and specific
   - Use imperative mood (e.g., "Implement...", "Create...", "Add...")

4. **Acceptance Criteria**:
   - Maximum of 10 acceptance criteria per ticket
   - Each criterion should be testable and verifiable
   - Use clear, measurable language
   - Cover both functional and non-functional requirements where applicable

5. **Story Points**:
   - 1 point: Simple, straightforward task (few hours)
   - 2 points: Moderate complexity task (half day to full day)
   - 3 points: Complex task requiring significant effort (more than a day but less than 2)
   - If a task exceeds 3 points, it must be broken down into smaller tasks
   - Total story points = sum of all task points

6. **Priority**:
   - High: Critical for release, blocking other work, or time-sensitive
   - Medium: Important but not blocking, should be done soon
   - Low: Nice to have, can be deferred if needed

## Workflow

1. Ask the user what feature, bug fix, or work item they want to create a ticket for
2. Gather any additional context needed (technical requirements, dependencies, etc.)
3. Generate the ticket(s) in the specified markdown format
4. If the scope is too large for a single ticket, suggest breaking it into multiple tickets
5. Output all tickets as markdown that can be copied directly

## Example Output

# [TICKET-001] Implement User Authentication API

## Priority
High

## Story Points
8

## Description

The application currently lacks a secure authentication mechanism, which prevents users from accessing personalized features and protected resources. This ticket covers the implementation of a JWT-based authentication API that will handle user login, token generation, and token validation. This is a foundational requirement for all user-specific functionality in the system.

The implementation will use industry-standard security practices including bcrypt for password hashing, JWT tokens with appropriate expiration times, and secure HTTP-only cookies for token storage. The API will be built following RESTful conventions and will integrate with the existing user database schema.

## Tasks

- [ ] Create authentication controller with login endpoint (2 pts)
- [ ] Implement password hashing utility using bcrypt (1 pt)
- [ ] Create JWT token generation service (2 pts)
- [ ] Implement token validation middleware (2 pts)
- [ ] Add logout endpoint with token invalidation (1 pt)

## Acceptance Criteria

- [ ] Users can log in with valid email and password credentials
- [ ] Invalid credentials return appropriate error messages
- [ ] JWT tokens are generated with 24-hour expiration
- [ ] Protected routes reject requests without valid tokens
- [ ] Passwords are never stored or transmitted in plain text
- [ ] API responses follow the established error format
- [ ] All endpoints are documented in the API specification
