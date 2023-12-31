# Build stage
FROM golang:1.21.5-alpine3.18 AS builder
WORKDIR /app 
COPY . .
RUN go env -w GOPROXY=https://goproxy.io,direct
RUN go build -o main main.go

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY *.env .
COPY db/migration db/migration

EXPOSE 8080 9090
CMD [ "/app/main"]