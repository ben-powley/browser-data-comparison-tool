FROM golang:latest

LABEL maintainer="Ben Powley <ben.powley97@icloud.com>"

WORKDIR /app

ADD data/* /data/

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]