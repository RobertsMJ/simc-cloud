# simc-cloud

Cloud application that runs [SimulationCraft](https://github.com/simulationcraft/simc) simulations. The initial phase of the project is focused on building out a system for simulating a World of Warcraft character's gear and talent combinations in various combat situations. In the long term, this project should support simulations of groups of players in realistic dungeon scenarios.

## Architecture

(TODO) A frontend interface for building gear combinations to simulate

- POSTs a request to the Job Creation service with a set of gear and gear constraints

Job Creation Service

- Creates a job and invokes the Gearset Generator with the provided set of items
- Validates the request is valid
- Returns the job ID to the frontend for status monitoring and results

(TODO) Gearset generator

- Given a list of gear options for each gear slot, generates gear combinations that satisfy the given constraints
    - Constraints could be similar to:
        - Minimum 4-piece set bonus
        - Minimum 2-piece set bonus
        - Limit 2 crafted embellishments
        - Upgrade currency (i.e. user has 60 crests, how many items can be upgraded with this currency)
- Puts the resulting valid gearsets as valid simulationcraft input strings to a queue to be simulated

Simulation Runner

- Consumes gearset simulation requests from SQS
- Runs the simulation with the given configuration
- Publishes the result to a results SQS queue for persistence

(TODO) Persistence Service

- Consumes simulation results from the results SQS queue
- Writes results to DynamoDB (result record + atomic job completion counter)
- Decoupled from the Simulation Runner so that DB write failures don't lose simulation results — unprocessable messages go to a dead-letter queue
- DynamoDB table: single-table design with `PK=job_id`, `SK=RESULT#<gearset_id>`
    - Each result item stores `statistics` (promoted from simc output for leaderboard queries without deserializing the full result blob) and the full `result` blob
- Access patterns:
    - Get all sim results for a Job by Job ID (leaderboard, sharing, bookmarking)
    - Get a single result by Job ID + Gearset ID (drilldown)

Long term:

- Group simulations with a dungeon route
- Dungeon route building, a la Mythic Dungeon Tools


## Local Development

### Prerequisites

- Go, Docker, AWS CLI, AWS SAM CLI, [Task](https://taskfile.dev)

