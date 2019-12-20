FROM golang:latest

WORKDIR /app

COPY . .
RUN go get ./...
RUN go build . -o frontend

EXPOSE 8082

ENTRYPOINT ["/app/frontend"]