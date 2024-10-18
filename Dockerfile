FROM golang:alpine

RUN apk add --no-cache git build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY .env .env

RUN go build -v -o server ./cmd

CMD ["/app/server"]
