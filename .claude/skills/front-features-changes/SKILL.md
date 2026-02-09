---
name: front-features-changes
description: Work on front-end React/TypeScript features based on Jira tickets
args: ticket-id
---

# Front-End Features and Changes Skill

When the user invokes `/front-features-changes --ticket-id <ID>`, fetch the Jira ticket details via MCP and work on the front-end implementation.

## Input

- `--ticket-id`: The Jira ticket ID (e.g., `MEM-1` or just `1`). If only a number is provided, prefix with `MEM-`.

## Front-End Context

### Stack
- **Framework**: React with TypeScript
- **Styling**: Tailwind CSS with shadcn/ui components
- **Location**: `/web` directory

### Directory Structure
```
web/src/
├── components/    # Reusable UI components
├── pages/         # Page-level components and routes
├── hooks/         # Custom React hooks
├── lib/           # Utilities, helpers, and SRS algorithm previews
└── styles/        # Global styles and Tailwind configuration
```

### UI Guidelines
- Dark theme is the default
- Mobile-first responsive design
- WCAG 2.1 AA accessibility compliance
- Touch swipe gestures for mobile card review
- Keyboard shortcuts on desktop (1/2/3 for rating, space to flip, e to edit)

### Component Library
- Use shadcn/ui components as the foundation
- Follow existing component patterns in `components/`
- Maintain consistent spacing and color usage via Tailwind

### Key Features to Be Aware Of
- **Rich text editor**: Tiptap or Lexical for card content
- **Card types**: basic, cloze deletion, MCQ, image occlusion, audio, reverse
- **Media support**: images, audio, LaTeX (KaTeX), code highlighting (Shiki)
- **PWA**: Service worker for offline review capability

## Workflow

1. **Fetch the Jira ticket**: Use the Atlassian MCP tool to fetch the ticket details:
   - Call `mcp__atlassian__getJiraIssue` with `cloudId: "45d5c0c7-c3c3-468a-a8fb-7f19720b0424"` and `issueIdOrKey: "<ticket-id>"`
   - Extract: title (summary), description, tasks, acceptance criteria
   - If the ticket cannot be found, inform the user and stop

2. **ALWAYS sync with main and use git worktree**: Before creating a new branch, you MUST ensure you're on main with the latest code and use git worktree for isolation:
   - Run: `cd /Users/israelmendesdasilva/Desktop/code/memwright` (go to main worktree)
   - Run: `git checkout main`
   - Run: `git pull`
   - **NEVER create a branch from another feature branch**

3. **Create a feature branch with worktree**: Based on the ticket title, create a branch name and worktree:
   - Format: `<ticket-id-lowercase>-<title-slug>`
   - Convert the title to lowercase, replace spaces with hyphens, remove special characters
   - Example: Ticket `MEM-1` with title "Project Scaffolding and Docker Setup" → branch `mem-1-project-scaffolding-and-docker-setup`
   - Clean up existing worktree if present: `git worktree remove ../memwright-<branch-name> --force 2>/dev/null || true`
   - Create new worktree: `git worktree add ../memwright-<branch-name> -b <branch-name>`
   - Change to worktree: `cd ../memwright-<branch-name>`

4. **Analyze the ticket**: Review the requirements and identify:
   - Which components need to be created or modified
   - Which pages are affected
   - Required hooks or utilities
   - Any new dependencies needed

5. **Explore the codebase**: Before implementing, explore the existing front-end code to:
   - Understand current patterns and conventions
   - Find similar components to reference
   - Identify shared utilities to reuse

6. **Plan the implementation**: Create a task list based on the ticket tasks, mapping each to specific files and changes

7. **Implement the changes**:
   - Follow existing code style and patterns
   - Use TypeScript strictly (no `any` types unless absolutely necessary)
   - Write components as functional components with hooks
   - Ensure responsive design works on mobile and desktop
   - Add appropriate accessibility attributes

8. **Verify acceptance criteria**: Check each acceptance criterion from the ticket

9. **Stop for review**: Present a summary of changes to the user

**IMPORTANT**: Do NOT run any git commands after making code changes. No `git add`, no `git commit`, no `git push`, no staging, no pull requests. Only create the branch at the start, then make code changes and stop.

## Rules

1. **Always explore first**: Before writing code, use the Explore agent to understand the current state of the front-end codebase

2. **Match existing patterns**: Follow the conventions already established in the codebase

3. **No hardcoded strings**: Use constants or i18n keys for user-facing text where applicable

4. **Accessibility**: Include proper ARIA labels, keyboard navigation, and focus management

5. **Responsive design**: Test that changes work on mobile viewports (mobile-first approach)

6. **Component reuse**: Check for existing components before creating new ones

7. **Type safety**: Define proper TypeScript interfaces for props and state

8. **TypeScript best practices**:
    - Use strict mode, avoid `any` types
    - Prefer `interface` for object shapes, `type` for unions/intersections
    - Use discriminated unions for state management
    - Leverage type inference where possible
    - Export types alongside components that use them

9. **React best practices**:
    - Use functional components with hooks
    - Keep components small and focused (single responsibility)
    - Lift state up only when necessary
    - Use `useMemo` and `useCallback` sparingly, only for actual performance issues
    - Prefer composition over prop drilling
    - Co-locate related code (component, styles, types, tests)

10. **Tailwind best practices**:
    - Use utility classes directly, avoid unnecessary abstractions
    - Follow mobile-first responsive design (sm:, md:, lg:)
    - Use consistent spacing scale (p-4, m-2, gap-3)
    - Leverage shadcn/ui components as base
    - Extract repeated patterns into reusable components, not CSS classes

11. **Avoid unnecessary comments**:
    - Don't add comments that just restate what the code does
    - Only add comments for non-obvious business logic
    - Let the code be self-documenting through clear naming
    - No redundant JSDoc on self-explanatory props

12. **Git worktree workflow**:
    - Always use git worktree for feature branches to avoid context switching
    - Each worktree is isolated with its own working directory, node_modules, build artifacts
    - Branches are global across worktrees, so commits/pushes work the same way
    - Worktrees enable parallel development: work on multiple features simultaneously
    - Clean up worktrees after PR merge (user will handle this separately)

## Example Interaction

**User**: /front-features-changes --ticket-id MEM-5

**Assistant**:
1. Fetches ticket MEM-5 from Jira via MCP
2. Displays the ticket summary:
   - Title: "Dashboard Stats Component"
   - Description: Create a stats dashboard showing review metrics...
   - Tasks: Create StatsCard component, Implement charts...
3. Syncs with main: `cd /Users/israelmendesdasilva/Desktop/code/memwright && git checkout main && git pull`
4. Creates worktree: `git worktree remove ../memwright-mem-5-dashboard-stats-component --force 2>/dev/null || true && git worktree add ../memwright-mem-5-dashboard-stats-component -b mem-5-dashboard-stats-component && cd ../memwright-mem-5-dashboard-stats-component`
5. Explores the codebase and plans implementation
6. Implements changes following React/TypeScript patterns
7. Stops for user review
