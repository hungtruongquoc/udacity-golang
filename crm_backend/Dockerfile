FROM golang:1.22-alpine3.19
WORKDIR /app
COPY . .
RUN go build -o main .

EXPOSE 80
CMD ["/app/main"]