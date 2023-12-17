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
COPY app.env .
RUN mkdir -p db/schema
COPY db/schema/0000_init_schema.up.sql db/schema

EXPOSE 8080
CMD [ "/app/main"]