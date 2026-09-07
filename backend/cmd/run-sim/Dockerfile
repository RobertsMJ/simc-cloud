ARG SIMC_VERSION=latest
FROM golang:alpine AS build-deps
WORKDIR /app

COPY go.work .

COPY test-utils/go.mod test-utils/go.mod
COPY simc/go.mod simc/go.mod

COPY backend/go.mod backend/go.mod
COPY backend/go.sum backend/go.sum

RUN go mod download


FROM build-deps AS source
WORKDIR /app
COPY simc/ simc/
COPY test-utils/ test-utils/
COPY backend/ backend/


FROM source AS dev-deps
WORKDIR /app
RUN go install github.com/air-verse/air@v1.67.4
RUN go install github.com/go-delve/delve/cmd/dlv@v1.27.1
RUN \
    apk --no-cache add --virtual build_dependencies \
    libcurl \
    libgcc \
    libstdc++
COPY --from=simulationcraftorg/simc:latest /app/SimulationCraft/simc simc


FROM dev-deps AS dev
ENTRYPOINT [ \
    "air", \
    "--build.cmd", "CGO_ENABLED=0 GOOS=linux go build -gcflags=all=\"-N -l\" -o main ./backend/cmd/run-sim", \
    "--build.full_bin", "dlv exec ./main --headless --listen=:2345 --api-version=2 --accept-multiclient --continue --log" \
    ]


FROM source AS build
RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build \
    -ldflags="-s -w" \
    -o main ./backend/cmd/run-sim/


FROM alpine:latest AS simc-runtime
RUN \
    apk --no-cache add --virtual build_dependencies \
    libcurl \
    libgcc \
    libstdc++
COPY --from=simulationcraftorg/simc:latest /app/SimulationCraft/simc simc


FROM simc-runtime AS prod
WORKDIR /app
COPY --from=build /app/main main
ENTRYPOINT ["./main"]