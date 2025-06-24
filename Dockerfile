FROM golang:1.24.0-alpine AS buider

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/auth/main.go

FROM alpine:3.22.0

WORKDIR /app

COPY --from=buider /app/server .
COPY --from=buider /app/prod.env .

EXPOSE 40051

CMD ["./server"]