# Documentation Rules

These rules apply to CLAUDE.md, files under `.claude/rules/`, READMEs, and code comments — the durable documentation surfaces, as opposed to `TODO.md` or a PR description.

- Don't restate what's directly discoverable by reading the code or directory structure (file layout, which packages exist, which addons are installed, etc.). It adds token and maintenance cost without adding information, and it's the first thing to drift out of sync with the code it's describing.
- Only write down what isn't obvious from the code itself:
    - Rationale — why a decision was made, what it trades off against
    - Constraints and gotchas that aren't visible from a single file (cross-service invariants, platform limitations)
    - Standing rules meant to constrain future decisions
- Don't document current implementation progress (what's done vs. pending, in-progress design decisions) in these docs — that belongs in `TODO.md` or a PR description. If a doc describes a desired state that isn't fully true yet, say so explicitly (a clear marker, not a blend of target-state and current-state left for the reader to untangle).
- If the same fact needs to be true in two docs, state it once and link to it rather than duplicating the explanation.
