FROM golang:1.20-alpine3.16

WORKDIR /app

COPY ./ ./

RUN go build -o /bin/app main.go

ENTRYPOINT ["app"]
