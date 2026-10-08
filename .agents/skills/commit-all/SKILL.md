---
name: commit-all
description: Commit pending repository changes in logical groups with evidence-based messages and limited diff inspection. Use when asked to commit all changes or organize changes into commits.
---

# Commit All

Create separate commits for logical groups unless the user specifies another grouping or a single commit.

## Select changes

- Inspect the complete staged, unstaged, and untracked file inventory using Git status and name/status summaries. Handle filenames safely, including spaces and renames.
- If anything is unstaged or untracked, ask which changes to commit before staging or committing anything. Distinguish staged and unstaged changes, including files with both. Invoking this skill alone does not answer this question; an explicit selection already provided by the user does.
- If everything is staged, proceed with the staged changes without confirmation. Preserve changes outside the selected scope.

## Understand and group

- Start with paths, change types, and diff statistics. Group related changes based on this evidence and explicit user context. Keep tests and generated sources with the changes they support when the evidence establishes that relationship; do not split a coherent feature merely by directory or file type.
- Keep analysis lightweight. If there are too many selected files for a brief review, skip diff content entirely. Judge size from file count and change volume; dozens of changed files normally qualify.
- For smaller inventories, inspect only a little diff content when needed: a bounded sample of relevant hunks, never a full review. Do not repeatedly sample until the entire diff has effectively been read.
- Choose logical groups unless the user has specified them. One coherent group warrants one commit. Ask about ambiguous grouping only when it prevents proceeding reliably.
- Write concise messages based only on what you understood from inspected evidence or explicit user context. Do not invent intent, behavior, fixes, verification, or outcomes. Prefer a narrower factual message when that is all the evidence supports.
- If completely unsure what a group changes and unable to write an honest message, stop before committing it. Tell the user you cannot do it and they must commit those changes themselves. Do not override the limited-analysis rule to manufacture certainty.

## Commit

- Commit each selected group with its own message. Use scoped Git operations so each commit contains exactly the selected changes for that group. Preserve staged and unstaged content excluded by the selection, including partial staging within a file.
- Do not discard work, amend existing commits, bypass hooks, or push unless separately requested. If a commit fails, stop and explain the blocker instead of blindly retrying or modifying unrelated files.
- Verify each successful commit's included files. Count only commits created during this invocation. Reassess unexpected working-tree or index changes before continuing.

## Response

After successful completion, output exactly `made X commits`, replacing `X` with the number created, including zero. Add no report, hashes, messages, or explanation. Required selection questions, uncertainty notices, and failure explanations must still be communicated when they arise; do not hide them behind a success-only response.
