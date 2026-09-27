---
paths:
  - "backend/**"
---

# Backend Design Rules

- Keep the backend platform-agnostic. Infrastructure (message queue, database, compute runtime) sits behind repository interfaces and transport adapters, so business logic moves to a new runtime without changes.
- Modularity and separation of concerns are learning goals here in their own right. Any component should be swappable without affecting the rest of the system.
- Shared cross-cutting concerns (admin HTTP server, health/readiness, metrics, logging) belong in `backend/platform`, composed into each service through its `fx.Module`. A new service wires that module in and doesn't reimplement these.

## Design Constraints

- Repository interfaces and transport handlers use domain types only (`models.*`, `context.Context`). Infrastructure types must never appear in interfaces or domain packages, only in their concrete implementations.
