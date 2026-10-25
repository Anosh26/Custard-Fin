# Git Workflow Rule

## Strict Human Approval for Git Mutations
- **NEVER** run `git commit` or `git push` automatically or proactively.
- All Git commits and pushes must be explicitly requested or confirmed by the human before executing.
- The agent may inspect repository status with read-only commands (`git status`, `git diff`, `git log`).
- When changes are ready, show the files changed or git diff, suggest a commit message, and let the human make or approve the commit and push.
