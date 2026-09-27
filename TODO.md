
- [ ] k8s infrastructure initialization
    - [x] Move kiali to a proper helm install
    - [x] move prometheus to a proper helm install
    - [x] add grafana
    - [ ] add argocd
    - [ ] deploy nats
    - [ ] deploy scylladb
    - [ ] better solution for gateway-api CRD installation in-cluster
- [ ] rearchitect backend around k8s deployment
    - [x] add fx
    - [ ] add config
    - [x] add metrics endpoints
    - [ ] platformize REST api and queue subscribers
    - [ ] rearchitect run-sim
    - [ ] rearchitect result-writer
        - [ ] decide atomic result-write + counter-increment strategy for ScyllaDB — DynamoDB's `TransactWriteItems` (used today in `backend/db/job.go`) has no direct equivalent: counter columns can't share a table with regular columns, and logged batches can't mix counter/non-counter mutations even within a single partition
    - [ ] rearchitect jobs

