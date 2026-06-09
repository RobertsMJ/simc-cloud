# Backend Design Rules

- The Go backend is designed to be platform-agnostic: AWS infrastructure (Lambda, SQS, DynamoDB) is isolated behind repository interfaces and transport adapters so the core business logic can be ported to other runtimes (e.g. Kubernetes + Kafka + Postgres) without modification.
- The backend should be designed with a focus on modularity and separation of concerns as a learning exercise for the developer. Any component should be able to be swapped out without affecting the overall system.
- The backend should be designed with a focus on performance and scalability to handle high traffic loads efficiently.
- The backend should be designed to be maintainable and testable, with clear interfaces and effective tests.

## Design Constraints

- Repository interfaces and transport handlers must be defined in terms of domain types only (`models.*`, `context.Context`). AWS-specific types (`events.SQSEvent`, `dynamodbav` tags, etc.) must never appear in interfaces or domain packages — only in their concrete implementations.
