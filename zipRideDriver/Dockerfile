# build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates && update-ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/zipride ./cmd

# runtime stage
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates
COPY --from=builder /bin/zipride /usr/local/bin/zipride
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/zipride"]
