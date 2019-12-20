FROM golang:latest

WORKDIR /app

COPY . .
RUN go build -o backend

EXPOSE 8080

ENTRYPOINT ["/app/backend"]