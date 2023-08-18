FROM golang:1.21-alpine AS builder

WORKDIR /notes-app

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . ./
RUN go build -o ./cmd/app ./cmd/

FROM alpine

COPY --from=builder /notes-app/cmd/app /
COPY config/prod.yml /config.yml

CMD ["/app"]
