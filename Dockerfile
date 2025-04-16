FROM golang:1.23rc2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server . \
  && chmod +x server \
  && ls -lh server

FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/public ./public

EXPOSE 8080

CMD ["/app/server"]

