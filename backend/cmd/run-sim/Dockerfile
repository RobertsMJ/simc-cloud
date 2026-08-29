ARG SIMC_VERSION=latest
ARG GOLANG_VERSION=1.26.7
FROM golang:${GOLANG_VERSION} AS build
WORKDIR /app

COPY go.work .

COPY test-utils/go.mod test-utils/go.mod
COPY simc/go.mod simc/go.mod

COPY backend/go.mod backend/go.mod
COPY backend/go.sum backend/go.sum

RUN go mod download

COPY simc/ simc/
COPY test-utils/ test-utils/
COPY backend/ backend/

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build \
    -tags lambda.norpc \
    -o main ./backend/cmd/run-sim/

FROM alpine:latest AS prod-deps
RUN \
    apk --no-cache add --virtual build_dependencies \
        libcurl \
        libgcc \
        libstdc++

FROM simulationcraftorg/simc:${SIMC_VERSION} AS simc
FROM prod-deps
WORKDIR /app
COPY --from=simc /app/SimulationCraft/simc simc
COPY --from=build /app/main main
ENTRYPOINT ["./main"]