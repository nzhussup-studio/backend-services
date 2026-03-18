# BUILD STAGE
FROM golang:1.24.1 AS builder
ARG TARGETOS
ARG TARGETARCH
ENV CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH}
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bin/llm-service ./cmd/api

# RUN STAGE
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/bin/llm-service /app/llm-service
EXPOSE 8086
CMD ["/app/llm-service"]


