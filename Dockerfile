FROM golang:1.20-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/api

FROM alpine:3.18

RUN apk add --no-cache ca-certificates tzdata

RUN adduser -D -H -h /app appuser
WORKDIR /app

COPY --from=builder /app/api .

USER appuser

EXPOSE 8080

CMD ["./api"]