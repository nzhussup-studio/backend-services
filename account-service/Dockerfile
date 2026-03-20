FROM golang:1.24.1 AS builder
ARG TARGETOS
ARG TARGETARCH
ENV CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH}
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o bin/account-service ./cmd/api

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/account-service /app/account-service
EXPOSE 8087
CMD ["/app/account-service"]
