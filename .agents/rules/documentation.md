---
paths:
  - "**/*.md"
---

# Documentation Rules

- Don't restate what's directly discoverable by reading the code or directory structure (file layout, which packages exist, etc.).
- Only write down what isn't obvious from the code itself:
    - Rationale — why a decision was made, what it trades off against
    - Constraints and gotchas that aren't visible from a single file (cross-service invariants, platform limitations)
    - Standing rules meant to constrain future decisions
- Don't document current implementation progress
- If the same fact needs to be true in two docs, state it once and link to it rather than duplicating the explanation.
