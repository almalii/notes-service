FROM golang:1.21

WORKDIR /notes-app

COPY . .

RUN go build -o ./cmd/main ./cmd/

CMD ["./cmd/main"]
