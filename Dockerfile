FROM golang:1.21-alpine

WORKDIR /app
COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o main .
CMD ["/app/main"]
