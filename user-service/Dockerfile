FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o user ./user-service/main.go

EXPOSE 8082

CMD ["./user"]