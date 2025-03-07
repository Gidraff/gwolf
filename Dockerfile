FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

EXPOSE 8080

CMD ["./bin/app"]