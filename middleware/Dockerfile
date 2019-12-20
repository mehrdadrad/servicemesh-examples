FROM golang:latest

WORKDIR /app

COPY . .
RUN go build -o middleware

EXPOSE 8081

ENTRYPOINT ["/app/middleware"]