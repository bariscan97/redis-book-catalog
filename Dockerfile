FROM golang:alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY .env .env
RUN go build -v -o server ./cmd

CMD ["/app/server"]