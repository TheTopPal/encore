# Stage 1: build
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /bin/server ./cmd/server

# Stage 2: run
FROM alpine:3.21

COPY --from=builder /bin/server /bin/server

EXPOSE 8080

CMD ["/bin/server"]
