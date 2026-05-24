Serverless SIMC runner

Prerequisites:

- Install Docker (https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository)
- Install Go (`sudo snap install go --classic`)
- Install AWS CLI (`sudo snap install aws-cli --classic`)
- `aws login`
- Install AWS SAM (`brew install aws-sam-cli`)

[Rough architecture](https://excalidraw.com/#json=TYNR-QuLm4pt3ouZwGs1D,nWcATGd-Jg-XzRfETBOrRQ)


Test invocations:
- task sim-profile-sqs -- profiles/sample.simc