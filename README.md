# Cloud Native Simulationcraft Execution

For-fun educational project to exercise skills in system design, cloud architecture, and devops.
A cloud native platform for executing simulationcraft gear comparisons.

Application Scope:
1. UI for copy-pasting World of Warcraft Simc addon output. Select from available gear options to create a comparison job
1. API to parse simc addon input into a workable format
1. API to accept a simulation job with gear options
1. Event-driven service to generate valid gearsets for the simulation run
1. Event-driven simulation runner to execute the simulations
1. Persist simulation results in a schemaless DB

Platform scope:
1. Deploy a local kubernetes cluster with k3d
1. Deploy Prometheus/Grafana for basic observability
1. Use Istio service mesh for mTLS w/ kiali for observability
1. NATS Jetstream message queue - easy enough to deploy and maintain
1. ScyllaDB for job/results storage - nosql db, we'll see how it goes
1. Wrap application in a Helm chart for deployment
1. ArgoCD + gitops for deployment management
1. Tilt for development deployment

Data model:
1. Job - parent entity to group sim results, counts, and statuses of simulations for the run
1. Sim Result - status, metadata, and results of a single gearset

## Development

Prerequisites:

- Install Homebrew (`/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`)
- Install Docker (https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository)
- Install dependencies (`task install-deps`)
- Create cluster (`task cluster:create`)
