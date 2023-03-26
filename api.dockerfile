# syntax=docker/dockerfile:1

## Build
FROM golang:1.20.2-alpine3.17 AS build

WORKDIR /app

# Install dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source code
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg

# Build the binary
RUN go build -o /api ./cmd/api

## Deploy
FROM scratch

# Copy our static executable
COPY --from=build /api /api
#COPY backend.env /

EXPOSE 9091

# USER nonroot:nonroot

# Run the binary
ENTRYPOINT ["/api"]
