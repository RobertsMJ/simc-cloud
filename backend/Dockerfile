# syntax=docker/dockerfile:1
FROM simulationcraftorg/simc:latest@sha256:f1a582dc7d9e0712e43aa6d02341b5f84027d7cd3f57d34e010d8d3a2bea213b AS simc
FROM golang:1.27.3-alpine3.24 AS build-deps
WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download


ARG SERVICE
FROM build-deps AS build
COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 \
    GOOS=linux \
    go build \
    -ldflags="-s -w" \
    -o main ./cmd/${SERVICE}/


FROM alpine:3.24 AS simc-runtime
WORKDIR /app
RUN \
    apk --no-cache add  \
    libcurl \
    libgcc \
    libstdc++
COPY --from=simc /app/SimulationCraft/simc simc


FROM simc-runtime AS prod-simc
WORKDIR /app
COPY --from=build /app/main main
ENTRYPOINT ["./main"]
EXPOSE 8081

# For services that don't use simc at runtime
FROM alpine:3.24 AS prod
WORKDIR /app
COPY --from=build /app/main main
ENTRYPOINT ["./main"]
EXPOSE 8081
