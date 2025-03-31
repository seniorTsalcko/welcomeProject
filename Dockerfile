FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /welcome-app ./cmd/server


FROM alpine:latest

RUN apk add --no-cache tzdata

WORKDIR /app
COPY --from=builder /welcome-app /app/welcome-app

EXPOSE 8080
CMD ["/app/welcome-app"]