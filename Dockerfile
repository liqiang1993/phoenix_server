FROM golang:alpine

RUN mkdir /app
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o phoenix_server ./src/main.go
CMD ["./phoenix_server"]
