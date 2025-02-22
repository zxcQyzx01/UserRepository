FROM golang:1.21-alpine

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod ./
RUN go mod download

COPY . .

RUN swag init -g cmd/main.go

RUN CGO_ENABLED=0 go build -o main ./cmd

EXPOSE 8080

CMD ["./main"] 