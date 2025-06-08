# filepath: /workspaces/simc-cloud/lambdas/echo/Dockerfile
FROM golang:1.24 AS build
WORKDIR /echo

# Copy function code
COPY go.mod go.sum ./

COPY main.go .

# Build binary
RUN CGO_ENABLED=0 go build -o main main.go

FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build /echo/main ./main

ENTRYPOINT ["./main"]