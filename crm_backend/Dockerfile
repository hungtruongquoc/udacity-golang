# Build stage
FROM golang:1.22-alpine3.19 as builder
WORKDIR /app
COPY . .
RUN go build -o main .
RUN apk add curl
# Set up server stage
FROM alpine:3.19.1
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 80
CMD ["/app/main"]