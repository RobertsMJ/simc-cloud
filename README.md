# Serverless SIMC runner

A serverless simulationcraft runner for comparing World of Warcraft gear.

## Development

Prerequisites:

- Install Docker (https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository)
- Install Go (`sudo snap install go --classic`)
- Install AWS CLI (`sudo snap install aws-cli --classic`)
- `aws login`
- Install AWS SAM (`brew install aws-sam-cli`)

Test invocations:
- task sim-profile-sqs -- profiles/sample.simc
