# Multi-stage build for CustardFin Cloud API (apid)
FROM golang:1.26-alpine AS builder

WORKDIR /build

# Install git/build tools if needed
RUN apk --no-cache add ca-certificates git

# Copy root workspace and dependent modules
COPY go.work ./
COPY libs/contracts ./libs/contracts
COPY apps/api ./apps/api

# Build static binary
WORKDIR /build/apps/api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /build/apid ./cmd/apid

# Production minimal runner
FROM alpine:3.20

RUN apk --no-cache add ca-certificates curl tzdata

WORKDIR /app
COPY --from=builder /build/apid /app/apid

EXPOSE 8080

ENTRYPOINT ["/app/apid"]
