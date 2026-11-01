---
paths:
  - "backend/**"
---

# Backend Design Rules

- Keep the backend platform-agnostic. Infrastructure (message queue, database, compute runtime) sits behind repository interfaces and transport adapters, so business logic moves to a new runtime without changes.
- Modularity and separation of concerns are learning goals here in their own right. Any component should be swappable without affecting the rest of the system.
- Shared cross-cutting concerns (admin HTTP server, health/readiness, metrics, logging) belong in `backend/platform`, composed into each service through its `fx.Module`. A new service wires that module in and doesn't reimplement these.

## Testing

- Tests use testify/suite: a <Thing>TestSuite struct embedding suite.Suite, a Test<Thing>TestSuite runner, and one method per case.
- Test packages are external (package x_test).

## Design Constraints

- Repository interfaces and transport handlers use domain types only (`models.*`, `context.Context`). Infrastructure types must never appear in interfaces or domain packages, only in their concrete implementations.


## Example:

api → domain → transport ← platform
platform and transport never import api or domain packages.
cmd/* is the only composition root, the one place that sees every layer.
