# BUILD STAGE
FROM golang:1.24.1 AS builder
ARG TARGETOS
ARG TARGETARCH
ENV CGO_ENABLED=1 GOOS=${TARGETOS} GOARCH=${TARGETARCH}
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bin/image-service ./cmd/api

# RUN STAGE
FROM alpine:latest
WORKDIR /app

# Install dependencies required by CGO
RUN apk add --no-cache libc6-compat libstdc++ gcc libgcc

COPY --from=builder /app/bin/image-service /app/image-service
COPY ./var/ ./var/
EXPOSE 8085
CMD ["/app/image-service"]


