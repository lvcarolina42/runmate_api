# Dockerfile
FROM golang:1.22-alpine

WORKDIR /app
COPY . .

RUN go build -o runmate_api .

EXPOSE 8080

CMD ["./runmate_api"]
