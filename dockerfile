FROM golang:1.24.2-alpine AS builder

WORKDIR /build

ADD go.mod .

COPY . .

RUN go build -o authServic cmd/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /build/authServic /app/authServic
COPY --from=builder /build/config.yml /app/config.yml
COPY --from=builder /build/migrations /app/migrations

CMD [ "./authServic"]
