FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN go install github.com/a-h/templ/cmd/templ@latest

COPY . .

RUN templ generate

RUN go build -o main .

EXPOSE 3000

CMD ["./main"]
