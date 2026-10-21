# syntax=docker/dockerfile:1
ARG SIMC_VERSION=latest
FROM simulationcraftorg/simc:${SIMC_VERSION} AS simc
FROM golang:alpine AS build-deps
WORKDIR /app

COPY go.work .
COPY go.work.sum .

COPY test-utils/go.mod test-utils/go.mod

COPY simc/go.mod simc/go.mod
COPY simc/go.sum simc/go.sum

COPY backend/go.mod backend/go.mod
COPY backend/go.sum backend/go.sum

RUN go mod download


ARG SERVICE
FROM build-deps AS build
COPY simc/ simc/
COPY test-utils/ test-utils/
COPY backend/ backend/

RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 \
    GOOS=linux \
    go build \
    -ldflags="-s -w" \
    -o main ./backend/cmd/${SERVICE}/


FROM alpine:latest AS simc-runtime
WORKDIR /app
RUN \
    apk --no-cache add --virtual build_dependencies \
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
FROM alpine:latest AS prod
WORKDIR /app
COPY --from=build /app/main main
ENTRYPOINT ["./main"]
EXPOSE 8081
