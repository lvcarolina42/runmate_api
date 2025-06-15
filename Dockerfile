FROM golang:latest AS builder
WORKDIR /app

COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/main ./

EXPOSE 3000
CMD ["./main"]
