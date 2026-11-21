---
paths:
  - "backend/**"
---

# Backend Design Rules

## Design Constraints

- Dependencies point inward: business logic never imports adapters, and imports never cycle.
- Prefer declaring an interface in the package that consumes it (e.g. `api.JobService`, `job.Repo`). Contracts shared across the messaging boundary (`queue.Handler`, `queue.Publisher`) live in `transport`. `cmd/*` is the only composition root and binds implementations with `fx.As`.
- Introduce an interface where a dependency has side effects or more than one plausible implementation, not to avoid importing a pure function or value type.
- Infrastructure SDKs (e.g. NATS, Prometheus, huma, `os/exec`, DB drivers) are imported only by adapter packages, which translate between infrastructure and business logic. Interfaces that business logic depends on never expose infrastructure types.
- `platform` and `transport` never import `api` or business-logic packages.
- Domain constructors are plain Go. `fx` appears only in `cmd/*` and `platform`; dependency injection never shapes a domain API.
- Shared cross-cutting concerns (admin HTTP server, health/readiness, metrics, logging) belong in `backend/platform`, composed into each service through its `fx.Module`. A new service wires that module in and doesn't reimplement these.
- Consumer handlers wrap deterministic failures (bad input, validation and render errors) with `queue.ErrFatal` so they aren't redelivered.
- The app never emits addon-format text, and round-trip fidelity in `simc` is a non-goal: parsing may be lossy, and parse and render paths don't need to be symmetric.

## Comments

- Code should be self-describing. Add a comment only to explain why, when the code can't; needing one is a sign the code may be over-complicated.

## Validation

- Validate changes with `task backend:format backend:test backend:lint` from the repo root (`task format test lint` inside `backend/`).

## Testing

- Tests use testify/suite: a <Thing>TestSuite struct embedding suite.Suite, a Test<Thing>TestSuite runner, and one method per behaviour.
- Repetitive test cases should be parameterized and run following the pattern of (suite.Suite) Run(name string, subtest func())
- Prefer external test packages (`package x_test`); use the internal package only to test unexported behaviour.
