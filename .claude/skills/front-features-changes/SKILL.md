---
name: front-features-changes
description: Work on front-end React/TypeScript features based on Jira tickets
---

# Front-End Features and Changes Skill

When the user invokes `/front-features-changes`, provide context about the front-end architecture and then work on the specified Jira ticket.

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

1. **Create a new branch**: Create a new git branch for this work (e.g., `feature/TICKET-123-description`)

2. **Receive the Jira ticket**: Ask the user to provide the Jira ticket details (ticket ID, description, tasks, acceptance criteria)

3. **Analyze the ticket**: Review the requirements and identify:
   - Which components need to be created or modified
   - Which pages are affected
   - Required hooks or utilities
   - Any new dependencies needed

4. **Explore the codebase**: Before implementing, explore the existing front-end code to:
   - Understand current patterns and conventions
   - Find similar components to reference
   - Identify shared utilities to reuse

5. **Plan the implementation**: Create a task list based on the ticket tasks, mapping each to specific files and changes

6. **Implement the changes**:
   - Follow existing code style and patterns
   - Use TypeScript strictly (no `any` types unless absolutely necessary)
   - Write components as functional components with hooks
   - Ensure responsive design works on mobile and desktop
   - Add appropriate accessibility attributes

7. **Verify acceptance criteria**: Check each acceptance criterion from the ticket

8. **Stop for review**: Once all changes are complete, stop and present the changes to the user for review

**IMPORTANT**: Do NOT commit, push, or create pull requests. Only create the branch and make code changes, then stop for user review.

## Rules

1. **Always explore first**: Before writing code, use the Explore agent to understand the current state of the front-end codebase

2. **Match existing patterns**: Follow the conventions already established in the codebase

3. **No hardcoded strings**: Use constants or i18n keys for user-facing text where applicable

4. **Accessibility**: Include proper ARIA labels, keyboard navigation, and focus management

5. **Responsive design**: Test that changes work on mobile viewports (mobile-first approach)

6. **Component reuse**: Check for existing components before creating new ones

7. **Type safety**: Define proper TypeScript interfaces for props and state

## Example Interaction

**User**: /front-features-changes

**Assistant**: I'll help you work on a front-end feature or change. Please provide the Jira ticket details including:
- Ticket ID and title
- Description
- Tasks
- Acceptance criteria

Once you share the ticket, I'll analyze the requirements, explore the relevant parts of the codebase, and implement the changes following the project's patterns and conventions.
