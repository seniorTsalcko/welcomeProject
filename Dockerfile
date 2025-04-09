FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /welcome-app ./cmd/server


FROM alpine:latest


WORKDIR /app
COPY --from=builder /taskList ./

CMD ["./taskList"]